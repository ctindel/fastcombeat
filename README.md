# Fastcombeat

![Fastcombeat Dashboard Screenshot](/fastcombeat_screenshot.png?raw=true "Fastcombeat Dashboard")

Fastcombeat estimates your internet download speed by performing a series of downloads from Netflix servers. Since there is no way for your ISP to cheat it or give the traffice priority without making your Netflix speeds faster in general, it is a very nice way to measure your real download speeds.  For more information see [https://fast.com](https://fast.com/).

Besides all of the typical beats configuration options, there is one more configuration option fastcombeat.use_ssl which should be set either to true or false to tell fastcombeat whether to use http or https endpoints for the speed test.

Credit to https://github.com/gesquive/fast-cli as the basis for the underlying bandwidth meter go code.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/ctindel/fastcombeat`

## Getting Started with Fastcombeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Fastcombeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Fastcombeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/ctindel/fastcombeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Fastcombeat run the command below. This will generate a binary
in the same directory with the name fastcombeat.

```
make
```


### Run

To run Fastcombeat with debugging output enabled, run:

```
./fastcombeat -c fastcombeat.yml -e -d "*"
```


### Test

To test Fastcombeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Fastcombeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Fastcombeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/ctindel/fastcombeat
git clone https://github.com/ctindel/fastcombeat ${GOPATH}/src/github.com/ctindel/fastcombeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
