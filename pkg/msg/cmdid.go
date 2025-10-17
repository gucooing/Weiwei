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
	"reflect"
)

const (
	loginReq = iota + 1
	loginRsp
	pingReq
	pingRsp
)

func init() {
	msgCmdMap = make(map[uint16]reflect.Type)
	cmdMsgMap = make(map[reflect.Type]uint16)

	RegisterMsg(loginReq, LoginReq{})
	RegisterMsg(loginRsp, LoginRsp{})
	RegisterMsg(pingReq, PingReq{})
	RegisterMsg(pingRsp, PingRsp{})
}
