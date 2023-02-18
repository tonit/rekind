package augment

import "testing"

func TestAdd(t *testing.T) {
	input := FlagAugmentationInput{
		Name: "image",
	}
	extractAugmentOption(input, []string{"create", "cluster", "--image=2"})

}
