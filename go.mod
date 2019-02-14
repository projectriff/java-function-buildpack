module github.com/projectriff/java-function-buildpack

require (
	github.com/buildpack/libbuildpack v1.10.0
	github.com/cloudfoundry/jvm-application-buildpack v1.0.0-M4
	github.com/cloudfoundry/libcfbuildpack v1.39.0
	github.com/cloudfoundry/openjdk-buildpack v1.0.0-M4
	github.com/onsi/gomega v1.4.3
	github.com/projectriff/riff-buildpack v0.1.1-0.20190207162816-5a5a4a635c5e
	github.com/sclevine/spec v1.2.0
)

replace github.com/projectriff/riff-buildpack => github.com/scothis/riff-buildpack v0.1.1-0.20190214002255-a465c569683a
