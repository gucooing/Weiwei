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
	tcpLenSize = 2
)

func (c *TCPConn) Read() (n int, b []byte, err error) {
	lenBytes := make([]byte, tcpLenSize)
	_, err = c.buf.Read(lenBytes)
	if err != nil {
		return
	}
	headLen := binary.BigEndian.Uint16(lenBytes)

	buf := c.getBuffer(int(headLen))
	defer c.putBuffer(buf)
	_, err = c.buf.Read(buf)
	if err != nil {
		return
	}
	b, err = c.BaseRead(buf)
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
	b, err = c.BaseWrite(b)
	if err != nil {
		return 0, err
	}
	headLen := len(b)

	if headLen > 0xFFFF {
		return 0, errors.New("data too large")
	}

	buf := c.getBuffer(tcpLenSize + headLen)
	defer c.putBuffer(buf)

	binary.BigEndian.PutUint16(buf[:tcpLenSize], uint16(headLen))
	copy(buf[tcpLenSize:], b)

	n, err = c.Conn.Write(buf[:tcpLenSize+headLen])
	return
}
