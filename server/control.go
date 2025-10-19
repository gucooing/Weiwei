// Copyright 2025 gucooing, gucooing@alsl.xyz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/auth"
	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util"
	"github.com/gucooing/weiwei/pkg/util/backoff"
)

var (
	ErrRepeatControl = errors.New("repeat control")
)

type ControlManager struct {
	mu       sync.Mutex
	contrils map[int64]*Control
}

func NewControlManager() *ControlManager {
	cm := &ControlManager{
		contrils: make(map[int64]*Control),
	}
	return cm
}

func (cm *ControlManager) AddControl(runId int64, control *Control) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.contrils[runId]; ok {
		return ErrRepeatControl
	}

	cm.contrils[runId] = control
	return nil
}

func (cm *ControlManager) GetControl(runId int64) (*Control, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cry, ok := cm.contrils[runId]
	return cry, ok
}

func (cm *ControlManager) DelControl(runId int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cry, ok := cm.contrils[runId]; ok {
		cry.Close()
		delete(cm.contrils, runId)
	}
}

func (cm *ControlManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	var lastErr error
	for _, control := range cm.contrils {
		err := control.Close()
		if err != nil {
			lastErr = err
		}
	}
	cm.contrils = nil

	return lastErr
}

type Control struct {
	// conn client net conn
	conn net.Conn
	// runId client id
	runId int64
	// seed
	seed int64
	// dispatcher msg handler
	dispatcher *msg.Dispatcher
	// lasePing lase ping time time.Time
	lasePing atomic.Value
	// doneChan
	doneChan chan struct{}
	// net conn pool
	connPool net.Pooler
	// work verifier
	workVerifier auth.Verifier
}

func NewControl(conn net.Conn) *Control {
	c := &Control{
		conn:       conn,
		runId:      util.NewRunId(),
		dispatcher: msg.NewDispatcher(conn),
		lasePing:   atomic.Value{},
		doneChan:   make(chan struct{}),
	}
	c.seed = rand.Int63n(time.Now().UnixNano() ^ c.runId)
	c.lasePing.Store(time.Now())

	// dispatcher
	c.dispatcher.RegisterMsg(&msg.CSPingReq{}, c.handlerPing)

	// pool
	c.connPool = net.NewConnPool(&net.Options{
		Dialer:          c.reqAddWorkConn,
		PoolSize:        10,
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: 24 * time.Hour,
	})

	c.workVerifier = auth.NewToken(strconv.FormatInt(c.seed, 10))

	slog.Infof("addr:%s runId:%v new weic",
		c.conn.RemoteAddr().String(), c.runId)
	return c
}

func (c *Control) Start() {
	go c.keepController()
	go c.dispatcher.Start()

	// close
	<-c.dispatcher.DoneChan()
	close(c.doneChan)
}

func (c *Control) keepController() {
	if config.Server.WeicTimeout <= 0 {
		return
	}
	backoff.BackoffStart(
		func() error {
			if time.Since(c.lasePing.Load().(time.Time)) >
				time.Duration(config.Server.WeicTimeout)*time.Second {

				slog.Infof("runId:%v weic timeout", c.runId)
				c.conn.Close()
				return nil
			}
			return errors.New("")
		},
		c.doneChan,
		&backoff.ExponentialBackoff{
			BaseInterval: 1 * time.Second,
			MaxRetries:   0,
			MaxInterval:  1 * time.Second,
		},
	)
}

func (c *Control) Close() error {
	slog.Infof("addr:%s runId:%v weic stop",
		c.conn.RemoteAddr().String(), c.runId)

	err := c.conn.Close()

	return err
}

func (c *Control) reqAddWorkConn(ctx context.Context) error {
	err := c.dispatcher.Send(&msg.SCAddWorkConnReq{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Control) addWorkConn(conn net.Conn, req *msg.CSAddWorkConnRsp) error {
	// auth
	if err := c.workVerifier.VerifyLogin(req.Timestamp, req.LoginKey); err != nil {
		return err
	}

	// add
	err := c.connPool.AddConn(conn)
	if err != nil {
		slog.Errorf("runId:%v addWorkConn err:%v", c.runId, err)
	}

	return err
}
