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

package msg

import (
	"io"
	"reflect"

	"github.com/gucooing/weiwei/pkg/net"
)

type Dispatcher struct {
	conn  net.Conn
	close bool

	doneChan    chan struct{}
	sendChan    chan Message // send
	readChan    chan Message // read
	msgHandlers map[reflect.Type]func(Message)
}

func NewDispatcher(conn net.Conn) *Dispatcher {
	d := &Dispatcher{
		conn:        conn,
		doneChan:    make(chan struct{}),
		sendChan:    make(chan Message, 100),
		readChan:    make(chan Message, 100),
		msgHandlers: make(map[reflect.Type]func(Message)),
	}
	return d
}

func (d *Dispatcher) Start() {
	go d.sendThread()
	go d.readThread()
}

func (d *Dispatcher) sendThread() {
	for {
		select {
		case <-d.doneChan:
			return
		case rawMsg := <-d.sendChan:
			WriteMsg(d.conn, rawMsg)
		}
	}
}

func (d *Dispatcher) readThread() {
	for {
		rawMsg, err := ReadMsg(d.conn)
		if err != nil {
			close(d.doneChan)
			return
		}
		if handler, ok := d.msgHandlers[reflect.TypeOf(rawMsg).Elem()]; ok {
			handler(rawMsg)
		} else {
			// TODO
		}
	}
}

func (d *Dispatcher) DoneChan() chan struct{} {
	return d.doneChan
}

func (d *Dispatcher) SendChan() chan Message {
	return d.sendChan
}

func (d *Dispatcher) Send(msg Message) error {
	select {
	case <-d.doneChan:
		return io.EOF
	case d.sendChan <- msg:
		return nil
	}
}

func (d *Dispatcher) RegisterMsg(rawMsg Message, handler func(Message)) {
	d.msgHandlers[reflect.TypeOf(rawMsg).Elem()] = handler
}
