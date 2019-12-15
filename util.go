package main

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	pl "github.com/gertd/go-pluralize"
)

var pluralizeClient *pl.Client

func init() {
	rand.Seed(time.Now().UnixNano())
	pluralizeClient = pl.NewClient()
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func resolveOptions(text string) string {
	if strings.Contains(text, "[") {
		options := getTextBetweenTags(text, "[", "]")
		option := pickRandomFromList(strings.Split(options, ","))
		text = replaceTextBetweenTags(text, option, "[", "]")
		// recursively fill in all options
		return resolveOptions(text)
	}
	return text
}

func getTextBetweenTags(text, startTag, endTag string) string {
	return strings.SplitN(strings.SplitN(text, startTag, 2)[1], endTag, 2)[0]
}

func replaceTextBetweenTags(text, replacement, startTag, endTag string) string {
	startIndex := strings.Index(text, startTag)
	endIndex := startIndex + strings.Index(text[startIndex+1:], endTag)
	return text[:startIndex] + replacement + text[endIndex+2:]
}

func pickRandomFromList(list []string) string {
	return list[rand.Int31n(int32(len(list)))]
}

func randomChance(probability float32) bool {
	return rand.Float32() < probability
}

func indefiniteArticle(word string) string {
	word = strings.TrimSpace(word)

	// exceptions:
	if strings.HasPrefix(word, "one") || strings.HasPrefix(word, "uni") {
		return "a"
	}

	//return 'an' if word starts with vowel, otherwise 'a'
	vowels := "aeiou"
	if strings.Contains(vowels, word[:1]) {
		return "an"
	}

	return "a"
}

func pluralize(word string) string {
	return pluralizeClient.Plural(word)
}

func formatOutput(result string) string {
	result = strings.TrimSpace(result)
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ") // replace multiple spaces with single space
	result = strings.ReplaceAll(result, " -", "-")                   // remove accidental space between hyphenated words
	result = strings.ReplaceAll(result, "- ", "-")
	result = strings.ReplaceAll(result, " ,", ",")
	result = strings.ToUpper(result[:1]) + result[1:]
	if strings.HasSuffix(result, ",") {
		result = result[:len(result)-1]
	}
	result = result + "."
	return result
}
