FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mediscan ./cmd/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mediscan .

CMD ["./mediscan"]