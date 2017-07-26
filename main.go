package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type LsofPlugin struct {
	Prefix string
}

func (l LsofPlugin) MetricKeyPrefx() string {
	if l.Prefix == "" {
		l.Prefix = "lsof.Counts"
	}
	return l.Prefix
}

func (l LsofPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(l.Prefix)
	return map[string](mp.Graphs){
		l.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "counts.process.java", Label: "Java"},
			},
		},
	}
}

func (l LsofPlugin) FetchMetrics() (map[string]interface{}, error) {
	ut, err := lsofGet()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch lsof metrics: %s", err)
	}

	return map[string]interface{}{"counts.process.java": ut}, nil
}

func lsofGet() (uint64, error) {
	cmdstr := "lsof -p $(pgrep java) | wc -l"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	str := strings.TrimSpace(string(out))
	i, _ := strconv.ParseUint(str, 10, 32)
	return i, err
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "lsof", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := LsofPlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(u)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
