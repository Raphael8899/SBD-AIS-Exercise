# build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /app/ordersystem

# runtime stage  
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /app/frontend /app/frontend
EXPOSE 3000
CMD ["/app/ordersystem"]