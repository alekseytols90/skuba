/*
 * Copyright (c) 2019 SUSE LLC. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package node

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/SUSE/caaspctl/internal/pkg/caaspctl/deployments"
	"github.com/SUSE/caaspctl/internal/pkg/caaspctl/deployments/ssh"
	node "github.com/SUSE/caaspctl/pkg/caaspctl/actions/node/join"
)

type joinOptions struct {
	role                  string
	ignorePreflightErrors string
}

func NewJoinCmd() *cobra.Command {
	joinOptions := joinOptions{}
	target := ssh.Target{}

	cmd := &cobra.Command{
		Use:   "join <node-name>",
		Short: "Joins a new node to the cluster",
		Run: func(cmd *cobra.Command, nodenames []string) {
			joinConfiguration := deployments.JoinConfiguration{
				KubeadmExtraArgs: map[string]string{"ignore-preflight-errors": joinOptions.ignorePreflightErrors},
			}

			joinConfiguration.Role = deployments.MustGetRoleFromString(joinOptions.role)

			if err := node.Join(joinConfiguration, target.GetDeployment(nodenames[0])); err != nil {
				klog.Fatalf("error joining node %s: %s", nodenames[0], err)
			}
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().AddFlagSet(target.GetFlags())
	cmd.Flags().StringVarP(&joinOptions.role, "role", "r", "", "Role that this node will have in the cluster (master|worker)")
	cmd.Flags().StringVar(&joinOptions.ignorePreflightErrors, "ignore-preflight-errors", "", "Comma separated list of preflight errors to ignore")

	cmd.MarkFlagRequired("role")

	return cmd
}
