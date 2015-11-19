package terms

import (
    "regexp"
    "strings"
)

var PunctRegex, _ = regexp.Compile(`[^\w\s]`)
var URLRegex, _ = regexp.Compile(`(^|\s)((https?:\/\/)?[\w-]+(\.[\w-]+)+\.?(:\d+)?(\/\S*)?)`)
var MentionRegex, _ = regexp.Compile(`(@[a-z\d_]+)`)
var HashtagRegex, _ = regexp.Compile(`(#[\S\W]+)`)

func removePunctuation( text string ) string {
    return PunctRegex.ReplaceAllString( text, "" )
}

func removeURLs( text string ) string {
    return URLRegex.ReplaceAllString( text, "" )
}

func removeMentions( text string ) string {
    return MentionRegex.ReplaceAllString( text, "" )
}

func removeHashtags( text string ) string {
    return HashtagRegex.ReplaceAllString( text, "" )
}

func removeStopWords( words []string ) []string {
    validWords := make( []string, len(words) )
    i := 0
    for _, word := range words {
        if !IsStopWord( word ) {
            validWords[i] = word
            i++
        }
    }
    return validWords[0:i]
}

func ExtractTerms( text string ) []string {
    text = removeURLs( text ) // remove urls first
    text = removeMentions( text ) // then hashtags
    text = removeHashtags( text ) // then mentions
    text = removePunctuation( text ) // finally leftover punctuation
    return removeStopWords( strings.Fields( strings.ToLower( text ) ) ) // filter out stopwords
}
