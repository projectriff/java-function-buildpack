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

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/projectriff/libfnbuildpack"
)

type Detect struct{}

func (Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{
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
	}

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if ok, err := libfnbuildpack.IsRiff(context.Application.Path, cr); err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to determine if application is riff\n%w", err)
	} else if !ok {
		return result, nil
	}

	metadata, err := libfnbuildpack.Metadata(context.Application.Path, cr)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to read riff metadata\n%w", err)
	}

	result.Plans[0].Requires = append(result.Plans[0].Requires, libcnb.BuildPlanRequire{
		Name:     "riff-java",
		Metadata: metadata,
	})

	return result, nil
}
