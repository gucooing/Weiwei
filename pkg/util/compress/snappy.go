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

package compress

import (
	"github.com/golang/snappy"
)

type Snappy struct{}

func (*Snappy) Compress(src []byte) ([]byte, error) {
	data := snappy.Encode(nil, src)
	return data, nil
}

func (*Snappy) Decompress(src []byte) ([]byte, error) {
	return snappy.Decode(nil, src)
}
