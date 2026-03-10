FROM golang:1.21.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /grid-simulation .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /grid-simulation /app/grid-simulation

ENTRYPOINT ["/app/grid-simulation"]
