#!/bin/bash

PLATFORMS=( "$@" )

for _platform in "${PLATFORMS[@]}"; do
	IFS=',' read -ra _os_arch <<< "$_platform"
	_os="${_os_arch[0]}"
	_arch="${_os_arch[1]}"
	GOOS="$_os" GOARCH="$_arch" make build_dist
done
