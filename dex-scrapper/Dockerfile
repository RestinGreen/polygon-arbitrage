FROM golang:latest

WORKDIR /app
COPY . .
RUN go build -o dex-scrapper cmd/main/main.go

CMD ["./dex-scrarpper"]