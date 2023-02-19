package main

import (
	"github.com/fatih/color"
	"github.com/tonit/rekind/pkg/augment"
	"github.com/tonit/rekind/pkg/docker"
	"github.com/tonit/rekind/pkg/images"
	"os"
	"strings"
)

func main() {
	color.Magenta("Using reKinD - an augmented flavour of KinD.")
	args := os.Args[1:]

	commands := []augment.CommandAugmentationInput{
		{Name: "get images", Run: func(args []string) {
			var id, err = docker.FindContainer("kind") // use name from flags really...
			if err != nil {
				panic(err)
			}
			var images = docker.ListImages(id) // use name from flags really...
			for _, imgs := range images.Images {
				color.Magenta(imgs.ID + "," + strings.Join(imgs.RepoTags, ",") + "," + imgs.Size)
			}

		}},
		{Name: "", Run: func(args []string) {
			color.Magenta("Running augmented: kind " + strings.Join(args, " "))
			augment.OneOffCommand("kind", args)
		}},
	}

	augmenter := []augment.FlagAugmentationInput{
		{Name: "kubeversion", Replace: func(match augment.AugmentationResult) (string, string) {
			return "image", getNodeImage(augment.RunForValue("kind", []string{"version", "-q"}), match.Value)
		}},
		{Name: "withNamespaces", Erase: true, After: func(match augment.AugmentationResult) {}},
		{Name: "image", Erase: true},
		{Name: "name", Erase: false}, // TODO: will copy it to the context
	}

	augment.BuildAndRun(commands, augmenter, args)
}

func getNodeImage(kindVersion string, k8s string) string {
	return images.KindVersions[kindVersion+";"+k8s]
}
