cd "$(dirname "$0")" || exit
kubectl create namespace argocd
kubectl apply -k argocd
kubectl apply -n argocd -f argocd/app.yml
