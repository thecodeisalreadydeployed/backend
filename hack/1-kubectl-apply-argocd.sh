cd $(dirname $0)
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/thecodeisalreadydeployed/argo-cd/v2.0.5/manifests/install.yaml
kubectl apply -f .argocd.yaml
