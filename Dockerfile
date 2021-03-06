FROM golang:1.16 as build-env
WORKDIR /__w
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build main.go

FROM gcr.io/distroless/base
WORKDIR /__w
COPY --from=build-env /__w/main /__w
EXPOSE 3000
ENV APP_ENV=PROD
CMD ["/__w/main"]
