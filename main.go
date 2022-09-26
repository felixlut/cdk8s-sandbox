package main

import (
	"example.com/hello/imports/actionssummerwinddev"
	// "example.com/hello/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func generateEnv(envAsMap map[string]string) *[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv {
	var k8sEnv []*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv
	for k, v := range envAsMap {
		k8sEnv = append(k8sEnv, &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersEnv{
			Name:  jsii.String(k),
			Value: jsii.String(v),
		})
	}
	return &k8sEnv
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	runnersetName := "runnerset-a"
	dockerdWithinRunnerContainer := true
	ephemeral := true
	githubRunnerGroup := "gaas"
	githubEnterprise := "skandinaviska-enskilda-banken-ab"
	githubRunnerLabels := []string{"runnerset-a", "gaas", "dependabot"} // TODO: Fix actual labels
	serviceName := "pipeline"
	serviceAccountName := "arc-runnerset"
	runnerContainerName := "runner"
	runnerImage := "tfs-custom-runners-docker.repo7.sebank,se/arc-runner-dind" // TODO: Image tag
	runnerEnvMap := map[string]string{
		"HTTP_PROXY":            "http://gias.sebank.se:8080",
		"HTTP_PROXYS":           "http://gias.sebank.se:8080",
		"NO_PROXY":              ".sebank.se, localhost, .seb.net,*.sebank.se,*.seb.net",
		"DISABLE_RUNNER_UPDATE": "false",
	}
	runnerEnv := generateEnv(runnerEnvMap)

	// runnerCpu := "2000m"
	// runnerMemory := "7Gi"
	runnerPrivileged := true
	runnerDockerVolumeName := "var-lib-docker-a"
	runnerDockerVolumePath := "/var/lib/docker"
	runnerDockerVolumeStorageClass := "arc-lib-docker"
	runnerDockerVolumeAccess := "ReadWriteOnce"
	runnerToolsVolumeName := "runner-tool-cache-a"
	runnerToolsVolumePath := "/opt/hostedtoolcache"
	runnerToolsVolumeStorageClass := "arc-runner-tool-cache"
	runnerToolsVolumeAccess := "ReadWriteOnce"

	actionssummerwinddev.NewRunnerSet(chart, jsii.String(runnersetName), &actionssummerwinddev.RunnerSetProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: jsii.String(runnersetName),
		},
		Spec: &actionssummerwinddev.RunnerSetSpec{
			DockerdWithinRunnerContainer: jsii.Bool(dockerdWithinRunnerContainer),
			Ephemeral:                    jsii.Bool(ephemeral),
			Group:                        jsii.String(githubRunnerGroup),
			Enterprise:                   jsii.String(githubEnterprise),
			Labels:                       jsii.Strings(githubRunnerLabels...),
			Selector: &actionssummerwinddev.RunnerSetSpecSelector{
				MatchLabels: &map[string]*string{
					"app": jsii.String(runnersetName),
				},
			},
			ServiceName: jsii.String(serviceName),
			Template: &actionssummerwinddev.RunnerSetSpecTemplate{
				Metadata: &actionssummerwinddev.RunnerSetSpecTemplateMetadata{
					Labels: &map[string]*string{
						"app": jsii.String(runnersetName),
					},
				},
				Spec: &actionssummerwinddev.RunnerSetSpecTemplateSpec{
					ServiceAccountName: jsii.String(serviceAccountName),
					Containers: &[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainers{
						{
							Name:    jsii.String(runnerContainerName),
							Env:     runnerEnv,
							Image:   jsii.String(runnerImage),
							// TODO: Resources
							// Resources:                &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResources{
							// 	Limits:   &map[string]actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResourcesLimits{
							// 		"cpu":
							// 	},
							// 	Requests: &map[string]actionssummerwinddev.RunnerSetSpecTemplateSpecContainersResourcesRequests{},
							// },
							SecurityContext: &actionssummerwinddev.RunnerSetSpecTemplateSpecContainersSecurityContext{Privileged: jsii.Bool(runnerPrivileged)},
							VolumeMounts: &[]*actionssummerwinddev.RunnerSetSpecTemplateSpecContainersVolumeMounts{
								{Name: jsii.String(runnerDockerVolumeName), MountPath: jsii.String(runnerDockerVolumePath)},
								{Name: jsii.String(runnerToolsVolumeName), MountPath: jsii.String(runnerToolsVolumePath)},
							},					
						},
					},
				},
			},
			VolumeClaimTemplates: &[]*actionssummerwinddev.RunnerSetSpecVolumeClaimTemplates{
				{
					Metadata: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesMetadata{
						Name: jsii.String(runnerDockerVolumeName),
					},
					Spec: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpec{
						AccessModes: &[]*string{jsii.String(runnerDockerVolumeAccess)},
						// Resources:        &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpecResources{
						// 	Requests: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpecResourcesRequests{

						// 	},
						// },
						StorageClassName: jsii.String(runnerDockerVolumeStorageClass),
					},
				},
				{
					Metadata: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesMetadata{
						Name: jsii.String(runnerToolsVolumeName),
					},
					Spec: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpec{
						AccessModes: &[]*string{jsii.String(runnerToolsVolumeAccess)},
						// Resources:        &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpecResources{
						// 	Requests: &actionssummerwinddev.RunnerSetSpecVolumeClaimTemplatesSpecResourcesRequests{

						// 	},
						// },
						StorageClassName: jsii.String(runnerToolsVolumeStorageClass),
					},
				},
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "hello", nil)
	app.Synth()
}
