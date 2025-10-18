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
	"errors"
	"sync/atomic"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util"
	"github.com/gucooing/weiwei/pkg/util/backoff"
)

type Control struct {
	// conn client net conn
	conn net.Conn

	// runId client id
	runId int64

	// dispatcher msg handler
	dispatcher *msg.Dispatcher

	// lasePing lase ping time time.Time
	lasePing atomic.Value

	// doneChan
	doneChan chan struct{}
}

func NewControl(conn net.Conn) *Control {
	c := &Control{
		conn:       conn,
		runId:      util.NewRunId(),
		dispatcher: msg.NewDispatcher(conn),
		lasePing:   atomic.Value{},
		doneChan:   make(chan struct{}),
	}
	c.lasePing.Store(time.Now())

	// dispatcher
	c.dispatcher.RegisterMsg(&msg.PingReq{}, c.handlerPing)

	// keepController
	go c.keepController()
	slog.Infof("addr:%s runId:%v new weic",
		c.conn.RemoteAddr().String(), c.runId)
	return c
}

func (c *Control) Start() {
	go c.dispatcher.Start()

	// close
	<-c.dispatcher.DoneChan()
	slog.Infof("addr:%s runId:%v weic stop",
		c.conn.RemoteAddr().String(), c.runId)
	c.conn.Close()

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
