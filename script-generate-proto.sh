#!/usr/bin/env bash
set -e

<<COMMENT
go get -v github.com/golang/protobuf/protoc-gen-go

printf "generating proto in raft/raftpb\n"
protoc --go_out=. raft/raftpb/*.proto

printf "generating proto in wal/walpb\n"
protoc --go_out=. wal/walpb/*.proto
COMMENT

# for now, be conservative about what version of protoc we expect
if ! [[ $(protoc --version) =~ "3.0.0" ]]; then
	echo "could not find protoc 3.0.0, is it installed + in PATH?"
	exit 255
fi

echo "Installing gogo/protobuf..."
GOGOPROTO_ROOT="$GOPATH/src/github.com/gogo/protobuf"
rm -rf $GOGOPROTO_ROOT
go get -u github.com/gogo/protobuf/{proto,protoc-gen-gogo,gogoproto,protoc-gen-gofast}
go get -u golang.org/x/tools/cmd/goimports
pushd "${GOGOPROTO_ROOT}"
	git reset --hard HEAD
	make install
popd

printf "Generating raftpb\n"
protoc --gofast_out=plugins=grpc:. \
	--proto_path=$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf:. \
	raftpb/*.proto;

printf "Generating walpb\n"
protoc --gofast_out=plugins=grpc:. \
	--proto_path=$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf:. \
	walpb/*.proto;
