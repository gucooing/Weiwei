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
)

type ClientConfig struct {
	Log           *Log        `json:"log" yaml:"log" toml:"log"`
	ServerNetwork string      `json:"serverNetwork" yaml:"serverNetwork" toml:"serverNetwork"`
	ServerAddr    string      `json:"serverAddr" yaml:"serverAddr" toml:"serverAddr"`
	Auth          *AuthConfig `json:"auth" toml:"auth" yaml:"auth"`
}

func (c *ClientConfig) Init() error {
	if c == nil {
		return errors.New("config is nil")
	}

	c.Log.Init()
	c.Auth.Init()

	return nil
}
