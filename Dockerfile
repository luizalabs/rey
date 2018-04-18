FROM golang:1.9 as builder
WORKDIR /go/src/github.com/luizalabs/rey
ADD . /go/src/github.com/luizalabs/rey
RUN CGO_ENABLED=0 go build -o rey

FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /go/src/github.com/luizalabs/rey/rey .
CMD ["./rey"]
