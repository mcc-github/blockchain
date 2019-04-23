#!/bin/bash -eu
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

# ----------------------------------------------------------------
# set custom limits
# ----------------------------------------------------------------

cat <<EOF >/etc/security/limits.d/99-mcc-github.conf
# custom limits for mcc-github development

*       soft    nofile          10000
*       hard    nofile          16384
EOF
