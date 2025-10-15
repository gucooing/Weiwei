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
	"errors"
)

var (
	ErrCryptTypeUn = errors.New("crypt type unknown")
)

type Crypt interface {
	Encryption(data []byte) (encrypted []byte, err error)
	Decrypt(encrypted []byte) (decrypted []byte, err error)
}

type CryptType string

const (
	CryptTypeNone CryptType = "none"
	CryptTypeRsa  CryptType = "rsa"

	// CryptTypeXor Only used after security verification
	CryptTypeXor CryptType = "xor"
)

func NewCrypt(cryptType CryptType, conf interface{}) (Crypt, error) {
	switch cryptType {
	case CryptTypeNone:
		return newCryptNone(conf)
	case CryptTypeRsa:
		return newCryptRsa(conf)
	case CryptTypeXor:
		return newCryptXor(conf)
	default:
		return nil, ErrCryptTypeUn
	}
}
