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
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed        = errors.New("pool is closed")
	ErrPoolExhausted = errors.New("connection pool exhausted")
	ErrPoolTimeout   = errors.New("connection pool timeout")
)

type Pooler interface {
	AddConn(conn Conn) error
	Get(ctx context.Context) (Conn, error)
	Close() error
}

type Options struct {
	Dialer func(ctx context.Context) error

	PoolSize        int
	DialTimeout     time.Duration
	ConnMaxLifetime time.Duration
}

type ConnPool struct {
	cfg *Options

	queue    chan struct{}
	connsMu  sync.Mutex
	conns    []Conn
	poolSize int

	_closed uint32 // atomic
}

func NewConnPool(opt *Options) *ConnPool {
	p := &ConnPool{
		cfg: opt,

		queue:    make(chan struct{}, opt.PoolSize),
		conns:    make([]Conn, 0, opt.PoolSize),
		poolSize: 0,
	}

	p.connsMu.Lock()
	p.checkMinIdleConns()
	p.connsMu.Unlock()

	return p
}

func (p *ConnPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *ConnPool) checkMinIdleConns() {
	if p.cfg.PoolSize-p.poolSize <= 0 {
		return
	}

	for i := 0; i < p.cfg.PoolSize-p.poolSize; i++ {
		err := p._addConn()
		if err != nil && !errors.Is(err, ErrClosed) {

		}
	}
}

func (p *ConnPool) _addConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.DialTimeout)
	defer cancel()

	err := p.cfg.Dialer(ctx)
	if err != nil {
		return err
	}

	if p.closed() {
		return ErrClosed
	}

	return nil
}

func (p *ConnPool) AddConn(conn Conn) error {
	if p.closed() {
		conn.Close()
		return ErrClosed
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	if p.poolSize >= p.cfg.PoolSize {
		conn.Close()
		return ErrPoolExhausted
	}

	p.conns = append(p.conns, conn)
	p.poolSize++
	p.freeTurn()
	return nil
}

func (p *ConnPool) Get(ctx context.Context) (Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}
	if err := p.waitTurn(ctx); err != nil {
		return nil, err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	for {
		conn, err := p.popConn()
		if err != nil {
			return nil, err
		}

		if conn == nil {
			p._addConn()
			return nil, ErrPoolExhausted
		}

		p.freeTurn()

		if !p.isHealthyConn(conn) {
			conn.Close()
		}

		return conn, err
	}
}

func (p *ConnPool) waitTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case p.queue <- struct{}{}:
		return nil
	default:
	}

	// add
	err := p._addConn()
	if err != nil {
		return err
	}

	timer := time.NewTimer(p.cfg.DialTimeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case p.queue <- struct{}{}:
		return nil
	case <-timer.C:
		return ErrPoolTimeout
	}
}

func (p *ConnPool) freeTurn() {
	<-p.queue
}

func (p *ConnPool) popConn() (conn Conn, err error) {
	if p.closed() {
		return nil, ErrClosed
	}

	n := len(p.conns)
	if n == 0 {
		return nil, nil
	}

	index := n - 1
	cn := p.conns[index]

	p.conns = p.conns[:index]
	p.poolSize--

	p.checkMinIdleConns()

	return cn, nil
}

func (p *ConnPool) isHealthyConn(cn Conn) bool {
	now := time.Now()

	if p.cfg.ConnMaxLifetime > 0 && now.Sub(cn.CreatedAt()) >= p.cfg.ConnMaxLifetime {
		return false
	}

	return true
}

func (p *ConnPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}

	var firstErr error
	p.connsMu.Lock()
	for _, conn := range p.conns {
		if err := conn.Close(); err != nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.connsMu.Unlock()

	return firstErr
}
