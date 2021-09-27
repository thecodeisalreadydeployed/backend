cd "$(dirname "$0")" || exit
cd .. || exit
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39
golangci-lint run