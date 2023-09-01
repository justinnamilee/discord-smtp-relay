FROM golang:1.15-alpine as builder
WORKDIR /build
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o discord-smtp .

FROM golang:1.15-alpine as test
WORKDIR /test
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates
VOLUME /test
ENTRYPOINT CGO_ENABLED=0 GOOS=linux go test ./...

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/discord-smtp .
VOLUME /message.tt
ENTRYPOINT [ "./discord-smtp" ]