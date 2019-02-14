/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/projectriff/riff-buildpack/invoker"
	"github.com/projectriff/riff-buildpack/metadata"
)

type JavaBuildpack struct {
	name string
}

func (b *JavaBuildpack) Name() string {
	return b.name
}

func (b *JavaBuildpack) Detect(detect detect.Detect, metadata metadata.Metadata) (bool, error) {
	// Try java
	_, ok := detect.BuildPlan[jvmapplication.Dependency]
	return ok, nil
}

func (b *JavaBuildpack) BuildPlan(detect detect.Detect, metadata metadata.Metadata) buildplan.BuildPlan {
	return BuildPlanContribution(detect, metadata)
}

func (b *JavaBuildpack) Invoker(build build.Build) (invoker.Invoker, bool, error) {
	return NewJavaInvoker(build)
}

func NewBuildpack() invoker.Buildpack {
	return &JavaBuildpack{
		name: "java",
	}
}
