FROM golang:alpine3.17 as builder

ARG swagger_version="latest"

RUN apk update \
 && apk add --no-cache \
        bash \
        binutils \
        gcc \
        git \
        make \
        musl-dev

RUN go install github.com/swaggo/swag/cmd/swag@${swagger_version}

WORKDIR /tmp/go
ADD . /tmp/go
COPY config.yaml /tmp/go
ARG goreleaser_flags
RUN make build

FROM alpine:3.17
COPY --from=builder /tmp/go/NORSI-TRANS /tmp/go/NORSI-TRANS
COPY --from=builder /tmp/go/config.yaml .
COPY --from=builder /tmp/go/swagger/doc.json /app/swagger/doc.json
RUN apk update
WORKDIR .
CMD ["/tmp/go/NORSI-TRANS"]
EXPOSE 8000