FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o temperature-api main.go

FROM alpine:3.20

WORKDIR /root/
COPY --from=builder /app/temperature-api .

EXPOSE 8081
CMD ["./temperature-api"]