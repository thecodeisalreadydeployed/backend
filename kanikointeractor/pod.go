package kanikointeractor

import (
	"fmt"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const busyboxImage = "busybox:1.33.1"

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	workingDirectoryVolumeMount := apiv1.VolumeMount{
		MountPath: config.DefaultKanikoWorkingDirectory,
		Name:      "working-directory",
	}

	dotSSHVolumeMount := apiv1.VolumeMount{
		MountPath: "/root/.ssh",
		Name:      "dot-ssh",
	}

	podLabel := map[string]string{
		"codedeploy/component": "kaniko",
	}
	defaultPodLabel := model.PodLabel(it.DeploymentID)
	err := mergo.Merge(&podLabel, defaultPodLabel)
	if err != nil {
		panic(err)
	}

	buildScript, err := PresetNestJS(BuildOptions{
		InstallCommand:   "yarn install",
		BuildCommand:     "yarn build a",
		WorkingDirectory: "nx/",
		OutputDirectory:  "dist/apps/a",
		StartCommand:     "node dist/apps/a/main",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("buildScript: %v\n", buildScript)

	buildScriptPath := filepath.Join(workingDirectoryVolumeMount.MountPath, "codedeploy.Dockerfile")

	fmt.Printf("buildScriptPath: %v\n", buildScriptPath)

	podSpec := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("kaniko-%s", util.RandomString(5)),
			Labels: podLabel,
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: workingDirectoryVolumeMount.Name,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: dotSSHVolumeMount.Name,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
			InitContainers: []apiv1.Container{
				{
					Name:         "init-busybox",
					Image:        busyboxImage,
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount, dotSSHVolumeMount},
					Command: []string{
						"sh",
						"-c",
						"echo \"" + buildScript + "\" > " + buildScriptPath,
					},
				},
				{
					Name:         "init-git",
					Image:        "alpine/git:v2.30.2",
					Args:         []string{"clone", "--single-branch", "--", it.BuildContext, filepath.Join(workingDirectoryVolumeMount.MountPath, "code")},
					VolumeMounts: []apiv1.VolumeMount{workingDirectoryVolumeMount},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:v1.6.0",
					Args: []string{
						fmt.Sprintf("--dockerfile=%s", filepath.Join(workingDirectoryVolumeMount.MountPath, "codedeploy.Dockerfile")),
						fmt.Sprintf("--context=dir://%s", filepath.Join(workingDirectoryVolumeMount.MountPath, "code")),
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

	kanikoSecretVolumeMount := apiv1.VolumeMount{
		MountPath: "/kaniko",
		Name:      "kaniko-secret",
	}

	podSpec.Spec.Volumes = append(podSpec.Spec.Volumes, apiv1.Volume{
		Name: kanikoSecretVolumeMount.Name,
		VolumeSource: apiv1.VolumeSource{
			EmptyDir: &apiv1.EmptyDirVolumeSource{},
		},
	})

	podSpec.Spec.Containers[0].Env = append(podSpec.Spec.Containers[0].Env, apiv1.EnvVar{
		Name:  "GOOGLE_APPLICATION_CREDENTIALS",
		Value: "/kaniko/config.json",
	})

	podSpec.Spec.Containers[0].VolumeMounts = append(podSpec.Spec.Containers[0].VolumeMounts, kanikoSecretVolumeMount)

	fmt.Printf("it.Registry.Secret(): %v\n", it.Registry.Secret())

	podSpec.Spec.InitContainers = append(podSpec.Spec.InitContainers, apiv1.Container{
		Name:         "init-gcr-secret",
		Image:        busyboxImage,
		VolumeMounts: []apiv1.VolumeMount{kanikoSecretVolumeMount},
		Command: []string{
			"sh",
			"-c",
			"echo \"" + it.Registry.Secret() + "\" > " + "/kaniko/config.json",
		},
	})

	return podSpec
}

func (it *KanikoInteractor) ECRKanikoPodSpec() apiv1.Pod {
	podSpec := it.baseKanikoPodSpec()
	return podSpec
}
