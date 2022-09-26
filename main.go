package main

import (
	"example.com/hello/imports/actionssummerwinddev"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	runnersetName := "runnerset-a"

	actionssummerwinddev.NewRunnerSet(chart, jsii.String(runnersetName), &actionssummerwinddev.RunnerSetProps{
		Spec: &actionssummerwinddev.RunnerSetSpec{
			Selector: &actionssummerwinddev.RunnerSetSpecSelector{
				MatchLabels: &map[string]*string{
					"app": jsii.String(runnersetName),
				},
			},
			ServiceName: jsii.String(runnersetName),
			Template:    &actionssummerwinddev.RunnerSetSpecTemplate{},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "hello", nil)
	app.Synth()
}
