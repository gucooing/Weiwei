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

package client

import (
	"errors"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util/backoff"
)

type Control struct {
	// conn and weic network conn
	conn net.Conn
	// runId
	runId int64
	// seed
	seed int64
	// dispatcher msg handler
	dispatcher *msg.Dispatcher
	// doneChan
	doneChan chan struct{}
}

func NewControl(conn net.Conn) (*Control, error) {
	c := &Control{
		conn:       conn,
		dispatcher: msg.NewDispatcher(conn),
		doneChan:   make(chan struct{}),
	}
	// dispatcher
	c.dispatcher.RegisterMsg(&msg.PingRsp{}, c.handlerPing)
	slog.Infof("new weis control")
	return c, nil
}

func (c *Control) Run() {
	go c.keepController()
	go c.dispatcher.Start()

	<-c.dispatcher.DoneChan()
	close(c.doneChan)
	c.conn.Close()
	slog.Infof("weis control done")
}

func (c *Control) keepController() {
	backoff.BackoffStart(
		func() error {
			// TODO Try again?
			c.sendPingReq()
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
