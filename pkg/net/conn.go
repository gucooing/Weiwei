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
	"net"
	"time"

	"github.com/gucooing/weiwei/pkg/util/compress"
	"github.com/gucooing/weiwei/pkg/util/crypt"
)

type Conn interface {
	Read() (n int, b []byte, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	SetCrypt(crypt crypt.Crypt)
	SetCompress(compress compress.Compress)
	CreatedAt() time.Time
}
