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

package java

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type Invoker struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewInvoker(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) Invoker {
	return Invoker{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

func (i Invoker) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger

	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.ExtractZip(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.LaunchEnvironment.PrependPath("CLASSPATH", layer.Path)

		layer.Launch = true
		return layer, nil
	})
}

func (Invoker) Name() string {
	return "invoker"
}
