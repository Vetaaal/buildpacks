// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Implements java/appengine buildpack.
// The appengine buildpack sets the image entrypoint.
package main

import (
	"fmt"

	"github.com/GoogleCloudPlatform/buildpacks/pkg/appengine"
	"github.com/GoogleCloudPlatform/buildpacks/pkg/appstart"
	gcp "github.com/GoogleCloudPlatform/buildpacks/pkg/gcpbuildpack"
	"github.com/GoogleCloudPlatform/buildpacks/pkg/java"
)

func main() {
	gcp.Main(detectFn, buildFn)
}

func detectFn(ctx *gcp.Context) (gcp.DetectResult, error) {
	return gcp.OptInAlways(), nil
}

func buildFn(ctx *gcp.Context) error {
	return appengine.Build(ctx, "java", entrypoint)
}

func entrypoint(ctx *gcp.Context) (*appstart.Entrypoint, error) {
	webXMLExists, err := ctx.FileExists("WEB-INF", "appengine-web.xml")
	if err != nil {
		return nil, err
	}
	if webXMLExists {
		return &appstart.Entrypoint{
			Type:    appstart.EntrypointGenerated.String(),
			Command: "serve WEB-INF/appengine-web.xml",
		}, nil
	}

	executable, err := java.ExecutableJar(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding executable jar: %w", err)
	}

	return &appstart.Entrypoint{
		Type:    appstart.EntrypointGenerated.String(),
		Command: "serve " + executable,
	}, nil
}
