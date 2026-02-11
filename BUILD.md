# Safecastbeat

Welcome to Safecastbeat.

## Getting Started with Safecastbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.24+


### Clone

To clone Safecastbeat from the git repository, run the following commands:

```
git clone https://github.com/radoondas/safecastbeat
cd safecastbeat
```


### Build

To build the binary for Safecastbeat run the command below. This will generate a binary
in the same directory with the name safecastbeat.

```
go build -o safecastbeat .
```

Or via Make (delegates to libbeat's Makefile):

```
make
```


### Run

To run Safecastbeat with debugging output enabled, run:

```
./safecastbeat -c safecastbeat.yml -e -d "*"
```


### Test

To test Safecastbeat, run the following command:

```
go test ./...
```


### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean Safecastbeat source code, run the following commands:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


## Cross-Compilation

Build binaries for all supported platforms (linux, darwin, windows; amd64, arm64):

```
./scripts/build.sh
```

Binaries are output to `dist/`.


## Packaging

Create distributable archives (tar.gz, zip, sha512 checksums):

```
./scripts/package.sh
```

To also create deb/rpm packages, install [nfpm](https://nfpm.goreleaser.com/) first:

```
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.41.1
./scripts/package.sh
```

Verify all expected artifacts:

```
./scripts/verify-assets.sh
```


## Docker

Build a local Docker image:

```
docker buildx build --platform linux/amd64 -t safecastbeat:local .
```


## Release

Releases are automated via GitHub Actions. Push a tag to trigger:

```
git tag v9.3.0
git push origin v9.3.0
```

This builds all platforms, packages archives and deb/rpm, builds and pushes a multi-arch Docker image, and creates a GitHub Release with all artifacts.

Manual dry-run (without publishing) is available via `workflow_dispatch` in the Actions tab.
