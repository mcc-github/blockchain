# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
ARG GO_VER
ARG ALPINE_VER
FROM golang:${GO_VER}-alpine as golang

RUN apk add --no-cache \
	bash \
	gcc \
	git \
	make \
	musl-dev;

ADD . $GOPATH/src/github.com/mcc-github/blockchain
WORKDIR $GOPATH/src/github.com/mcc-github/blockchain
ENV EXECUTABLES go git

FROM golang as tools
RUN make configtxgen configtxlator cryptogen peer discover idemixgen

FROM golang:${GO_VER}-alpine
# git is required to support `go list -m`
RUN apk add --no-cache \
	bash \
	git \
	jq \
	tzdata;
ENV FABRIC_CFG_PATH /etc/mcc-github/blockchain
VOLUME /etc/mcc-github/blockchain
COPY --from=tools /go/src/github.com/mcc-github/blockchain/.build/bin /usr/local/bin
COPY --from=tools /go/src/github.com/mcc-github/blockchain/sampleconfig ${FABRIC_CFG_PATH}
