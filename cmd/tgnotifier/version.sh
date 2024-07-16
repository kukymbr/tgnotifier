#!/usr/bin/env bash

set -e

version_file="version.go"

echo "package main" > "$version_file"

{
  echo ""
  echo "// The tgnotifier version, generated in main.go."
  echo "const version = \"$(git describe --tags)\""
} >> "$version_file"