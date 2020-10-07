FROM golang:1.14-alpine AS builder

WORKDIR /build
COPY ./src/ /build/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:latest as alpine
RUN apk add --update --no-cache ca-certificates

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ab-testing-proxy /bin/ab-testing-proxy
ENTRYPOINT ["/bin/ab-testing-proxy"]
