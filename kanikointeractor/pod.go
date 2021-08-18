package kanikointeractor

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	workingDirectory := "__w"
	workingDirectoryVolumeMount := apiv1.VolumeMount{
		MountPath: fmt.Sprintf("/%s", workingDirectory),
		Name:      workingDirectory,
	}

	podSpec := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko",
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
			},
			InitContainers: []apiv1.Container{
				{
					Name:         "busybox",
					Image:        "busybox:1.33.1",
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount},
				},
				{
					Name:         "git",
					Image:        "alpine/git:v2.30.2",
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount},
					Command:      []string{"clone", it.BuildContext},
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
