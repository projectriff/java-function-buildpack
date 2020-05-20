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
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	ApplicationPath  string
	Handler          string
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger
}

func NewFunction(applicationPath string, handler string) (Function, error) {
	return Function{
		ApplicationPath: applicationPath,
		Handler:         handler,
		LayerContributor: libpak.NewLayerContributor(bard.FormatIdentity("Java", handler),
			map[string]interface{}{"handler": handler}),
	}, nil
}

func (f Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger

	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		if len(f.Handler) > 0 {
			if strings.ContainsAny(f.Handler, ".") {
				layer.LaunchEnvironment.Override("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.Handler)
			} else {
				layer.LaunchEnvironment.Override("SPRING_CLOUD_FUNCTION_DEFINITION", f.Handler)
			}
		}

		layer.LaunchEnvironment.Override("SPRING_CLOUD_FUNCTION_LOCATION", f.ApplicationPath)

		layer.Launch = true
		return layer, nil
	})
}

func (Function) Name() string {
	return "function"
}
