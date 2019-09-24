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

package java

import (
	"fmt"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Handler is the key identifying the riff handler metadata in the build plan
const Handler = "handler"

// Function represents the function to be executed.
type Function struct {
	application application.Application
	handler     string
	layer       layers.Layer
}

// Contributes makes the contribution to the launch layer.
func (f Function) Contribute() error {
	return f.layer.Contribute(marker{"Java", f.handler}, func(layer layers.Layer) error {
		if len(f.handler) > 0 {
			if strings.ContainsAny(f.handler, ".") {
				if err := layer.OverrideLaunchEnv("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.handler); err != nil {
					return err
				}
			} else {
				if err := layer.OverrideLaunchEnv("SPRING_CLOUD_FUNCTION_DEFINITION", f.handler); err != nil {
					return err
				}
			}
		}

		return layer.OverrideLaunchEnv("SPRING_CLOUD_FUNCTION_LOCATION", f.application.Root)
	}, layers.Launch)
}

// NewFunction creates a new instance returning true if the riff-invoker-java plan exists.
func NewFunction(build build.Build) (Function, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(Dependency)
	if err != nil {
		return Function{}, false, err
	}
	if !ok {
		return Function{}, false, nil
	}

	fa, ok := p.Metadata[Handler]
	if !ok {
		fa = ""
	}

	exec, ok := fa.(string)
	if !ok {
		return Function{}, false, fmt.Errorf("handler metadata of incorrect type: %v", p.Metadata[Handler])
	}

	return Function{
		build.Application,
		exec,
		build.Layers.Layer("java-function"),
	}, true, nil
}

type marker struct {
	Type    string `toml:"type"`
	Handler string `toml:"handler"`
}

func (m marker) Identity() (string, string) {
	return m.Type, m.Handler
}
