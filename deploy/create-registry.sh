cd "$(dirname "$0")" || exit

registryName='kind-registry'
registryPort='5001'

# Create registry container unless it already exists
if [ "$(docker inspect -f '{{.State.Running}}' "${registryName}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${registryPort}:5000" --name "${registryName}" \
    registry:2
fi
