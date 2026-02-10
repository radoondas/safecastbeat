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

alternatively:
```
make testsuite
make unit-tests
```

The test coverage is reported in the folder `./build/coverage/`

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


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/). To build packages of your beat, run the following command:

```
mage package
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
