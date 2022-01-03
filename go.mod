module github.com/thecodeisalreadydeployed

go 1.16

require (
	bou.ke/monkey v1.0.2
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/ProtonMail/go-crypto v0.0.0-20210707164159-52430bf6b52c // indirect
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/argoproj/argo-cd/v2 v2.2.2
	github.com/bxcodec/faker/v3 v3.6.0
	github.com/fasthttp/websocket v1.4.3 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gavv/httpexpect/v2 v2.3.1
	github.com/ghodss/yaml v1.0.0 //ct
	github.com/go-git/go-billy/v5 v5.3.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gofiber/fiber/v2 v2.16.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moby/buildkit v0.9.1
	github.com/nginxinc/kubernetes-ingress v1.12.0
	github.com/segmentio/ksuid v1.0.4
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/viper v1.10.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/subosito/gotenv v1.2.0
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/automaxprocs v1.4.0
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.13
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/kustomize/api v0.10.0
	sigs.k8s.io/kustomize/kyaml v0.12.0
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/argoproj/gitops-engine => github.com/argoproj/gitops-engine v0.4.0
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
	k8s.io/apiserver => k8s.io/apiserver v0.21.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.4
	k8s.io/client-go => k8s.io/client-go v0.21.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.4
	k8s.io/code-generator => k8s.io/code-generator v0.21.4
	k8s.io/component-base => k8s.io/component-base v0.21.4
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.4
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.4
	k8s.io/cri-api => k8s.io/cri-api v0.21.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.4
	k8s.io/kubectl => k8s.io/kubectl v0.21.4
	k8s.io/kubelet => k8s.io/kubelet v0.21.4
	k8s.io/kubernetes => k8s.io/kubernetes v1.21.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.4
	k8s.io/metrics => k8s.io/metrics v0.21.4
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.4
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.4
)
