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
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/msg"
)

func (c *Control) sendPingReq() error {
	_, err := msg.WriteMsg(c.conn, &msg.CSPingReq{
		ClientTimestamp: time.Now().UnixNano(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Control) handlerPing(rawMsg msg.Message) {
	rsp := rawMsg.(*msg.SCPingRsp)

	clientTime := time.Unix(0, rsp.ClientTimestamp)
	serverTime := time.Unix(0, rsp.ServerTimestamp)

	slog.Tracef("weis ping:%s", serverTime.Sub(clientTime).String())
}
