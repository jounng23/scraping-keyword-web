package strutil

import (
	"regexp"
	"strconv"
)

func CollectTotalSearchResultsFromStats(stats string) int {
	// Define a regular expression to match the number pattern
	re := regexp.MustCompile(`\d+`)

	// Find all matches of the number pattern in the text
	matches := re.FindAllString(stats, -1)

	// Join the matched digits to form the complete number string
	numericText := ""
	for _, match := range matches {
		numericText += match
	}

	// Convert the numeric string to int64
	number, _ := strconv.Atoi(numericText)
	return number
}
