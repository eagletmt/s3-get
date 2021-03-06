#!/bin/bash
set -ex

VERSION=0.1.1

rm -rf pkg

CGO_ENABLED=0 gox -verbose -output 'pkg/{{.OS}}_{{.Arch}}/s3-get'

cd pkg
for pkg in *
do
  zip -j s3-get_v${VERSION}_${pkg}.zip ${pkg}/*
  rm -r ${pkg}
done
sha256sum *.zip > SHA256SUMS
