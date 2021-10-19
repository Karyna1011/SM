FROM golang:1.15

WORKDIR /go/src/github.com/phanos1011/sm
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /usr/local/bin/output .

FROM alpine:3.9

COPY --from=0 /usr/local/bin/output /usr/local/bin/output
RUN apk add --no-cache ca-certificates

ENTRYPOINT [ "output" ]
