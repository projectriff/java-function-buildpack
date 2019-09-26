/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package java

import (
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/openjdk-cnb/jre"
	"github.com/projectriff/libfnbuildpack/function"
)

type Buildpack struct{}

func (b Buildpack) Build(build build.Build) (int, error) {
	if f, ok, err := NewFunction(build); err != nil {
		return build.Failure(function.Error_ComponentInitialization), err
	} else if ok {
		if err := f.Contribute(); err != nil {
			return build.Failure(function.Error_ComponentContribution), err
		}

		if streaming, err := NewStreamingHTTPAdapter(build); err != nil {
			return build.Failure(function.Error_ComponentInitialization), err
		} else {
			if err := streaming.Contribute(); err != nil {
				return build.Failure(function.Error_ComponentContribution), err
			}
		}

		if invoker, ok, err := NewInvoker(build); err != nil {
			return build.Failure(function.Error_ComponentInitialization), err
		} else if ok {
			if err := invoker.Contribute(); err != nil {
				return build.Failure(function.Error_ComponentContribution), err
			}
		}
	}

	return build.Success()
}

func (b Buildpack) Detect(detect detect.Detect, metadata function.Metadata) (int, error) {
	return detect.Pass(buildplan.Plan{
		Provides: []buildplan.Provided{
			{Name: Dependency},
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
				Name: Dependency,
				Metadata: map[string]interface{}{
					Handler: metadata.Handler,
				},
			},
		},
	})
}

func (b Buildpack) Id() string {
	return "java"
}
