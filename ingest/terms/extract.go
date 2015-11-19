package terms

import (
    "regexp"
    "strings"
)

func removePunctuation( text string ) string {
    reg, _ := regexp.Compile(`[\W]`)
    return reg.ReplaceAllString( text, "" )
}

func removeURLs( text string ) string {
    reg, _ := regexp.Compile(`(^|\s)((https?:\/\/)?[\w-]+(\.[\w-]+)+\.?(:\d+)?(\/\S*)?)`)
    return reg.ReplaceAllString( text, "" )
}

func removeMentions( text string ) string {
    reg, _ := regexp.Compile(`(@[a-z\d_]+)`)
    return reg.ReplaceAllString( text, "" )
}

func removeHashtags( text string ) string {
    reg, _ := regexp.Compile(`(#[\S\W]+)`)
    return reg.ReplaceAllString( text, "" )
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
    return removeStopWords( strings.Fields( text ) ) // filter out stopwords
}
