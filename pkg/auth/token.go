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

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/gucooing/weiwei/pkg/util"
)

type Token struct {
	Token string
}

func NewToken(token string) (*Token, error) {
	t := &Token{
		Token: token,
	}
	return t, nil
}

func (t *Token) SetVerifyLogin(timestamp int64) string {
	return t.GetAuthKey(t.Token, timestamp)
}

func (t *Token) VerifyLogin(timestamp int64, loginKey string) error {
	if strings.Compare(loginKey, t.GetAuthKey(t.Token, timestamp)) == 0 {
		return nil
	}
	return errors.New("invalid auth key")
}

func (t *Token) GetAuthKey(token string, timestamp int64) string {
	shaCtx := sha256.New()
	data := []byte(token)
	timeBytes := []byte(strconv.FormatInt(timestamp, 10))
	if len(timeBytes) > len(data) {
		timeBytes = timeBytes[len(timeBytes)-len(data):]
	}
	util.Xor(data, timeBytes)
	shaCtx.Write(data)
	bin := shaCtx.Sum(nil)
	return hex.EncodeToString(bin)
}
