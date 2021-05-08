package cron

import (
	"encoding/json"
	"fmt"
	"github.com/kubeflow/arena/pkg/apis/types"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var getCronTemplate = `
Name:               %v
Namespace:          %v
Type:               %v
Schedule:           %v
Suspend:            %v
ConcurrencyPolicy:  %v
CreationTimestamp:  %v
LastScheduleTime:   %v
Deadline:           %v
%v
`

func GetCronInfo(name, namespace string) (*types.CronInfo, error) {
	return GetCronHandler().GetCron(namespace, name)
}

func DisplayCron(cron *types.CronInfo, format types.FormatStyle) {
	switch format {
	case "json":
		data, _ := json.MarshalIndent(cron, "", "    ")
		fmt.Printf("%v", string(data))
		return
	case "yaml":
		data, _ := yaml.Marshal(cron)
		fmt.Printf("%v", string(data))
		return
	case "", "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		lines := []string{"\nHistory:", "NAME\tSTATUS\tTYPE\tCREATETIME\tFINISHTIME"}
		lines = append(lines, "----\t------\t----\t----------\t----------")

		if len(cron.History) > 0 {
			for _, item := range cron.History {
				lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
					item.Name, item.Status, item.Kind, item.CreateTime, item.FinishTime))
			}
		}

		printLine(w, fmt.Sprintf(strings.Trim(getCronTemplate, "\n"),
			cron.Name,
			cron.Namespace,
			cron.Type,
			cron.Schedule,
			strconv.FormatBool(cron.Suspend),
			cron.ConcurrencyPolicy,
			cron.CreationTimestamp,
			cron.LastScheduleTime,
			cron.Deadline,
			strings.Join(lines, "\n"),
		))

		_ = w.Flush()
		return
	}
}
