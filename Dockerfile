FROM golang:1.16 as build-env
WORKDIR /__w
ADD . /__w
RUN go build main.go

FROM gcr.io/distroless/base
WORKDIR /__w
COPY --from=build-env /__w/main /__w
EXPOSE 3000
CMD ["/__w/main"]