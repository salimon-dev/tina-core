# Build Stage
FROM golang:alpine as builder
WORKDIR /app
COPY . /app
RUN apk add --no-cache ca-certificates 
RUN go build -o bootstrap .

# Final Stage (Scratch)
FROM scratch
WORKDIR /app
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bootstrap ./bootstrap
ENTRYPOINT ["./bootstrap"]
