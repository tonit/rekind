package augment

import (
	execute "github.com/alexellis/go-execute/pkg/v1"
	util "github.com/tonit/rekind/pkg"
	"strings"
)

type AugmentationContext struct {
	memory map[string]string
	steps  map[int]AugmentationResult
}

type FlagAugmentationInput struct {
	Name    string
	Erase   bool
	Append  func(match AugmentationResult) string
	Replace func(match AugmentationResult) (string, string)
	After   func(match AugmentationResult)
}

type AugmentationResult struct {
	Input FlagAugmentationInput
	//Stripped []string
	Tombstone bool
	Start     int
	End       int
	Value     string
}

func BuildAndRun(augmenter []FlagAugmentationInput, args []string, executable string) {
	context := BuildContext(augmenter, args)
	cmd := MakeCommand(args, context)

	command := execute.ExecTask{
		Command:     executable,
		Args:        cmd,
		StreamStdio: true,
	}

	command.Execute()
}

func RunForValue(executable string, args []string) string {

	cmd := execute.ExecTask{
		Command:     executable,
		Args:        args,
		StreamStdio: false,
	}

	res, err := cmd.Execute()
	if err != nil {
		panic(err)
	}

	if res.ExitCode != 0 {
		panic("Non-zero exit code: " + res.Stderr)
	}

	return util.NormalizeVersionToMinor(res.Stdout)
}

func MakeCommand(args []string, context map[int]AugmentationResult) []string {
	var cmd []string
	var appendedStack []string

	for i, arg := range args {
		augmented, exists := context[i]
		if exists {
			if !augmented.Tombstone {
				//fmt.Printf("Found action for %d is %s = %s\n", i, augmented.Input.Name, augmented.Value)
				if augmented.Input.Replace != nil {
					key, value := augmented.Input.Replace(augmented)
					cmd = append(cmd, "--"+key+"="+value)
				} else if augmented.Input.Append != nil {
					appendedStack = append(appendedStack, augmented.Input.Append(augmented))
				}
			} else {
				//fmt.Printf("Tombstone for %d as it is %s\n", i, arg)
			}
		} else {
			cmd = append(cmd, arg)
			//fmt.Printf("We will leave %d as it is %s\n", i, arg)
		}
	}

	for _, a := range appendedStack {
		cmd = append(cmd, a)
	}
	return cmd
}

func BuildContext(augmenter []FlagAugmentationInput, args []string) map[int]AugmentationResult {
	var context = map[int]AugmentationResult{}

	//infix phase
	for _, a := range augmenter {
		thing := extractAugmentOption(a, args)
		if thing.Start >= 0 {
			context[thing.Start] = thing
			// add additional checks for error case. Can only be +1
			if thing.End > thing.Start {
				context[thing.End] = AugmentationResult{
					Tombstone: true,
				}
			}
		}
	}
	return context
}

func extractAugmentOption(searchFor FlagAugmentationInput, args []string) AugmentationResult {
	inFlag := false
	var foundValue string
	var augmentEnd = -1
	var augmentStart = -1

	for i, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--"+searchFor.Name) && strings.Contains(arg, "="):
			foundValue = strings.Split(arg, "=")[1]
			augmentStart = i
			augmentEnd = i
			continue
		case strings.HasPrefix(arg, "--"+searchFor.Name) && !strings.Contains(arg, "="):
			inFlag = true
			augmentStart = i
			augmentEnd = i + 1
			continue
		case inFlag:
			foundValue = arg
			inFlag = false
			continue
		}
	}
	res := AugmentationResult{
		Input: searchFor,
		Value: foundValue,
		Start: augmentStart,
		End:   augmentEnd,
	}
	return res
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}