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
	"github.com/gucooing/weiwei/pkg/util/compress"
	"github.com/gucooing/weiwei/pkg/util/crypt"
)

type baseConn struct {
	crypt    crypt.Crypt
	compress compress.Compress
}

func newBaseConn() *baseConn {
	b := &baseConn{
		crypt:    crypt.CryptNone,
		compress: compress.CompressNone,
	}

	return b
}

func (b *baseConn) SetCrypt(crypt crypt.Crypt) {
	b.crypt = crypt
}

func (b *baseConn) SetCompress(compress compress.Compress) {
	b.compress = compress
}
