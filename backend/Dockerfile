FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM golang:1.21 

WORKDIR /app

# Copy only the necessary artifacts from the builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/go.mod /app/go.mod
COPY --from=builder /app/go.sum /app/go.sum

EXPOSE 8080

CMD ["/app/main"]