# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
ARG GO_VER
ARG ALPINE_VER
FROM alpine:${ALPINE_VER} as base
RUN apk add --no-cache tzdata libltdl

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	gcc \
	musl-dev \
	git \
	libtool \
	bash \
	make;
ADD . $GOPATH/src/github.com/mcc-github/blockchain
WORKDIR $GOPATH/src/github.com/mcc-github/blockchain
ENV EXECUTABLES go git

FROM golang as orderer
ARG GO_TAGS
RUN make orderer GO_TAGS=${GO_TAGS}

FROM base
ENV FABRIC_CFG_PATH /etc/mcc-github/blockchain
VOLUME /etc/mcc-github/blockchain
VOLUME /var/mcc-github
COPY --from=orderer /go/src/github.com/mcc-github/blockchain/.build/bin /usr/local/bin
COPY --from=orderer /go/src/github.com/mcc-github/blockchain/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=orderer /go/src/github.com/mcc-github/blockchain/sampleconfig/orderer.yaml ${FABRIC_CFG_PATH}
COPY --from=orderer /go/src/github.com/mcc-github/blockchain/sampleconfig/configtx.yaml ${FABRIC_CFG_PATH}
EXPOSE 7050
CMD ["orderer"]
