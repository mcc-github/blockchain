#!/bin/sh
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


{
echo "## $2";
date;
echo ""
git log "$1..HEAD"  --oneline | grep -v Merge | sed -e "s/\[\(FAB-[0-9]*\)\]/\[\1\](https:\/\/jira.mcc-github.org\/browse\/\1\)/" -e "s/ \(FAB-[0-9]*\)/ \[\1\](https:\/\/jira.mcc-github.org\/browse\/\1\)/" -e "s/\([0-9|a-z]*\)/* \[\1\](https:\/\/github.com\/mcc-github\/blockchain\/commit\/\1)/"
echo ""
echo ""
cat CHANGELOG.md
} >> CHANGELOG.new
mv -f CHANGELOG.new CHANGELOG.md
