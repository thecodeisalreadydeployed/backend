FROM golang:1.16 as build-env
WORKDIR /__w
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags '-w -s' main.go

FROM scratch
WORKDIR /__w
ADD https://github.com/trif0lium/secrets-resolve/releases/download/v0.0.3/secrets-resolve_0.0.3_linux_amd64 /__w/secrets-resolve
RUN chmod +x /__w/secrets-resolve
ADD https://github.com/krallin/tini/releases/download/v0.19.0/tini-static /tini
ENTRYPOINT ["/tini", "--"]
COPY --from=build-env /__w/main /__w
EXPOSE 3000
ENV APP_ENV=PROD
CMD /__w/secrets-resolve && /__w/main
