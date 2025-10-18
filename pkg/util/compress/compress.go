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
	"errors"
)

var (
	CompressNone   = &None{}
	CompressSnappy = &Snappy{}
	CompressGzip   = &Gzip{}

	ErrCompressTypeNu = errors.New(`compress type nu`)
)

type Compress interface {
	Compress(src []byte) ([]byte, error)
	Decompress(src []byte) ([]byte, error)
}

type CompressType string

const (
	CompressTypeNone   CompressType = "none"
	CompressTypeGzip   CompressType = "gzip"
	CompressTypeSnappy CompressType = "snappy"
	CompressTypeZlib   CompressType = "zlib"
	CompressTypeBrotli CompressType = "brotli"
	CompressTypeLzo    CompressType = "lzo"
)

func NewCompress(compressType CompressType) (Compress, error) {
	switch compressType {
	case CompressTypeNone:
		return CompressNone, nil
	case CompressTypeGzip:
		return CompressGzip, nil
	case CompressTypeSnappy:
		return CompressSnappy, nil
	default:
		return nil, ErrCompressTypeNu
	}
}
