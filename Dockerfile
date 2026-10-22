FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/scheduler ./scheduler

COPY --from=builder /app/web ./web

ENV TODO_PORT=7540

ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

RUN mkdir -p /data

CMD ["./scheduler"]