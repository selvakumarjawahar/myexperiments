#!/bin/sh

set -e

docker run --rm -v $(pwd):/share golang:buster sh -c '
	apt-get update && \
	apt-get install -y protobuf-compiler && \
	go get github.com/golang/protobuf/protoc-gen-go && \
	protoc -I=/share/api/proto/ --go_out=/share/gen  /share/api/proto/netrounds_agent_untrusted_metrics.proto'