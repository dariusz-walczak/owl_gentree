#!/usr/bin/env bash

set -ex

gotest -v -cover -coverprofile /output/gentree_cover.out .
