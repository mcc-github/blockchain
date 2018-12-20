# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
ARG GO_VER=1.11.1
ARG ALPINE_VER=3.8
FROM alpine:${ALPINE_VER} as base
RUN apk add --no-cache tzdata libltdl

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	gcc \
	musl-dev \
	git \
	bash \
	make;
ADD . $GOPATH/src/github.com/mcc-github/blockchain
WORKDIR $GOPATH/src/github.com/mcc-github/blockchain
ENV EXECUTABLES go git

FROM golang as peer
RUN make peer

FROM base
ENV FABRIC_CFG_PATH /etc/mcc-github/blockchain
VOLUME /etc/mcc-github/blockchain
VOLUME /var/mcc-github
COPY --from=peer /go/src/github.com/mcc-github/blockchain/.build/bin /usr/local/bin
COPY --from=peer /go/src/github.com/mcc-github/blockchain/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=peer /go/src/github.com/mcc-github/blockchain/sampleconfig/core.yaml ${FABRIC_CFG_PATH}
EXPOSE 7051
CMD ["peer","node","start"]
