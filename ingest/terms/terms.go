package terms

import (
    "regexp"
    "sort"
    "strings"
    "sync"
)

// Term class to store a term and the occurance count
type Term struct {
    Text string
    Count uint64
}

// Terms sorting interface
type Terms []Term
func (terms Terms) Len() int {
    return len( terms )
}
func (terms Terms) Less(i, j int) bool {
    return terms[i].Count < terms[j].Count
}
func (terms Terms) Swap(i, j int) {
    terms[i], terms[j] = terms[j], terms[i]
}

var mutex = sync.Mutex{}
var termCounts = make( map[string]uint64 )

func ClearTerms() {
    termCounts = make( map[string]uint64 )
}

func RemovePunctuation( text string ) string {
    reg, _ := regexp.Compile("[.,-/#!$%^&*;:{}=-_`~()]")
    return reg.ReplaceAllString( text, "" )
}

func CountTerms( text string ) {
    words := strings.Fields( RemovePunctuation( text ) )
    mutex.Lock()
    for _, word := range words {
        if !IsStopWord( word ) {
            termCounts[ word ]++
        }
    }
    mutex.Unlock()
}

func GetTermCounts( num uint64 ) []Term {
    terms := make( Terms, len( termCounts ) )
    i := 0
    for term, count := range termCounts {
        terms[i] = Term{
            Text: term,
            Count: count,
        }
        i++
    }
    sort.Sort( sort.Reverse( terms ) )
    numTerms := uint64(len(terms))
    if num > numTerms {
        return terms[0:numTerms]
    }
    return terms[0:num]
}

func GetTopTerms( num uint64 ) []string {
    terms := GetTermCounts( num )
    topTerms := make( []string, len( termCounts ) )
    i := 0
    for _, term := range terms {
        topTerms[i] = term.Text
        i++
    }
    return topTerms[0:]
}

func GetTopTermsMap( num uint64 ) map[string]bool {
    topTerms := GetTopTerms( num )
    topTermsMap := make( map[string]bool )
    for _, term := range topTerms {
        topTermsMap[term] = true
    }
    return topTermsMap
}
