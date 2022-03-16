FROM golang:alpine AS builder
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
WORKDIR /build
COPY . .
RUN make app

FROM alpine
WORKDIR /app
COPY --from=builder /build/bin/app /bin/app
CMD app -d /var/lib/db/app.db -w :8080