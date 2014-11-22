#!/bin/sh

for i in *; do if [ -d $i ]; then cd $i; go fmt *; cd .. ; fi; done
go fmt *.go

