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

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/auth"
	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util/crypt"
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

	// weicLoginCrypt weic login crypt
	weicLoginCrypt crypt.Crypt

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

	slog.Debugf("new weicLoginCrypt...")
	if config.Server.WeicLogin != nil {
		s.weicLoginCrypt, err = crypt.NewCrypt(crypt.CryptTypeRsa, config.Server.WeicLogin.RsaPrivateKey)
		if err != nil {
			slog.Tracef("new weicLoginCrypt error:%v", err)
			return nil, err
		}
	}
	slog.Debugf("new weicLoginCrypt success")

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
		conn.SetCrypt(svr.weicLoginCrypt)
		go svr.loginWeic(conn)
	}
}

func (svr *Service) loginWeic(conn net.Conn) {
	slog.Debugf("login weic")
	// auth
	rawMsg, err := net.ReadMsg(conn)
	if err != nil {
		conn.Close()
		slog.Debugf("login weic read msg err:%v", err)
		return
	}
	loginReq, ok := rawMsg.(*msg.LoginReq)
	if !ok {
		conn.Close()
		slog.Debugf("login weic read msg no loginReq")
		return
	}
	slog.Debugf("login weic loginReq:%v", loginReq)
}
