/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_test

import (
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-cnb/jre"
	"github.com/onsi/gomega"
	"github.com/projectriff/java-function-buildpack/java"
	"github.com/projectriff/libfnbuildpack/function"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			b java.Buildpack
			f *test.DetectFactory
		)

		it.Before(func() {
			b = java.Buildpack{}
			f = test.NewDetectFactory(t)
		})

		when("id", func() {

			it("returns id", func() {
				g.Expect(b.Id()).To(gomega.Equal("java"))
			})
		})

		it("passes with handler", func() {
			g.Expect(b.Detect(f.Detect, function.Metadata{Handler: "test.handler"})).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: java.Dependency},
				},
				Requires: []buildplan.Required{
					{
						Name: jre.Dependency,
						Metadata: map[string]interface{}{
							jre.LaunchContribution: true,
						},
					},
					{Name: jvmapplication.Dependency},
					{
						Name: java.Dependency,
						Metadata: map[string]interface{}{
							java.Handler: "test.handler",
						},
					},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
