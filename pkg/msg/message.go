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
	"errors"
	"reflect"
)

type Message = interface{}

var (
	msgCmdMap map[uint16]reflect.Type
	cmdMsgMap map[reflect.Type]uint16

	ErrMsgCmd = errors.New("msg cmd is invalid")
)

func RegisterMsg(cmdId uint16, msg Message) {
	msgCmdMap[cmdId] = reflect.TypeOf(msg)
	cmdMsgMap[reflect.TypeOf(msg)] = cmdId
}

func GetMessageByCmdId(cmdId uint16) (Message, error) {
	t, ok := msgCmdMap[cmdId]
	if !ok {
		return nil, ErrMsgCmd
	}
	return reflect.New(t).Interface(), nil
}
