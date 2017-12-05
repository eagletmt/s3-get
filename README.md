# s3-get
Get object from S3

## Development
Use [dep](https://github.com/golang/dep) to manage dependency.

```
% dep ensure
% go build
% ./s3-get -help
```

## Package
Use [gox](https://github.com/mitchellh/gox) for cross compiling.

```
% tools/package.sh
```

## Release
Use [ghr](https://github.com/tcnksm/ghr) to upload packages to GitHub.

```
% tools/release.sh
```
