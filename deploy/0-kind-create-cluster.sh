cd "$(dirname "$0")" || exit
kind create cluster --config ./.kind.config.yaml
