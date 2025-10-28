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
	"errors"
	"io"
	"math"
	"net"
)

type TCPListener struct {
	net.Listener
}

func NewTCPListener(network, address string) (*TCPListener, error) {
	t := new(TCPListener)
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	t.Listener = listener

	return t, nil
}

func (l *TCPListener) Close() error {
	if l.Listener != nil {
		l.Listener.Close()
	}
	return nil
}

type TCPConn struct {
	*baseConn
	net.Conn
	buf *bufio.Reader
}

func (l *TCPListener) Accept() (Conn, error) {
	c := &TCPConn{
		baseConn: newBaseConn(),
	}
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	c.Conn = conn
	c.buf = bufio.NewReader(conn)

	return c, nil
}

var (
	tcpLenSize = 4
)

func (c *TCPConn) Read() (n int, bin []byte, err error) {
	lenBytes := make([]byte, tcpLenSize)
	if _, err = io.ReadFull(c.buf, lenBytes); err != nil {
		return 0, nil, err
	}
	headLen := binary.BigEndian.Uint32(lenBytes)

	buf := make([]byte, int(headLen))
	if _, err = io.ReadFull(c.buf, buf); err != nil {
		return 0, nil, err
	}
	bin, err = c.BaseRead(buf)
	if err != nil {
		return
	}
	n = tcpLenSize + int(headLen)
	return
}

func (c *TCPConn) Close() error {
	if c.Conn != nil {
		c.Conn.Close()
	}
	return nil
}

func TcpDial(network, address string) (Conn, error) {
	c := &TCPConn{
		baseConn: newBaseConn(),
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	c.Conn = conn
	c.buf = bufio.NewReader(conn)

	return c, nil
}

func (c *TCPConn) Write(b []byte) (n int, err error) {
	bin, err := c.BaseWrite(b)
	if err != nil {
		return 0, err
	}
	headLen := len(bin)

	if headLen > math.MaxUint32 {
		return 0, errors.New("data too large")
	}

	buf := c.getBuffer(tcpLenSize + headLen)
	defer c.putBuffer(buf)

	binary.BigEndian.PutUint32(buf[:tcpLenSize], uint32(headLen))
	copy(buf[tcpLenSize:], bin)

	totalWritten := 0
	for totalWritten < len(buf) {
		n, err = c.Conn.Write(buf[totalWritten:])
		if err != nil {
			return totalWritten, err
		}
		totalWritten += n
	}
	return
}
