package main

import (
	"strings"
	
	"example.com/hello/imports/actionssummerwinddev"
	// "example.com/hello/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

const (
	serviceName             = "pipeline"
	githubEnterprise        = "skandinaviska-enskilda-banken-ab"
	githubRunnerGroup       = "gaas"
	runnerSetServiceAccount = "arc-runnerset"
	runnerImage             = "tfs-custom-runners-docker.repo7.sebank,se/arc-runner-dind"

	runnerVolumeMountDockerNameBase     = "var-lib-docker"
	runnerVolumeMountDockerPath         = "/var/lib/docker"
	runnerVolumeMountDockerStorageClass = "arc-lib-docker"
	runnerVolumeMountToolsNameBase      = "runner-tool-cache"
	runnerVolumeMountToolsPath          = "/opt/hostedtoolcache"
	runnerVolumeMountToolsStorageClass  = "arc-runner-tool-cache"
)

// Maps and slices can't be constants
var (
	runnerEnv = map[string]string{
		"HTTP_PROXY":            "http://gias.sebank.se:8080",
		"HTTPS_PROXY":           "http://gias.sebank.se:8080",
		"NO_PROXY":              ".sebank.se, localhost, .seb.net,*.sebank.se,*.seb.net",
		"DISABLE_RUNNER_UPDATE": "false",
	}
	runnerResource = map[string]string{
		"cpu":    "2000m",
		"memory": "7Gi",
	}
)

func generateEnv(envMap map[string]string) *[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv {
	var k8sEnv []*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv
	for k, v := range envMap {
		k8sEnv = append(k8sEnv, &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv{
			Name:  jsii.String(k),
			Value: jsii.String(v),
		})
	}
	return &k8sEnv
}

func NewRunnerContainer(id, runnerImage string) *actionssummerwinddev.RunnerSetSpecTemplateSpecContainers {
	return &actionssummerwinddev.RunnerSetSpecTemplateSpecContainers{
		Name:  jsii.String(id),
		Image: jsii.String(runnerImage),
		SecurityContext: &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersSecurityContext{
			Privileged: jsii.Bool(true),
		},
		Env: generateEnv(runnerEnv),
		// Resources: &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResources{
		// 	Limits: &map[string]actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResourcesLimits{
		// 		"cpu": nil,
		// 		"memory": nil,
		// 	},
		// 	Requests: &map[string]actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResourcesRequests{
		// 		"cpu": nil,
		// 		"memory": nil,
		// 	},
		// },
	}
}

func NewRunnerSet(scope constructs.Construct, id *string, cached bool, prefix, postfix string) constructs.Construct {
	construct := constructs.NewConstruct(scope, id)

	runnerVolumeMountDockerName := prefix + runnerVolumeMountDockerNameBase + postfix
	runnerVolumeMountToolsName := prefix + runnerVolumeMountToolsNameBase + postfix

	var cacheString string
	if cached {
		cacheString = "cached"
	} else {
		cacheString = "cacheless"
	}
	labelGaas := strings.Join([]string{"gaas", cacheString, "v1"}, "-")
	labelRunnerset := prefix + cacheString + "-runnerset" + postfix

	baseRunnerContainer := NewRunnerContainer("runner", runnerImage)

	volumeClaimTemplates := []*actionssummerwinddev.RunnerSetSpecVolumeClaimTemplates{}
	if cached {
		baseRunnerContainer.VolumeMounts = &[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersVolumeMounts{
			{Name: jsii.String(runnerVolumeMountDockerName), MountPath: jsii.String(runnerVolumeMountDockerPath)},
			{Name: jsii.String(runnerVolumeMountToolsName), MountPath: jsii.String(runnerVolumeMountToolsPath)},
		}
		volumeClaimTemplates = []*actionssummerwinddev.RunnerSetSpecVolumeClaimTemplates{
			{
				Metadata: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesMetadata{
					Name: jsii.String(runnerVolumeMountDockerName),
				},
				Spec: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpec{
					AccessModes:      jsii.Strings("ReadWriteOnce"),
					StorageClassName: jsii.String(runnerVolumeMountDockerStorageClass),
				},
			},
			{
				Metadata: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesMetadata{
					Name: jsii.String(runnerVolumeMountToolsName),
				},
				Spec: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpec{
					AccessModes:      jsii.Strings("ReadWriteOnce"),
					StorageClassName: jsii.String(runnerVolumeMountToolsStorageClass),
				},
			},
		}
	}
	spec := &actionssummerwinddev.RunnerSetSpec{
		Selector: &actionssummerwinddev.RunnerSetSpecSelector{
			MatchLabels: &map[string]*string{"app": id},
		},
		ServiceName:                  jsii.String(serviceName),
		DockerdWithinRunnerContainer: jsii.Bool(true),
		Enterprise:                   jsii.String(githubEnterprise),
		Ephemeral:                    jsii.Bool(true),
		Group:                        jsii.String(githubRunnerGroup),
		Labels: &[]*string{
			jsii.String(labelRunnerset), jsii.String(labelGaas),
		},
		Template: &actionssummerwinddev.RunnerSetSpecTemplate{
			Metadata: &actionssummerwinddev.RunnerSetSpecTemplateMetadata{Name: id},
			Spec: &actionssummerwinddev.RunnerSetSpecTemplateSpec{
				Containers: &[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainers{
					baseRunnerContainer,
				},
				ServiceAccountName: jsii.String(runnerSetServiceAccount),
			},
		},
		VolumeClaimTemplates: &volumeClaimTemplates,
	}

	actionssummerwinddev.NewRunnerSet(construct, id, &actionssummerwinddev.RunnerSetProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: id},
		Spec:     spec,
	})

	return construct
}

func NewChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	NewRunnerSet(chart, jsii.String("cached-runnerset-a"), true, "", "-a")
	NewRunnerSet(chart, jsii.String("cached-runnerset-b"), true, "", "-b")
	NewRunnerSet(chart, jsii.String("cacheless-runnerset"), false, "", "")

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewChart(app, "runners", nil)
	app.Synth()
}
