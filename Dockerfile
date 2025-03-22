FROM golang:alpine as builder
WORKDIR /app
COPY . /app
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap .

FROM scratch
WORKDIR /app
COPY --from=builder /app/bootstrap ./bootstrap
ENTRYPOINT [ "./bootstrap" ]