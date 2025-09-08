FROM golang:1.24.4-alpine AS builder

WORKDIR /scheduler

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /app

FROM alpine:latest

WORKDIR /scheduler

COPY --from=builder /app .
COPY --from=builder /scheduler/web ./web

ENV TODO_PORT=7540 \
    TODO_DBFILE=/data/scheduler.db \
    TODO_PASSWORD=12345

EXPOSE 7540

CMD ["./app"]