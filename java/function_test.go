/*
 * Copyright 2018-2019 the original author or authors.
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

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/projectriff/java-function-buildpack/java"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestFunction(t *testing.T) {
	spec.Run(t, "Function", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan exists", func() {
			f.AddPlan(buildpackplan.Plan{
				Name: java.Dependency,
				Metadata: map[string]interface{}{
					java.Handler: "test.handler",
				},
			})

			_, ok, err := java.NewFunction(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeTrue())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := java.NewFunction(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeFalse())
		})

		it("contributes function class to launch", func() {
			f.AddPlan(buildpackplan.Plan{
				Name: java.Dependency,
				Metadata: map[string]interface{}{
					java.Handler: "test.handler",
				},
			})

			h, _, err := java.NewFunction(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(h.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("java-function")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(layer).To(test.HaveOverrideLaunchEnvironment("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", "test.handler"))
			g.Expect(layer).To(test.HaveOverrideLaunchEnvironment("SPRING_CLOUD_FUNCTION_LOCATION", f.Build.Application.Root))
		})

		it("contributes function definition to launch", func() {
			f.AddPlan(buildpackplan.Plan{
				Name: java.Dependency,
				Metadata: map[string]interface{}{
					java.Handler: "test-handler",
				},
			})

			h, _, err := java.NewFunction(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(h.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("java-function")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(layer).To(test.HaveOverrideLaunchEnvironment("SPRING_CLOUD_FUNCTION_DEFINITION", "test-handler"))
			g.Expect(layer).To(test.HaveOverrideLaunchEnvironment("SPRING_CLOUD_FUNCTION_LOCATION", f.Build.Application.Root))
		})
	}, spec.Report(report.Terminal{}))
}
