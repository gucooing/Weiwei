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

var (
	xorKeySize int = 4096
)

type XOR struct {
	Seed   int64
	XorKey []byte
}

func newCryptXor(conf interface{}) (Crypt, error) {
	x := new(XOR)
	seed, ok := conf.(int64)
	if !ok {
		return x, errors.New("conf err")
	}
	x.Seed = seed
	x.XorKey = SeedNewXorKey(seed)

	return x, nil
}

func SeedNewXorKey(seed int64) []byte {
	xorKey := make([]byte, xorKeySize)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < 4096>>3; i++ {
		binary.BigEndian.PutUint64(xorKey[i<<3:], r.Uint64())
	}
	return xorKey
}

func (x *XOR) Encryption(data []byte) (encrypted []byte, err error) {
	util.Xor(data, x.XorKey)
	return data, nil
}

func (x *XOR) Decrypt(encrypted []byte) (decrypted []byte, err error) {
	util.Xor(encrypted, x.XorKey)
	return encrypted, nil
}
