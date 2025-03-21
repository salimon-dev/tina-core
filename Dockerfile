FROM golang:alpine
WORKDIR /app
COPY . /app
RUN go build -o bootstrap .
ENTRYPOINT ["./bootstrap"]
