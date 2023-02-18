package main

import (
	"github.com/tonit/rekind/pkg/augment"
	"github.com/tonit/rekind/pkg/images"
	"os"
)

func main() {
	args := os.Args[1:]
	augmenter := []augment.FlagAugmentationInput{
		{Name: "kubeversion", Replace: func(match augment.AugmentationResult) (string, string) {
			return "image", getNodeImage(augment.RunForValue("kind", []string{"version", "-q"}), match.Value)
		}},
		{Name: "withNamespaces", Erase: true, After: func(match augment.AugmentationResult) {}},
		{Name: "image", Erase: true},
	}

	augment.BuildAndRun(augmenter, args, "kind")
}

func getNodeImage(kindVersion string, k8s string) string {
	return images.KindVersions[kindVersion+";"+k8s]
}
