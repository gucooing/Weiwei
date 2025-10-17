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

package config

import (
	"encoding/json"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"

	v1 "github.com/gucooing/weiwei/pkg/config/v1"
)

var (
	Server *v1.ServerConfig
	Client *v1.ClientConfig
)

func LoadConfig(b []byte, c any) error {
	if err := toml.Unmarshal(b, c); err == nil {
		return err
	}
	if err := yaml.Unmarshal(b, c); err == nil {
		return err
	}
	return json.Unmarshal(b, &c)
}

func loadConfFile(path string) ([]byte, error) {
	buff, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func LoadServerConfig(path string) error {
	buff, err := loadConfFile(path)
	if err != nil {
		return err
	}
	Server = new(v1.ServerConfig)
	if err := LoadConfig(buff, Server); err != nil {
		return err
	}
	Server.Init()
	return nil
}

func LoadClientConfig(path string) error {
	buff, err := loadConfFile(path)
	if err != nil {
		return err
	}
	Client = new(v1.ClientConfig)
	if err := LoadConfig(buff, Client); err != nil {
		return err
	}
	return Client.Init()
}
