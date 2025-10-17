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
	"context"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util/backoff"
)

type Control struct {
	// service context
	ctx context.Context
	// conn and weic network conn
	conn net.Conn
	// runId
	runId int64
	// doneChan
	doneChan chan struct{}
}

func NewControl(ctx context.Context) (*Control, error) {
	c := &Control{
		ctx:      ctx,
		doneChan: make(chan struct{}),
	}

	return c, nil
}

func (c *Control) Run() {
	go c.keepController()

	<-c.doneChan
}

func (c *Control) keepController() {
	backoff.BackoffStart(
		func() bool {
			_, err := msg.WriteMsg(c.conn, &msg.PingReq{
				ClientTimestamp: time.Now().UnixNano(),
			})
			if err != nil {
				slog.Errorf("to weis pingReq err:%s", err.Error())
				close(c.doneChan)
				return true
			}
			return false
		},
		c.doneChan,
		time.Second,
	)
}
