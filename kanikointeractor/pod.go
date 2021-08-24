package kanikointeractor

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	workingDirectory := "working-directory"
	workingDirectoryVolumeMount := apiv1.VolumeMount{
		MountPath: config.DefaultKanikoWorkingDirectory,
		Name:      workingDirectory,
	}

	dotSSH := "ssh"

	podLabel := map[string]string{
		"codedeploy/component": "kaniko",
	}
	defaultPodLabel := model.PodLabel(it.DeploymentID)
	err := mergo.Merge(&podLabel, defaultPodLabel)
	if err != nil {
		panic(err)
	}

	podSpec := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("kaniko-%s", util.RandomString(5)),
			Labels: podLabel,
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: workingDirectory,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: dotSSH,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
			InitContainers: []apiv1.Container{
				{
					Name:  "busybox",
					Image: "busybox:1.33.1",
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount, {
						MountPath: fmt.Sprintf("/%s", dotSSH),
						Name:      dotSSH,
					}},
				},
				{
					Name:  "git",
					Image: "alpine/git:v2.30.2",
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount, {
						MountPath: "/root/.ssh",
						Name:      dotSSH,
					}},
					Command: []string{"clone", it.BuildContext},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:v1.6.0",
					Args: []string{
						fmt.Sprintf("--dockerfile=%s", "codedeploy.Dockerfile"),
						fmt.Sprintf("--context=dir://%s", workingDirectory),
						fmt.Sprintf("--destination=%s", it.Destination),
					},
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount},
				},
			},
		},
	}
	return podSpec
}

func (it *KanikoInteractor) GCRKanikoPodSpec() apiv1.Pod {
	podSpec := it.baseKanikoPodSpec()
	return podSpec
}

func (it *KanikoInteractor) ECRKanikoPodSpec() apiv1.Pod {
	podSpec := it.baseKanikoPodSpec()
	return podSpec
}
