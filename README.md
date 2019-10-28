# `java-function-buildpack` [![Build Status](https://travis-ci.com/projectriff/java-function-buildpack.svg?branch=master)](https://travis-ci.com/projectriff/java-function-buildpack)

The Java Function Buildpack is a Cloud Native Buildpack V3 that provides the riff [Java Function Invoker](https://github.com/projectriff/java-function-invoker) to functions.

This buildpack is designed to work in collaboration with other buildpacks, which are tailored to
support (and know how to build / run) languages supported by riff.

## In Plain English

In a nutshell, when combined with the other buildpacks present in the [riff builder](https://github.com/projectriff/builder) what this means (and especially when dealing with the riff CLI which takes care of the creation of the `riff.toml` file for you):

- The presence of a `pom.xml` or `build.gradle` file will result in the compilation and execution of a java function, thanks to the [java invoker](https://github.com/projectriff/java-function-invoker)
  1. the `--handler` flag is optional in certain cases, as documented by the java invoker
- Ambiguity in the detection process will result in a build failure
- The presence of the `--invoker` flag will entirely bypass the detection mechanism and force a given language/invoker

## Detailed Buildpack Behavior

### Detection Phase

Detection passes if

- a `$APPLICATION_ROOT/riff.toml` exists and
- the build plan already contains a `jvm-application` key (typically because a JVM based application was detected by the [java buildpack](https://github.com/cloudfoundry/build-system-buildpack))

If detection passes, the buildpack will contribute an `openjdk-jre` key with `launch` metadata to instruct
the `openjdk-buildpack` to provide a JRE. It will also add a `riff-invoker-java` key and `handler`
metadata extracted from the riff metadata.

If several languages are detected simultaneously, the detect phase errors out.
The `override` key in `riff.toml` can be used to bypass detection and force the use of a particular invoker.

### Build Phase

If a java build has been detected

- Contributes the riff Java Invoker to a launch layer, set as the main java entry point with `function.uri = <build-directory>?handler=<handler>` set as an environment variable.

The function behavior is exposed _via_ standard buildpack [process types](https://github.com/buildpack/spec/blob/master/buildpack.md#launch):

- Contributes `web` process
- Contributes `function` process

## How to Build

### Prerequisites
To build the java-function-buildpack you'll need

- Go 1.13+
- to run acceptance tests:
  - a running local docker daemon
  - the [`pack`](https://github.com/buildpack/pack command line tool, [version](https://github.com/buildpack/pack/releases) `>= 0.5.0`.

You can build the buildpack by running

```bash
make
```

This will package (with pre-downloaded cache layers) the buildpack in the
`artifactory/io/projectriff/java/io.projectriff.java/latest` directory. That can be used as a `uri` in a `builder.toml`
file of a builder (see https://github.com/projectriff/builder)

## License

This buildpack is released under version 2.0 of the [Apache License](https://www.apache.org/licenses/LICENSE-2.0).
