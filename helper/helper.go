// Copyright 2016 TiKV Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"strings"
	"tihelper/helper/command"
)

func init() {
	cobra.EnablePrefixMatching = true
}

// Modified by Calvin 20220429
func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "TiHelper",
		Short: "TiDB Operation Helper",
	}

	rootCmd.PersistentFlags().StringP("pd", "u", "http://127.0.0.1:2379", "address of pd")

	rootCmd.AddCommand(
		command.NewClusterCommand(),
	)
	rootCmd.Flags().ParseErrorsWhitelist.UnknownFlags = true
	rootCmd.SilenceErrors = true

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	return rootCmd
}

// Modified by Calvin 20220429
func MainStart(args []string) {
	rootCmd := GetRootCmd()
	rootCmd.Flags().BoolP("interact", "i", false, "Run tihelper with readline.")
	rootCmd.Flags().BoolP("version", "V", false, "Print version information and exit.")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if flag, err := cmd.Flags().GetBool("version"); err == nil && flag {
			fmt.Println("Release Version:", "V1.0.0")
			return
		}
		if flag, err := cmd.Flags().GetBool("interact"); err == nil && flag {
			readlineCompleter := readline.NewPrefixCompleter(genCompleter(cmd)...)
			loop(cmd.PersistentFlags(), readlineCompleter)
		}
	}
	rootCmd.SetArgs(args)
	rootCmd.ParseFlags(args)
	rootCmd.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func loop(persistentFlags *pflag.FlagSet, readlineCompleter readline.AutoCompleter) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       "/tmp/readline.tmp",
		AutoComplete:      readlineCompleter,
		InterruptPrompt:   "^C",
		EOFPrompt:         "^D",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	getREPLCmd := func() *cobra.Command {
		rootCmd := GetRootCmd()
		persistentFlags.VisitAll(func(flag *pflag.Flag) {
			if flag.Changed {
				rootCmd.PersistentFlags().Set(flag.Name, flag.Value.String())
			}
		})
		rootCmd.LocalFlags().MarkHidden("pd")
		rootCmd.SetOutput(os.Stdout)
		return rootCmd
	}

	for {
		line, err := l.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				break
			} else if err == io.EOF {
				break
			}
			continue
		}
		if line == "exit" {
			os.Exit(0)
		}
		args, err := shellwords.Parse(line)
		if err != nil {
			fmt.Printf("parse command err: %v\n", err)
			continue
		}

		rootCmd := getREPLCmd()
		rootCmd.SetArgs(args)
		rootCmd.ParseFlags(args)
		if err := rootCmd.Execute(); err != nil {
			rootCmd.Println(err)
		}
	}
}

func genCompleter(cmd *cobra.Command) []readline.PrefixCompleterInterface {
	pc := []readline.PrefixCompleterInterface{}

	for _, v := range cmd.Commands() {
		if v.HasFlags() {
			flagsPc := []readline.PrefixCompleterInterface{}
			flagUsages := strings.Split(strings.Trim(v.Flags().FlagUsages(), " "), "\n")
			for i := 0; i < len(flagUsages)-1; i++ {
				flagsPc = append(flagsPc, readline.PcItem(strings.Split(strings.Trim(flagUsages[i], " "), " ")[0]))
			}
			flagsPc = append(flagsPc, genCompleter(v)...)
			pc = append(pc, readline.PcItem(strings.Split(v.Use, " ")[0], flagsPc...))
		} else {
			pc = append(pc, readline.PcItem(strings.Split(v.Use, " ")[0], genCompleter(v)...))
		}
	}
	return pc
}
