package pkg

import (
	"regexp"
	"strings"
)

func NormalizeVersionToMinor(given string) string {
	//fmt.Println("Normalizing " + given)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	submatchall := re.FindAllString(given, -1)
	if len(submatchall) >= 1 {
		return submatchall[0]
	} else {
		panic("Unable to extract version from " + given)
	}
}

func NormalizeVersionToMinor2(given string) string {
	_, result, _ := strings.Cut(given, "v")
	cutoff := strings.Split(result, ".")
	return cutoff[0] + "." + cutoff[1]
}
