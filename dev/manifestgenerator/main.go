package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/thecodeisalreadydeployed/util"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func main() {
	dpl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-nginx-controller",
			Namespace: "ingress-nginx",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "ingress-nginx",
					"app.kubernetes.io/instance":  "ingress-nginx",
					"app.kubernetes.io/component": "controller",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":      "ingress-nginx",
						"app.kubernetes.io/instance":  "ingress-nginx",
						"app.kubernetes.io/component": "controller",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "controller",
							Image:           "k8s.gcr.io/ingress-nginx/controller:v1.0.0",
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Args:            []string{"/nginx-ingress-controller", "--controller-class=k8s.io/ingress-nginx"},
							SecurityContext: &apiv1.SecurityContext{
								Capabilities: &apiv1.Capabilities{
									Drop: []apiv1.Capability{"ALL"},
									Add:  []apiv1.Capability{"NET_BIND_SERVICE"},
								},
								RunAsUser:                util.Int64Ptr(101),
								AllowPrivilegeEscalation: util.BoolPtr(true),
							},
							Env: []apiv1.EnvVar{
								{
									Name: "POD_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name:  "LD_PRELOAD",
									Value: "/usr/local/lib/libmimalloc.so",
								},
							},
							LivenessProbe: &apiv1.Probe{
								FailureThreshold: 5,
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt(10254),
										Scheme: apiv1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
						},
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(dpl)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(y))

	var spec appsv1.Deployment
	_ = yaml.Unmarshal(y, &spec)
	spew.Dump(spec)
}
