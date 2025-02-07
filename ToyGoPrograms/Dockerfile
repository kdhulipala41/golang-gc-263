FROM golang:1.23.5

WORKDIR /app

COPY . .

RUN go build -o simplealloc simplealloc.go

CMD ["./simplealloc"]