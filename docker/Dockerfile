FROM golang:1.10.3-alpine3.7 as golang-base

RUN apk update && \
    apk upgrade && \
    apk add --no-cache gcc linux-headers make musl-dev

FROM golang-base as microcosm-build

WORKDIR "/go/src/github.com/nkashy1/microcosm"

COPY . .

RUN go build

FROM ethereum/client-go:v1.8.12 as main

COPY --from=microcosm-build "/go/src/github.com/nkashy1/microcosm/microcosm" /usr/local/bin/

ENTRYPOINT []
CMD []