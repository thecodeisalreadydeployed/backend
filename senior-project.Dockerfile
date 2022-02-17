FROM golang:1.16 as build-env
WORKDIR /__w
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags '-w -s' main.go

FROM scratch
WORKDIR /__w
COPY --from=build-env /__w/main /__w
EXPOSE 3000
ENV APP_ENV=PROD
ENTRYPOINT /__w/main
