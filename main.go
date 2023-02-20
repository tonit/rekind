package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tonit/rekind/pkg/augment"
	"github.com/tonit/rekind/pkg/docker"
	"github.com/tonit/rekind/pkg/images"
	"os"
	"strings"
)

func main() {
	//color.HiCyan("Using reKinD - an augmented flavour of KinD.")
	var colored = color.HiCyanString("augmented flavour")
	fmt.Println(">> reKinD - " + colored + " of KinD <<")
	//fmt.Println("")

	args := os.Args[1:]

	commands := []augment.CommandAugmentationInput{
		{Name: "get images", Run: func(args []string) {
			var id, err = docker.FindContainer("kind") // use name from flags really...
			if err != nil {
				panic(err)
			}
			var imagesList = docker.ListImages(id) // use name from flags really...
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Image", "Checksum", "Size"})
			t.SetColumnConfigs([]table.ColumnConfig{
				{
					Name:  "Size",
					Align: text.AlignRight,
					//Colors: text.Colors{text.FgHiCyan},
				},
			})

			for _, imgs := range imagesList.Images {
				t.AppendRow([]interface{}{strings.Join(imgs.RepoTags, ","), imgs.ID, imgs.Size})
			}
			t.SetStyle(table.StyleLight)
			t.Style().Color.Header = text.Colors{text.FgHiCyan}
			t.Style().Color.Row = text.Colors{text.FgHiCyan}

			t.Render()

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
