FROM golang:1.23-alpine AS builder

WORKDIR /dyn-dns

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /dyn-dns .

FROM scratch

WORKDIR /app
COPY --from=builder /dyn-dns /app/dyn-dns
ENTRYPOINT [ "/app/dyn-dns" ]
