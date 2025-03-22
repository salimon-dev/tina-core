FROM golang:alpine as builder
WORKDIR /app
COPY . /app

ARG GOOS=linux
ARG GOARCH=amd64

RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o bootstrap .

FROM scratch
WORKDIR /app
COPY --from=builder /app/bootstrap ./bootstrap
ENTRYPOINT [ "./bootstrap" ]