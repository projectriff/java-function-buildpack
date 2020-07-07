/*
 * Copyright 2018-2020 the original author or authors.
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
	"io/ioutil"
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/projectriff/java-function-buildpack/java"
)

func testFunction(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "function-application")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "function-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes class function", func() {
		f, err := java.NewFunction(ctx.Application.Path, "test.handler")
		Expect(err).NotTo(HaveOccurred())

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = f.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(layer.LaunchEnvironment["SPRING_CLOUD_FUNCTION_FUNCTION_CLASS.override"]).To(Equal("test.handler"))
		Expect(layer.LaunchEnvironment["SPRING_CLOUD_FUNCTION_LOCATION.override"]).To(Equal(ctx.Application.Path))
	})

	it("contributes definition function", func() {
		f, err := java.NewFunction(ctx.Application.Path, "test-handler")
		Expect(err).NotTo(HaveOccurred())

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = f.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(layer.LaunchEnvironment["SPRING_CLOUD_FUNCTION_DEFINITION.override"]).To(Equal("test-handler"))
		Expect(layer.LaunchEnvironment["SPRING_CLOUD_FUNCTION_LOCATION.override"]).To(Equal(ctx.Application.Path))
	})

}
