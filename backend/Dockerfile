FROM golang:1.27-alpine as builder
WORKDIR /src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app ./app
CMD ["./app"]
