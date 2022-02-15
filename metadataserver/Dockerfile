FROM golang:1.16 as build-env
WORKDIR /__w
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o metadataserver -ldflags '-w -s' main.go

FROM scratch
WORKDIR /__w
COPY --from=build-env /__w/metadataserver /__w
CMD /__w/metadataserver
