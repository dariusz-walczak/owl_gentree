#!/usr/bin/env bash

gotest . -v -cover -coverprofile /output/gentree_cover.out
go tool cover -html=/output/gentree_cover.out -o /output/cover.html
