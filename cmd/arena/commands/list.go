// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"io"

	"github.com/kubeflow/arena/cmd/arena/commands/flags"
	"github.com/kubeflow/arena/pkg/client"
	"github.com/kubeflow/arena/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

func NewListCommand() *cobra.Command {
	var allNamespaces bool
	var command = &cobra.Command{
		Use:   "list",
		Short: "list all the training jobs",
		Run: func(cmd *cobra.Command, args []string) {
			kubeClient, err := client.GetClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			client := kubeClient.GetClientset()

			if err != nil {
				log.Errorf("Failed due to %v", err)
				os.Exit(1)
			}

			namespace := flags.GetProjectFlagIncludingAll(cmd, kubeClient, allNamespaces)

			jobs := []TrainingJob{}
			trainers := NewTrainers(client)
			for _, trainer := range trainers {
				trainingJobs, err := trainer.ListTrainingJobs(namespace)
				if err != nil {
					log.Errorf("Failed due to %v", err)
					os.Exit(1)
				}
				jobs = append(jobs, trainingJobs...)
			}

			jobs = makeTrainingJobOrderdByAge(jobs)

			displayTrainingJobList(jobs, false)
		},
	}

	command.Flags().BoolVarP(&allNamespaces, "all-projects", "A", false, "list from all projects")

	return command
}

func displayTrainingJobList(jobInfoList []TrainingJob, displayGPU bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	labelField := []string{"NAME", "STATUS", "AGE", "NODE", "IMAGE", "INTERACTIVE", "PROJECT", "USER", "CREATED BY CLI", "SERVICE URL(S)"}

	PrintLine(w, labelField...)

	for _, jobInfo := range jobInfoList {
		status := GetJobRealStatus(jobInfo)
		hostIP := jobInfo.HostIPOfChief()
		PrintLine(w, jobInfo.Name(),
			status,
			util.ShortHumanDuration(jobInfo.Age()),
			hostIP, jobInfo.Image(), jobInfo.Interactive(), jobInfo.Project(), jobInfo.User(), strconv.FormatBool(jobInfo.CreatedByCLI()), strings.Join(jobInfo.ServiceURLs(), ", "))
	}
	_ = w.Flush()
}

func PrintLine(w io.Writer, fields ...string) {
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	buffer := strings.Join(fields, "\t")
	fmt.Fprintln(w, buffer)
}
