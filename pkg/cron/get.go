package cron

import (
	"encoding/json"
	"fmt"
	"github.com/kubeflow/arena/pkg/apis/config"
	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"os"
	"strconv"
	"text/tabwriter"
)

func GetCronInfo(name, namespace string) (*types.CronInfo, error) {
	config := config.GetArenaConfiger().GetRestConfig()

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ret, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	b, err := ret.MarshalJSON()
	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(b)

	creationTimestamp := r.Get("metadata").Get("creationTimestamp").String()
	createTime, err := formatTime(creationTimestamp)

	c := &types.CronInfo{
		Name: r.Get("metadata").Get("name").String(),
		Namespace: r.Get("metadata").Get("namespace").String(),
		Type: r.Get("spec").Get("template").Get("kind").String(),
		Schedule: r.Get("spec").Get("schedule").String(),
		ConcurrencyPolicy: r.Get("spec").Get("concurrencyPolicy").String(),
		HistoryLimit: r.Get("spec").Get("historyLimit").Int(),
		Deadline: r.Get("spec").Get("deadline").String(),
		Suspend: r.Get("spec").Get("suspend").Bool(),
		CreationTimestamp: createTime.Unix(),
	}

	return c, nil
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
		var header []string
		header = append(header, []string{"NAMESPACE", "NAME", "TYPE", "SCHEDULE", "SUSPEND", "DEADLINE", "CONCURRENCYPOLICY"}...)
		PrintLine(w, header...)

		var items []string
		items = append(items, []string{
			cron.Namespace,
			cron.Name,
			cron.Type,
			cron.Schedule,
			strconv.FormatBool(cron.Suspend),
			cron.Deadline,
			cron.ConcurrencyPolicy,
		}...)
		PrintLine(w, items...)

		PrintLine(w, "\nHistory")

		_ = w.Flush()
		return
	}
}
