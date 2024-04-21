# Stage 1
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN mkdir /build

ADD . /build/
WORKDIR /build

RUN go mod download
RUN go build -o main .

# Stage 2
FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app
COPY --from=builder /build/main /app/
WORKDIR /app

EXPOSE 6767

CMD ["./main"]