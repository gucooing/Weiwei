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

type ServerConfig struct {
	Log         *Log        `json:"log" yaml:"log" toml:"log"`
	ApiNetwork  string      `json:"apiNetwork" yaml:"apiNetwork" toml:"apiNetwork"`
	ApiAddress  string      `json:"apiAddress" yaml:"apiAddress" toml:"apiAddress"`
	Auth        *AuthConfig `json:"auth" toml:"auth" yaml:"auth"`
	WeicTimeout int64       `json:"weicTimeout" yaml:"weicTimeout" toml:"weicTimeout"`
}

func (s *ServerConfig) Init() error {
	if s == nil {
		return errors.New("config is nil")
	}
	s.Log.Init()
	s.Auth.Init()
	return nil
}
