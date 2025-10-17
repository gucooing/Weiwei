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

package server

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/auth"
	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/env"
	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util/crypt"
)

const (
	connReadTimeout time.Duration = 10 * time.Second
)

var (
	ErrWeicLoginTime = errors.New("weic login timeout")
	ErrUnknownClient = errors.New("unknown client")
)

type Service struct {
	// ctx
	ctx context.Context
	// cancel
	cancel context.CancelFunc

	// weiListener service discovery listener
	weiListener net.Listener

	// multiListener business listener
	multiListener *net.MultiListener

	// weicLoginVerifier weic login auth
	weicLoginVerifier auth.Verifier
}

func NewService() (*Service, error) {
	slog.Debugf("new server service...")
	s := new(Service)

	slog.Debugf("new weiListener...")
	wln, err := net.Listen(config.Server.WeiNet.Network, config.Server.WeiNet.Address)
	if err != nil {
		slog.Tracef("new server listen err: %v", err)
		return nil, err
	}
	slog.Debugf("network:%s address:%s new weiListener success", config.Server.WeiNet.Network, config.Server.WeiNet.Address)
	s.weiListener = wln

	slog.Debugf("new multiListener...")

	slog.Debugf("server service success")
	return s, nil
}

func (svr *Service) Run(ctx context.Context) {
	slog.Infof("server service run")
	ctx, cancel := context.WithCancel(ctx)
	svr.ctx = ctx
	svr.cancel = cancel

	svr.mainHandle()
	<-svr.ctx.Done()
	// service context
	svr.Close()
}

func (svr *Service) Close() {
	slog.Debugf("server service close...")

	slog.Debugf("server service close success")
}

func (svr *Service) mainHandle() {
	slog.Debugf("run server service mainHandle")
	for {
		conn, err := svr.weiListener.Accept()
		if err != nil {
			slog.Printf("server service weiListener accept err:%v", err)
			return
		}
		conn.SetCrypt(config.Server.WeicLogin.Crypt)
		go func(conn net.Conn) {
			lerr := svr.newConn(conn)
			if lerr != nil {
				conn.Close()
				slog.Errorf("addr:%s new conn  err:%v", conn.RemoteAddr().String(), lerr)
			}
		}(conn)
	}
}

func (svr *Service) newConn(conn net.Conn) error {
	ctx := context.Background()
	loginCtx, cancel := context.WithTimeout(ctx, connReadTimeout)
	defer cancel()

	select {
	case <-loginCtx.Done():
		ctx.Err()
		return ErrWeicLoginTime
	default:
		rawMsg, err := msg.ReadMsg(conn)
		if err != nil {
			return err
		}
		switch m := rawMsg.(type) {
		case *msg.LoginReq: // new weic
			return svr.loginWeic(conn, m)
		default:
			return ErrUnknownClient
		}
	}

}

func (svr *Service) loginWeic(conn net.Conn, loginReq *msg.LoginReq) error {
	// auth
	if err := config.Server.Auth.Verifier.VerifyLogin(loginReq); err != nil {
		return err
	}
	slog.Debugf("addr:%s loginReq version:%s token:%s",
		conn.RemoteAddr().String(), loginReq.Version, loginReq.LoginKey)
	// new weic
	cl := NewControl(conn)
	loginRsp := &msg.LoginRsp{
		Version: env.Version,
		Seed:    rand.Int63n(time.Now().UnixNano() ^ cl.runId),
		RunId:   cl.runId,
	}
	cry, err := crypt.NewCrypt(crypt.CryptTypeXor, loginRsp.Seed)
	if err != nil {
		return err
	}
	defer conn.SetCrypt(cry)
	slog.Debugf("addr:%s loginRsp version:%s runId:%v seed:%v",
		conn.RemoteAddr().String(), loginRsp.Version, loginRsp.RunId, loginRsp.Seed)
	_, err = msg.WriteMsg(conn, loginRsp)
	if err != nil {
		return err
	}
	go cl.Start()
	return nil
}
