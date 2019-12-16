FROM alpine:3.10

RUN apk add --no-cache ca-certificates

ADD ./installation-operator /installation-operator

ENTRYPOINT ["/installation-operator"]
