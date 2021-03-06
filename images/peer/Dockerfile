# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG ALPINE_VER

FROM alpine:${ALPINE_VER} as peer-base
RUN apk add --no-cache tzdata

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	bash \
	gcc \
	git \
	make \
	musl-dev
ADD . $GOPATH/src/github.com/mcc-github/blockchain
WORKDIR $GOPATH/src/github.com/mcc-github/blockchain

FROM golang as peer
ARG GO_TAGS
RUN make peer GO_TAGS=${GO_TAGS}

FROM peer-base
ENV FABRIC_CFG_PATH /etc/mcc-github/blockchain
VOLUME /etc/mcc-github/blockchain
VOLUME /var/mcc-github
COPY --from=peer /go/src/github.com/mcc-github/blockchain/.build/bin /usr/local/bin
COPY --from=peer /go/src/github.com/mcc-github/blockchain/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=peer /go/src/github.com/mcc-github/blockchain/sampleconfig/core.yaml ${FABRIC_CFG_PATH}
EXPOSE 7051
CMD ["peer","node","start"]
