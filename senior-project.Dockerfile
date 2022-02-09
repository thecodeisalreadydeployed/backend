FROM golang:1.16 as build-env
WORKDIR /__w
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build main.go

FROM golang:1.16-alpine
WORKDIR /__w
ADD https://github.com/trif0lium/secrets-resolve/release/download/v0.0.3/secrets-resolve_0.0.3_linux_amd64 /__w/secrets-resolve
RUN chmod +x /__w/secrets-resolve
ADD https://github.com/krallin/tini/releases/download/v0.19.0/tini-static /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]
COPY --from=build-env /__w/main /__w
EXPOSE 3000
ENV APP_ENV=PROD
CMD /__w/secrets-resolve && /__w/main
