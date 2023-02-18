package augment

import "testing"

func TestAdd(t *testing.T) {
	input := FlagAugmentationInput{
		Name: "image",
	}
	result := extractAugmentOption(input, []string{"create", "cluster", "--image=2"})
	if result.Start != 1 {
		t.Fail()

	}
}
