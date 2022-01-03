cd "$(dirname "$0")" || exit
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.2.2/manifests/core-install.yaml
kubectl apply -f .argocd.yaml
