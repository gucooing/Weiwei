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

package net

import (
	"bufio"
	"encoding/binary"
	"net"

	"github.com/gucooing/weiwei/pkg/util/crypt"
)

type TCPListener struct {
	net.Listener
}

func NewTCPListener(address string) (*TCPListener, error) {
	t := new(TCPListener)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	t.Listener = listener

	return t, nil
}

type TCPConn struct {
	conn  net.Conn
	buf   *bufio.Reader
	crypt crypt.Crypt
}

func (l *TCPListener) Accept() (Conn, error) {
	c := new(TCPConn)
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.buf = bufio.NewReader(conn)

	return c, nil
}

var (
	tcpLenSize = 2
)

func (c *TCPConn) Read() (n int, b []byte, err error) {
	lenBytes := make([]byte, tcpLenSize)
	n, err = c.buf.Read(lenBytes)
	if err != nil {
		return
	}
	headLen := binary.BigEndian.Uint16(lenBytes)

	b = make([]byte, headLen)

	n, err = c.buf.Read(b)
	if err != nil {
		return
	}
	b, err = c.crypt.Decrypt(b)
	if err != nil {
		return
	}

	return
}

func (c *TCPConn) Close() error {
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

func (c *TCPConn) SetCrypt(crypt crypt.Crypt) {
	c.crypt = crypt
}
