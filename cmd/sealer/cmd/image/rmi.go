// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package image

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/sealerio/sealer/pkg/clusterfile"
	"github.com/sealerio/sealer/pkg/define/options"
	"github.com/sealerio/sealer/pkg/imageengine"
)

var removeOpts *options.RemoveImageOptions

var longNewRmiCmdDescription = ``

var exampleForRmiCmd = `
  sealer rmi docker.io/sealerio/kubernetes:v1-22-15-sealerio-2

prune dangling images:
  sealer rmi --prune/-p

force removal of the image and any containers using the image:
  sealer rmi docker.io/sealerio/kubernetes:v1-22-15-sealerio-2 --force/-f

`

// NewRmiCmd rmiCmd represents the rmi command
func NewRmiCmd() *cobra.Command {
	rmiCmd := &cobra.Command{
		Use:   "rmi",
		Short: "remove local images",
		// TODO: add long description.
		Long:    longNewRmiCmdDescription,
		Example: exampleForRmiCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !removeOpts.Force {
				args = filterClusterImage(args)
			}
			return runRemove(args)
		},
	}
	removeOpts = &options.RemoveImageOptions{}
	flags := rmiCmd.Flags()
	flags.BoolVarP(&removeOpts.Prune, "prune", "p", false, "prune dangling images")
	flags.BoolVarP(&removeOpts.Force, "force", "f", false, "force removal of the image and any containers using the image")
	return rmiCmd
}

func runRemove(images []string) error {
	removeOpts.ImageNamesOrIDs = images
	engine, err := imageengine.NewImageEngine(options.EngineGlobalConfigurations{})
	if err != nil {
		return err
	}

	return engine.RemoveImage(removeOpts)
}

// getRunningClusterImages get cluster image name and id
// if cluster image not exist, return empty string
func getRunningClusterImages() (string, string, error) {
	cf, _, err := clusterfile.GetActualClusterFile()
	if err != nil {
		return "", "", err
	}

	imageName := cf.GetCluster().Spec.Image
	engine, err := imageengine.NewImageEngine(options.EngineGlobalConfigurations{})
	if err != nil {
		return "", "", err
	}

	image, err := engine.Inspect(&options.InspectOptions{ImageNameOrID: imageName})
	if err != nil {
		return "", "", err
	}

	return imageName, image.ID, nil
}

// filterClusterImage remove cluster image from images
func filterClusterImage(images []string) (newImages []string) {
	cImageName, cImageID, _ := getRunningClusterImages()
	if cImageName == "" && cImageID == "" {
		return images
	}
	for _, i := range images {
		if i != cImageName && !strings.Contains(cImageID, i) {
			newImages = append(newImages, i)
		} else {
			logrus.Errorf("Image used by %v: image is in use by the cluster", cImageID)
		}
	}
	return
}
