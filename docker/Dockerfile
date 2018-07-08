FROM golang:1.10.3-alpine3.7 as go-ethereum-build

ARG REPO_URL="https://github.com/ethereum/go-ethereum.git"
ARG CHECKOUT_TARGET="master"

RUN apk update && \
    apk upgrade && \
    apk add --no-cache gcc git linux-headers make musl-dev

RUN git clone ${REPO_URL} "/go/src/github.com/ethereum/go-ethereum"

WORKDIR "/go/src/github.com/ethereum/go-ethereum"

RUN git checkout ${CHECKOUT_TARGET}

RUN make all

FROM alpine:3.7 as main

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=go-ethereum-build "/go/src/github.com/ethereum/go-ethereum/build/bin/" /usr/local/bin/
