FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o main cmd/todo/main.go
CMD ["./main"]