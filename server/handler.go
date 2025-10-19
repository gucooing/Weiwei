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
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/msg"
)

func (c *Control) handlerPing(rawMsg msg.Message) {
	req := rawMsg.(*msg.CSPingReq)

	clientTime := time.Unix(0, req.ClientTimestamp)
	serverTime := time.Now()

	_, err := msg.WriteMsg(c.conn, &msg.SCPingRsp{
		ClientTimestamp: req.ClientTimestamp,
		ServerTimestamp: serverTime.UnixNano(),
	})
	if err != nil {
		slog.Errorf("runId:%v weic pingRsp write err: %s", c.runId, err.Error())
		return
	}
	c.lasePing.Store(time.Now())
	slog.Tracef("runId:%v weic ping:%s", c.runId, serverTime.Sub(clientTime).String())
}
