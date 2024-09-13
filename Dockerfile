# the first parse: build the application
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o canoe canoe/cmd/canoe
RUN ls -la

# the second parse: deploy the application
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/.env .
COPY --from=builder /app/.env.dev .
COPY --from=builder /app/canoe .
EXPOSE 9000
CMD ["./canoe"]