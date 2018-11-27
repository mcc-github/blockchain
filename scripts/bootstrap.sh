#! /bin/bash -ue
set -o pipefail

#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# if version not passed in, default to latest released version
export VERSION=1.3.0
# if ca version not passed in, default to latest released version
export CA_VERSION=$VERSION
# current version of thirdparty images (couchdb, kafka and zookeeper) released
export THIRDPARTY_IMAGE_VERSION=0.4.13
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')")
export MARCH=$(uname -m)

printHelp() {
cat <<EOF

    Usage: $(basename $0) [OPTION]... [version [ca_version [thirdparty_version]]]

    OPTIONS:
        -h|--help Print this message
        -d        Bypass docker image download
        -s        Bypass blockchain-samples repo clone
        -b        Bypass download of platform-specific binaries

    EXAMPLES:
        To download docker images for version $VERSION
        ==> $0 -sb $VERSION

EOF
}

dockerFabricPull() {
  local FABRIC_TAG=$1
  for IMAGES in peer orderer ccenv javaenv tools; do
      echo "==> FABRIC IMAGE: $IMAGES"
      echo
      docker pull mcc-github/blockchain-$IMAGES:$FABRIC_TAG
      docker tag mcc-github/blockchain-$IMAGES:$FABRIC_TAG mcc-github/blockchain-$IMAGES
  done
}

dockerThirdPartyImagesPull() {
  local THIRDPARTY_TAG=$1
  for IMAGES in couchdb kafka zookeeper; do
      echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
      echo
      docker pull mcc-github/blockchain-$IMAGES:$THIRDPARTY_TAG
      docker tag mcc-github/blockchain-$IMAGES:$THIRDPARTY_TAG mcc-github/blockchain-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> FABRIC CA IMAGE"
      echo
      docker pull mcc-github/blockchain-ca:$CA_TAG
      docker tag mcc-github/blockchain-ca:$CA_TAG mcc-github/blockchain-ca
}

samplesInstall() {
  # clone (if needed) mcc-github/blockchain-samples and checkout corresponding
  # version to the binaries and docker images to be downloaded
  if [ -d first-network ]; then
    # if we are in the blockchain-samples repo, checkout corresponding version
    echo "===> Checking out v${VERSION} of mcc-github/blockchain-samples"
    git checkout v${VERSION}
  elif [ -d blockchain-samples ]; then
    # if blockchain-samples repo already cloned and in current directory,
    # cd blockchain-samples and checkout corresponding version
    echo "===> Checking out v${VERSION} of mcc-github/blockchain-samples"
    cd blockchain-samples && git checkout v${VERSION}
  else
    echo "===> Cloning mcc-github/blockchain-samples repo and checkout v${VERSION}"
    git clone -b master https://github.com/mcc-github/blockchain-samples.git && cd blockchain-samples && git checkout v${VERSION}
  fi
}

# Incrementally downloads the .tar.gz file locally first, only decompressing it
# after the download is complete. This is slower than binaryDownload() but
# allows the download to be resumed.
binaryIncrementalDownload() {
      local BINARY_FILE=$1
      local URL=$2
      curl -f -s -C - ${URL} -o ${BINARY_FILE} || rc=$?
      # Due to limitations in the current Nexus repo:
      # curl returns 33 when there's a resume attempt with no more bytes to download
      # curl returns 2 after finishing a resumed download
      # with -f curl returns 22 on a 404
      if [ "$rc" = 22 ]; then
	  # looks like the requested file doesn't actually exist so stop here
	  return 22
      fi
      if [ -z "$rc" ] || [ $rc -eq 33 ] || [ $rc -eq 2 ]; then
          # The checksum validates that RC 33 or 2 are not real failures
          echo "==> File downloaded. Verifying the md5sum..."
          localMd5sum=$(md5sum ${BINARY_FILE} | awk '{print $1}')
          remoteMd5sum=$(curl -s ${URL}.md5)
          if [ "$localMd5sum" == "$remoteMd5sum" ]; then
              echo "==> Extracting ${BINARY_FILE}..."
              tar xzf ./${BINARY_FILE} --overwrite
	      echo "==> Done."
              rm -f ${BINARY_FILE} ${BINARY_FILE}.md5
          else
              echo "Download failed: the local md5sum is different from the remote md5sum. Please try again."
              rm -f ${BINARY_FILE} ${BINARY_FILE}.md5
              exit 1
          fi
      else
          echo "Failure downloading binaries (curl RC=$rc). Please try again and the download will resume from where it stopped."
          exit 1
      fi
}

