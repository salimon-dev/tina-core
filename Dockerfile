FROM golang:alpine as builder
WORKDIR /app
COPY . /app

ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

RUN go build -o bootstrap .

FROM scratch
WORKDIR /app
COPY --from=builder /app/bootstrap ./bootstrap
ENTRYPOINT [ "./bootstrap" ]