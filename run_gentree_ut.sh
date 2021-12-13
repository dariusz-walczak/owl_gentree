#!/usr/bin/env bash

set -ex

go test -v -cover -coverprofile /output/gentree_cover.out .
