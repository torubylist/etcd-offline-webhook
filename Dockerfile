FROM golang:1.10.4-alpine as builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /etcd-offline-webhook

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/etcd-offline-webhook

# Runtime image
FROM scratch AS base
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/etcd-offline-webhook /bin/etcd-offline-webhook