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

	"github.com/panjf2000/gnet/v2"
)

type GTCPListener struct {
	gnet.BuiltinEventEngine

	connCh  chan gnet.Conn
	closeCh chan struct{}
}

func NewGTCPListener(network, address string) (*GTCPListener, error) {
	t := &GTCPListener{
		BuiltinEventEngine: gnet.BuiltinEventEngine{},
		connCh:             make(chan gnet.Conn),
		closeCh:            make(chan struct{}),
	}
	go gnet.Run(
		t,
		"tcp://"+address,
		gnet.WithMulticore(true),
	)

	return t, nil
}

func (l *GTCPListener) Close() error {
	close(l.closeCh)
	return nil
}

type GTCPConn struct {
	*baseConn
	gnet.Conn
	buf *bufio.Reader
}

func (l *GTCPListener) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	select {
	case <-l.closeCh:
		return nil, gnet.Close
	default:
		l.connCh <- c
		return
	}
}

func (l *GTCPListener) Accept() (Conn, error) {
	c := &GTCPConn{
		baseConn: newBaseConn(),
	}
	conn, ok := <-l.connCh
	if !ok {
		return nil, io.EOF
	}
	c.Conn = conn
	c.buf = bufio.NewReader(conn)

	return c, nil
}

func (c *GTCPConn) Close() error {
	if c.Conn != nil {
		c.Conn.Close()
	}
	return nil
}

func (c *GTCPConn) Read() (n int, b []byte, err error) {
	lenBytes := make([]byte, tcpLenSize)
	if _, err = io.ReadFull(c.buf, lenBytes); err != nil {
		return 0, nil, err
	}
	headLen := binary.BigEndian.Uint32(lenBytes)

	buf := c.getBuffer(int(headLen))
	defer c.putBuffer(buf)
	if _, err = io.ReadFull(c.buf, buf); err != nil {
		return 0, nil, err
	}
	b, err = c.BaseRead(buf)
	if err != nil {
		return
	}
	n = tcpLenSize + int(headLen)
	return
}

func (c *GTCPConn) Write(b []byte) (n int, err error) {
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
