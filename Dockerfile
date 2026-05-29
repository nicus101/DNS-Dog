FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0

RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/dyn-dns .

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /app/dyn-dns /app/dyn-dns
ENTRYPOINT [ "/app/dyn-dns", "-w" ]
