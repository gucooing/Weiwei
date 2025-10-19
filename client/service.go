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
	err := backoff.BackoffStart(
		func() error {
			if err := svr.loginWeis(); err != nil {
				slog.Debugf("login weis err: %v", err)
				return err
			}
			return nil
		},
		svr.ctx.Done(),
		&backoff.ExponentialBackoff{
			BaseInterval: 5 * time.Second,
			MaxRetries:   0,
			MaxInterval:  30 * time.Second,
		},
	)
	if err != nil {
		slog.Errorf("login weis err: %v", err)
	}
}

func (svr *Service) loginWeis() error {
	slog.Debugf("new weisConn...")
	conn, err := net.Dial(config.Client.WeisNet.Network, config.Client.WeisNet.Address)
	if err != nil {
		return err
	}
	conn.SetCrypt(config.Client.WeicLogin.Crypt)
	slog.Debugf("network:%s address:%s new weisConn success", config.Client.WeisNet.Network, config.Client.WeisNet.Address)

	// login
	timestamp := time.Now().UnixNano()
	loginReq := &msg.CSLoginReq{
		Version:   env.Version,
		Timestamp: timestamp,
		LoginKey:  config.Client.Auth.Verifier.SetVerifyLogin(timestamp),
	}

	slog.Debugf("token:%s start login...", loginReq.LoginKey)
	_, err = msg.WriteMsg(conn, loginReq)
	if err != nil {
		return err
	}
	rawMsg, err := msg.ReadMsg(conn)
	if err != nil {
		return err
	}
	loginRsp, ok := rawMsg.(*msg.SCLoginRsp)
	if !ok {
		return errors.New("login weis read msg no loginRsp")
	}
	cry, err := crypt.NewCrypt(crypt.CryptTypeXor, loginRsp.Seed)
	if err != nil {
		return err
	}
	defer conn.SetCrypt(cry)

	ctl, err := NewControl(conn)
	if err != nil {
		return err
	}
	slog.Debugf("loginRsp version:%s runId:%v seed:%v",
		loginRsp.Version, loginRsp.RunId, loginRsp.Seed)

	ctl.runId = loginRsp.RunId
	ctl.seed = loginRsp.Seed
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
