module github.com/projectriff/java-function-buildpack

require (
	github.com/buildpack/libbuildpack v1.13.0
	github.com/cloudfoundry/jvm-application-cnb v1.0.0-M7
	github.com/cloudfoundry/libcfbuildpack v1.51.0
	github.com/cloudfoundry/openjdk-cnb v1.0.0-M7
	github.com/onsi/gomega v1.5.0
	github.com/projectriff/libfnbuildpack v0.2.0
	github.com/sclevine/spec v1.2.0
)

replace github.com/projectriff/libfnbuildpack => github.com/scothis/libfnbuildpack v0.2.1-0.20190501155641-93021491f2ee
