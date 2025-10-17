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

package client

import (
	"context"
	"errors"
	"time"

	"github.com/gookit/slog"

	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/env"
	"github.com/gucooing/weiwei/pkg/msg"
	"github.com/gucooing/weiwei/pkg/net"
	"github.com/gucooing/weiwei/pkg/util/backoff"
	"github.com/gucooing/weiwei/pkg/util/crypt"
)

type Service struct {
	// ctx
	ctx context.Context
	// cancel
	cancel context.CancelFunc
	// control
	control *Control
}

func NewService() (*Service, error) {
	slog.Debugf("new client service...")
	s := &Service{}

	slog.Debugf("new client service success")
	return s, nil
}

func (svr *Service) Run(ctx context.Context) error {
	slog.Infof("client service run")
	ctx, cancel := context.WithCancel(ctx)
	svr.ctx = ctx
	svr.cancel = cancel

	// login weis
	svr.cycleLoginWeis()
	if svr.control == nil {
		return errors.New("weic login weis error")
	}
	go svr.keepController()

	<-svr.ctx.Done()
	// service context
	svr.Close()

	return nil
}

func (svr *Service) Close() {
	slog.Debugf("client service close...")

	slog.Debugf("client service close success")
}

func (svr *Service) cycleLoginWeis() {
	backoff.BackoffStart(
		func() bool {
			if err := svr.loginWeis(); err != nil {
				slog.Errorf("login weis err: %v", err)
				return false
			}
			return true
		},
		svr.ctx.Done(),
		5*time.Second,
	)
}

func (svr *Service) loginWeis() error {
	slog.Debugf("new weisConn...")
	ctl, err := NewControl(svr.ctx)
	if err != nil {
		return err
	}
	wsc, err := net.Dial(config.Client.WeisNet.Network, config.Client.WeisNet.Address)
	if err != nil {
		return err
	}
	slog.Debugf("network:%s address:%s new weisConn success", config.Client.WeisNet.Network, config.Client.WeisNet.Address)
	ctl.conn = wsc
	ctl.conn.SetCrypt(config.Client.WeicLogin.Crypt)

	// login
	loginReq := &msg.LoginReq{
		Version:   env.Version,
		Timestamp: time.Now().UnixNano(),
		LoginKey:  "",
	}
	config.Client.Auth.Verifier.SetVerifyLogin(loginReq)
	slog.Debugf("token:%s start login...", loginReq.LoginKey)
	_, err = msg.WriteMsg(ctl.conn, loginReq)
	if err != nil {
		return err
	}
	rawMsg, err := msg.ReadMsg(ctl.conn)
	if err != nil {
		return err
	}
	loginRsp, ok := rawMsg.(*msg.LoginRsp)
	if !ok {
		return errors.New("login weis read msg no loginRsp")
	}
	cry, err := crypt.NewCrypt(crypt.CryptTypeXor, loginRsp.Seed)
	if err != nil {
		return err
	}
	defer ctl.conn.SetCrypt(cry)

	slog.Debugf("loginRsp version:%s runId:%v seed:%v",
		loginRsp.Version, loginRsp.RunId, loginRsp.Seed)

	ctl.runId = loginRsp.RunId
	svr.control = ctl

	go ctl.Run()
	return nil
}

func (svr *Service) keepController() {
	for {
		select {
		case <-svr.ctx.Done():
			return
		case <-svr.control.doneChan:
			svr.cycleLoginWeis()
		}
	}
}
