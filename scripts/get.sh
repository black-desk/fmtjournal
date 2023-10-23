#!/bin/bash

set -e
set -x

# shellcheck disable=SC1090
. <(curl -fL "https://raw.githubusercontent.com/black-desk/get/master/get.sh") \
	black-desk fmtjournal

curl -fLO "https://raw.githubusercontent.com/black-desk/fmtjournal/v$VERSION/tools/journalctl"

$SUDO install -m755 -D "$TMP_DIR/fmtjournal" "$PREFIX/bin/fmtjournal"
$SUDO install -m755 -D "$TMP_DIR/journalctl" "$PREFIX/bin/journalctl"
