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
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/projectriff/java-function-buildpack/java"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect java.Detect
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "command")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
	})

	it("passes without riff.toml", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "riff-java"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jre", Metadata: map[string]interface{}{"launch": true}},
						{Name: "jvm-application"},
						{Name: "streaming-http-adapter"},
					},
				},
			},
		}))
	})

	it("passes with with riff.toml", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "riff.toml"), []byte(`
artifact = "test-artifact"
`), 0644))

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "riff-java"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jre", Metadata: map[string]interface{}{"launch": true}},
						{Name: "jvm-application"},
						{Name: "streaming-http-adapter"},
						{Name: "riff-java", Metadata: map[string]interface{}{"artifact": "test-artifact"}},
					},
				},
			},
		}))
	})

}
