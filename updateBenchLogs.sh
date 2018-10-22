#!/bin/bash

cd base64

go test -bench . > _bench.log

cd ../between

go test -bench . > _bench.log

cd ../contains

go test -bench . > _bench.log

cd ../foreach

go test -bench . > _bench.log

cd ..