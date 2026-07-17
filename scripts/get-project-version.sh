#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Chen Linxuan <me@black-desk.cn>
#
# SPDX-License-Identifier: MIT

# This script prints the project version to STDOUT
# for automatically creating a tag in CI.
#
# The version is read from internal/version/version.go so that the Go
# source remains the single source of truth. The same constant is
# overridden at release time via goreleaser ldflags.

set -e
set -o pipefail

CURRENT_SOURCE_FILE_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]:-$0}")" && pwd)"
REPO_ROOT="$(cd -- "$CURRENT_SOURCE_FILE_DIR/.." && pwd)"

VERSION_FILE="$REPO_ROOT/internal/version/version.go"

if [ ! -f "$VERSION_FILE" ]; then
	echo "::error::Version source file not found: $VERSION_FILE" >&2
	exit 1
fi

version="$(sed -n 's/^var[[:space:]]*Version[[:space:]]*=[[:space:]]*"\([^"]*\)".*/\1/p' "$VERSION_FILE" | head -n1)"

if [ -z "$version" ]; then
	echo "::error::Failed to parse Version from $VERSION_FILE" >&2
	exit 1
fi

printf '%s\n' "$version"
