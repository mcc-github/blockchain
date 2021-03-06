#!/bin/sh
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# This script publishes the blockchain docker images to mcc-github dockerhub and nexus3 repositories.
# when publishing the images to dockerhub the values for NS_PULL & NS_PUSH variables set to default values.
# when publishing the images to nexus3 repository the values for NS_PULL & NS_PUSH variables set like below
# NS_PULL=nexus3.mcc-github.org:10001/mcc-github & NS_PUSH=nexus3.mcc-github.org:10002/mcc-github
# since nexus has separate port numbers to pull and push the images to nexus3.

usage() {
    echo "Usage: $0 <username> <password>"
    echo "<username> and <password> credentials for the repository"
    echo "ENV:"
    echo "  NS_PULL=$NS_PULL"
    echo "  NS_PUSH=$NS_PUSH"
    echo "  VERSION=$VERSION"
    echo "  TWO_DIGIT_VERSION=$TWO_DIGIT_VERSION"
    exit 1
}

missing() {
    echo "Error: some image(s) missing from registry"
    echo "ENV:"
    echo "  NS_PULL=$NS_PULL"
    echo "  NS_PUSH=$NS_PUSH"
    echo "  VERSION=$VERSION"
    echo "  TWO_DIGIT_VERSION=$TWO_DIGIT_VERSION"
    exit 1
}

failed() {
    echo "Error: multiarch manifest push failed"
    echo "ENV:"
    echo "  NS_PULL=$NS_PULL"
    echo "  NS_PUSH=$NS_PUSH"
    echo "  VERSION=$VERSION"
    echo "  TWO_DIGIT_VERSION=$TWO_DIGIT_VERSION"
    exit 1
}

USER=${1:-nobody}
PASSWORD=${2:-nohow}
NS_PULL=${NS_PULL:-mcc-github}
NS_PUSH=${NS_PUSH:-mcc-github}
VERSION=${BASE_VERSION:-1.3.0}
TWO_DIGIT_VERSION=${TWO_DIGIT_VERSION:-1.3}

if [ "$#" -ne 2 ]; then
    usage
fi

# verify that manifest-tool is installed and found on PATH
if ! [ -x "$(command -v manifest-tool)" ]; then
    echo "manifest-tool not installed or not found on PATH"
    exit 1
fi

IMAGES="blockchain-baseos blockchain-peer blockchain-orderer blockchain-ccenv blockchain-tools"

# check that all images have been published
for image in ${IMAGES}; do
    docker pull "${NS_PULL}/${image}:amd64-${VERSION}" || missing
    docker pull "${NS_PULL}/${image}:s390x-${VERSION}" || missing
done

# push the multiarch manifest and tag with just $VERSION
for image in ${IMAGES}; do
    manifest-tool --username "${USER}" --password "${PASSWORD}" push from-args \
        --platforms linux/amd64,linux/s390x --template "${NS_PULL}/${image}:ARCH-${VERSION}" \
        --target "${NS_PUSH}/${image}:${VERSION}"
    manifest-tool --username "${USER}" --password "${PASSWORD}" push from-args \
        --platforms linux/amd64,linux/s390x --template "${NS_PULL}/${image}:ARCH-${VERSION}" \
        --target "${NS_PUSH}/${image}:${TWO_DIGIT_VERSION}"
done

# test that manifest is working as expected
for image in ${IMAGES}; do
    docker pull "${NS_PULL}/${image}:${VERSION}" || failed
    docker pull "${NS_PULL}/${image}:${TWO_DIGIT_VERSION}" || failed
done

echo "Successfully pushed multiarch manifest"
exit 0
