#!/bin/bash
#
rm -f main
go build *.go
sleep 1
scp main jenkins:/data/app/falcon-test/main
