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

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gookit/slog"
	"github.com/spf13/cobra"

	_ "github.com/gucooing/weiwei/pkg/msg"

	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/env"
	"github.com/gucooing/weiwei/server"
)

func init() {
	weisCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "weis.json", "config file")
	weisCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version")
}

var (
	cfgFile     string
	showVersion bool

	weisCmd = &cobra.Command{
		Use:                        env.WeiS,
		Short:                      fmt.Sprintf("%s is %s's server (%s)", env.WeiS, env.Name, env.Git),
		Example:                    fmt.Sprintf("%s --help", env.WeiS),
		ValidArgs:                  nil,
		ValidArgsFunction:          nil,
		Args:                       nil,
		ArgAliases:                 nil,
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Annotations:                nil,
		Version:                    env.Version,
		PersistentPreRun:           nil,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        nil,
		RunE:                       initWeis,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
		CompletionOptions:          cobra.CompletionOptions{},
		TraverseChildren:           false,
		Hidden:                     false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}
)

func main() {
	if err := weisCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initWeis(cmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Println(env.Version)
	}
	// load config
	if err := config.LoadServerConfig(cfgFile); err != nil {
		fmt.Println(err)
		return err
	}
	// run
	if err := runWeis(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func runWeis() error {
	slog.Configure(func(l *slog.SugaredLogger) {
		f := l.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		l.ChannelName = env.WeiS
		l.Level = config.Server.Log.Level
	})

	// slog.Panicf("%s is starting...",env.WeiS)
	// slog.Fatalf("%s is starting...",env.WeiS)
	// slog.Errorf("%s is starting...",env.WeiS)
	// slog.Warnf("%s is starting...",env.WeiS)
	// slog.Noticef("%s is starting...",env.WeiS)
	slog.Infof("%s is starting...", env.WeiS)
	// slog.Debugf("%s is starting...",env.WeiS)
	// slog.Tracef("%s is starting...",env.WeiS)

	svr, err := server.NewService()
	if err != nil {
		slog.Panicf("server.NewService err: %v", err)
		return err
	}
	slog.Infof("%s is startup success", env.WeiS)

	ctx := context.Background()
	svr.Run(ctx)

	return nil
}
