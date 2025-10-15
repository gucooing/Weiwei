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

package net

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/msg"
)

var (
	msgCmdSize = 2

	ErrNetWorkNu = errors.New("network unknown")
	ErrCmdSize   = fmt.Errorf("msg len err")
)

func Listen(network, address string) (listener Listener, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		listener, err = NewTCPListener(address)
	default:
		return nil, ErrNetWorkNu
	}
	return
}

func ReadMsg(conn Conn) (message msg.Message, err error) {
	n, buffer, err := conn.Read()
	if err != nil {
		return nil, err
	}
	if n <= msgCmdSize {
		return nil, ErrCmdSize
	}
	cmdId := binary.BigEndian.Uint16(buffer[0:msgCmdSize])
	message, err = msg.GetMessageByCmdId(cmdId)
	if err != nil {
		slog.Tracef("read msg :%s", buffer)
		return nil, err
	}
	err = json.Unmarshal(buffer[msgCmdSize:], message)
	return
}
