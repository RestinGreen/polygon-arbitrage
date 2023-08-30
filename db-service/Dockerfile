FROM golang:latest

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o db-service cmd/main/main.go


CMD ["./db-service"]