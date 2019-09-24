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

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Dependency is the key identifying the riff java invoker in the buildpack plan.
const Dependency = "riff-invoker-java"

// Invoker represents the Java invoker contributed by the buildpack.
type Invoker struct {
	layer  layers.DependencyLayer
	layers layers.Layers
}

// Contributes makes the contribution to the launch layer.
func (i Invoker) Contribute() error {
	if err := i.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.Body("Expanding to %s", layer.Root)
		return helper.ExtractZip(artifact, layer.Root, 0)
	}, layers.Launch); err != nil {
		return err
	}

	streamingCommand := fmt.Sprintf("java -cp %s $JAVA_OPTS org.springframework.boot.loader.JarLauncher", i.layer.Root)
	command := fmt.Sprintf("streaming-http-adapter %s", streamingCommand)

	return i.layers.WriteApplicationMetadata(layers.Metadata{
		Processes: layers.Processes{
			layers.Process{Type: "function", Command: command},
			layers.Process{Type: "streaming-function", Command: streamingCommand},
			layers.Process{Type: "web", Command: command},
		},
	})
}

// NewInvoker creates a new instance returning true if the riff-invoker-java plan exists.
func NewInvoker(build build.Build) (Invoker, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(Dependency)
	if err != nil {
		return Invoker{}, false, err
	}
	if !ok {
		return Invoker{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return Invoker{}, false, err
	}

	dep, err := deps.Best(Dependency, p.Version, build.Stack)
	if err != nil {
		return Invoker{}, false, err
	}

	return Invoker{
		build.Layers.DependencyLayer(dep),
		build.Layers,
	}, true, nil
}
