// Copyright Calvin Yang.
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

package command

import (
	"github.com/spf13/cobra"
)

func NewClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "show the cluster information",
	}
	cmd.AddCommand(newCheckClusterStatus())
	return cmd
}

func newCheckClusterStatus() *cobra.Command {
	r := &cobra.Command{
		Use:   "check",
		Short: "check the cluster status",
		Run:   checkClusterStatusCommandFunc,
	}
	return r
}

func checkClusterStatusCommandFunc(cmd *cobra.Command, args []string) {
	cmd.Println("checkClusterStatusCommandFunc")
}