# This will attempt to download the .tar.gz all at once, but will trigger the
# binaryIncrementalDownload() function upon a failure, allowing for resume
# if there are network failures.
binaryDownload() {
      local BINARY_FILE=$1
      local URL=$2
      echo "===> Downloading: " ${URL}
      # Check if a previous failure occurred and the file was partially downloaded
      if [ -e ${BINARY_FILE} ]; then
          echo "==> Partial binary file found. Resuming download..."
          binaryIncrementalDownload ${BINARY_FILE} ${URL}
      else
          if ! curl --output - ${URL} | tar -xz; then
              echo "==> There was an error downloading the binary file. Switching to incremental download."
              echo "==> Downloading file..."
              binaryIncrementalDownload ${BINARY_FILE} ${URL}
	  else
	      echo "==> Done."
          fi
      fi
}

binariesInstall() {
  echo "===> Downloading version ${FABRIC_TAG} platform specific blockchain binaries"
  binaryDownload ${BINARY_FILE} https://nexus.mcc-github.org/content/repositories/releases/org/mcc-github/blockchain/mcc-github-blockchain/${ARCH}-${VERSION}/${BINARY_FILE}
  if [ $? -eq 22 ]; then
     echo
     echo "------> ${FABRIC_TAG} platform specific blockchain binary is not available to download <----"
     echo
   fi

  echo "===> Downloading version ${CA_TAG} platform specific blockchain-ca-client binary"
  binaryDownload ${CA_BINARY_FILE} https://nexus.mcc-github.org/content/repositories/releases/org/mcc-github/blockchain-ca/mcc-github-blockchain-ca/${ARCH}-${CA_VERSION}/${CA_BINARY_FILE}
  if [ $? -eq 22 ]; then
     echo
     echo "------> ${CA_TAG} blockchain-ca-client binary is not available to download  (Available from 1.1.0-rc1) <----"
     echo
   fi
}

dockerInstall() {
  which docker >& /dev/null
  NODOCKER=$?
  if [ "${NODOCKER}" == 0 ]; then
	  echo "===> Pulling blockchain Images"
	  dockerFabricPull ${FABRIC_TAG}
	  echo "===> Pulling blockchain ca Image"
	  dockerCaPull ${CA_TAG}
	  echo "===> Pulling thirdparty docker images"
	  dockerThirdPartyImagesPull ${THIRDPARTY_TAG}
	  echo
	  echo "===> List out mcc-github docker images"
	  docker images | grep mcc-github*
  else
    echo "========================================================="
    echo "Docker not installed, bypassing download of Fabric images"
    echo "========================================================="
  fi
}

###########  End of Function Definitions ######################################

DOCKER=true
SAMPLES=true
BINARIES=true

# Use 'getopt' to parse the command line arguments
options=$(getopt -o b,d,h,s --long help -- "$@")

# Extract all flags from command line
# Leave arguments
eval set -- "$options"
while [[ $# > 0 ]]; do
    case "$1" in
    -h|--help) printHelp ; exit  ;;
    -b)        BINARIES=false    ;;
    -d)        DOCKER=false      ;;
    -s)        SAMPLES=false     ;;
    --)        shift ; break     ;;
    esac
    shift
done

# All that is left is arguments (no flags)
case $# in
     0)  ;;
     1)  VERSION=$1 ;;
     2)  VERSION=$1 ; CA_VERSION=$2 ;;
     3)  VERSION=$1 ; CA_VERSION=$2 ; THIRDPARTY_IMAGE_VERSION=$3 ;;
     *)  echo "ERROR: Unexpected argument: $4" ; printHelp ; exit 1 ;;
esac

# prior to 1.2.0 architecture was determined by uname -m
if [[ $VERSION =~ ^1\.[0-1]\.* ]]; then
  export FABRIC_TAG=${MARCH}-${VERSION}
  export CA_TAG=${MARCH}-${CA_VERSION}
  export THIRDPARTY_TAG=${MARCH}-${THIRDPARTY_IMAGE_VERSION}
else
  # starting with 1.2.0, multi-arch images will be default
  : ${CA_TAG:="$CA_VERSION"}
  : ${FABRIC_TAG:="$VERSION"}
  : ${THIRDPARTY_TAG:="$THIRDPARTY_IMAGE_VERSION"}
fi

BINARY_FILE=mcc-github-blockchain-${ARCH}-${VERSION}.tar.gz
CA_BINARY_FILE=mcc-github-blockchain-ca-${ARCH}-${CA_VERSION}.tar.gz

if $SAMPLES; then
  echo
  echo "Installing mcc-github/blockchain-samples repo"
  echo
  samplesInstall
fi
if $BINARIES; then
  echo
  echo "Installing Hyperledger Fabric binaries"
  echo
  binariesInstall
fi
if $DOCKER; then
  echo
  echo "Installing Hyperledger Fabric docker images"
  echo
  dockerInstall
fi
