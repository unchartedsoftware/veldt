package terms

import (
	"regexp"
	"strings"
)

var (
	punctRegex   = regexp.MustCompile(`[^\w\s]`)
	urlRegex     = regexp.MustCompile(`(^|\s)((https?:\/\/)?[\w-]+(\.[\w-]+)+\.?(:\d+)?(\/\S*)?)`)
	mentionRegex = regexp.MustCompile(`(@[a-z\d_]+)`)
	hashtagRegex = regexp.MustCompile(`(#[\S\W]+)`)
)

func removePunctuation(text string) string {
	return punctRegex.ReplaceAllString(text, "")
}

func removeURLs(text string) string {
	return urlRegex.ReplaceAllString(text, "")
}

func removeMentions(text string) string {
	return mentionRegex.ReplaceAllString(text, "")
}

func removeHashtags(text string) string {
	return hashtagRegex.ReplaceAllString(text, "")
}

func removeStopWords(words []string) []string {
	validWords := make([]string, len(words))
	i := 0
	for _, word := range words {
		if !IsStopWord(word) {
			validWords[i] = word
			i++
		}
	}
	return validWords[0:i]
}

// ExtractTerms will extract meaningful terms from a string of text.
func ExtractTerms(text string) []string {
	text = removeURLs(text)                                       // remove urls first
	text = removeMentions(text)                                   // then hashtags
	text = removeHashtags(text)                                   // then mentions
	text = removePunctuation(text)                                // finally leftover punctuation
	return removeStopWords(strings.Fields(strings.ToLower(text))) // filter out stopwords
}
