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

package msg

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gucooing/weiwei/pkg/net"
)

var (
	msgCmdSize = 2

	ErrCmdSize = fmt.Errorf("msg cmd size err")
)

func ReadMsg(conn net.Conn) (message Message, err error) {
	n, buffer, err := conn.Read()
	if err != nil {
		return nil, err
	}
	if n <= msgCmdSize {
		return nil, ErrCmdSize
	}
	cmdId := binary.BigEndian.Uint16(buffer[0:msgCmdSize])
	message, err = GetMessageByCmdId(cmdId)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer[msgCmdSize:], message)
	return
}

func WriteMsg(conn net.Conn, message Message) (n int, err error) {
	cmdId, err := GetCmdIdByMessage(message)
	if err != nil {
		return 0, err
	}
	data, err := json.Marshal(message)
	if err != nil {
		return 0, err
	}
	buffer := make([]byte, msgCmdSize+len(data))
	binary.BigEndian.PutUint16(buffer[:msgCmdSize], cmdId)
	copy(buffer[msgCmdSize:], data)
	n, err = conn.Write(buffer)
	if err != nil {
		return 0, err
	}
	return n, err
}
