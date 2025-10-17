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

	"github.com/gucooing/weiwei/client"
	"github.com/gucooing/weiwei/pkg/config"
	"github.com/gucooing/weiwei/pkg/env"
	_ "github.com/gucooing/weiwei/pkg/msg"
)

func init() {
	weicCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "weic.json", "config file")
	weicCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version")
}

var (
	cfgFile     string
	showVersion bool

	weicCmd = &cobra.Command{
		Use:                        env.WeiS,
		Short:                      fmt.Sprintf("%s is %s's client (%s)", env.WeiC, env.Name, env.Git),
		Example:                    fmt.Sprintf("%s --help", env.WeiC),
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
		RunE:                       initWeic,
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
	if err := weicCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initWeic(cmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Println(env.Version)
	}
	// load config
	if err := config.LoadClientConfig(cfgFile); err != nil {
		fmt.Println(err)
		return err
	}
	// run
	if err := runWeic(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func runWeic() error {
	slog.Configure(func(l *slog.SugaredLogger) {
		f := l.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		l.ChannelName = env.WeiS
		l.Level = config.Client.Log.Level
	})

	slog.Infof("%s is starting...", env.WeiC)

	svr, err := client.NewService()
	if err != nil {
		slog.Panicf("client.NewService err: %v", err)
		return err
	}
	slog.Infof("%s is startup success", env.WeiC)

	ctx := context.Background()

	return svr.Run(ctx)
}
