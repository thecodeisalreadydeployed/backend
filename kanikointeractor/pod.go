package kanikointeractor

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	__w := "__w"
	podSpec := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko",
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: __w,
				},
			},
			InitContainers: []apiv1.Container{
				{
					Name:  "git",
					Image: "alpine/git:v2.30.2",
					VolumeMounts: []apiv1.VolumeMount{
						{
							MountPath: fmt.Sprintf("/%s", __w),
							Name:      __w,
						},
					},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:v1.6.0",
					Args: []string{
						fmt.Sprintf("--dockerfile=%s", "codedeploy.Dockerfile"),
						fmt.Sprintf("--context=%s", it.BuildContext),
						fmt.Sprintf("--destination=%s", it.Destination),
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							MountPath: fmt.Sprintf("/%s", __w),
							Name:      __w,
						},
					},
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
