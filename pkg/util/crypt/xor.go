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

package crypt

import (
	"encoding/binary"
	"errors"
	"math/rand"

	"github.com/gucooing/weiwei/pkg/util"
)

type XOR struct {
	seed   int64
	xorKey []byte
}

func newCryptXor(conf interface{}) (Crypt, error) {
	x := new(XOR)
	seed, ok := conf.(int64)
	if !ok {
		return x, errors.New("conf err")
	}
	x.seed = seed
	x.xorKey = make([]byte, 4096)
	r := rand.New(rand.NewSource(x.seed))
	for i := 0; i < 4096>>3; i++ {
		binary.BigEndian.PutUint64(x.xorKey[i<<3:], r.Uint64())
	}

	return x, nil
}

func (x *XOR) Encryption(data []byte) (encrypted []byte, err error) {
	util.Xor(data, x.xorKey)
	return data, nil
}

func (x *XOR) Decrypt(encrypted []byte) (decrypted []byte, err error) {
	util.Xor(encrypted, x.xorKey)
	return encrypted, nil
}
