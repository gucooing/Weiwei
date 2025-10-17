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

package v1

import (
	"errors"

	"github.com/gucooing/weiwei/pkg/util"
	"github.com/gucooing/weiwei/pkg/util/crypt"
)

type ClientConfig struct {
	Log       *Log        `json:"log" yaml:"log" toml:"log"`
	Auth      *AuthConfig `json:"auth" toml:"auth" yaml:"auth"`
	WeisNet   *Net        `json:"weisNet" yaml:"weisNet" toml:"weisNet"`
	WeicLogin *WeicLogin  `json:"weicLogin" yaml:"weicLogin" toml:"weicLogin"`
}

func (c *ClientConfig) Init() error {
	if c == nil {
		return errors.New("config is nil")
	}
	c.Log.Init()
	c.Auth.Init()
	c.WeicLogin = util.EmptyDefault(c.WeicLogin, &WeicLogin{CryptType: string(crypt.CryptTypeNone)})
	c.WeisNet.Init()
	c.WeicLogin.Init()

	return nil
}
