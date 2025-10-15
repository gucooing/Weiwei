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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RSA struct {
	privKey *rsa.PrivateKey
}

func newCryptRsa(conf interface{}) (Crypt, error) {
	r := new(RSA)
	pem := conf.(string)
	privKey, err := ParsePrivKeyPem([]byte(pem))
	if err != nil {
		return nil, err
	}
	r.privKey = privKey

	return r, nil
}

func (R RSA) Encryption(data []byte) (encrypted []byte, err error) {
	// TODO implement me
	panic("implement me")
}

func (R RSA) Decrypt(encrypted []byte) (decrypted []byte, err error) {
	// TODO implement me
	panic("implement me")
}

func ParsePrivKeyPem(privKeyPem []byte) (privKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(privKeyPem)
	if block == nil {
		return nil, errors.New("invalid rsa private key")
	}
	privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
