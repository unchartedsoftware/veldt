package terms

import (
    "sync"
)

var mutex = sync.Mutex{}
var termCounts = make( map[string]uint64 )
var savedTopTerms map[string]bool

// AddTerms adds the terms of the provided string to the current term counts.
func AddTerms( text string ) {
    terms := ExtractTerms( text )
    mutex.Lock()
    for _, term := range terms {
        termCounts[ term ]++
    }
    mutex.Unlock()
}

// SaveTopTerms saves the top N terms of the current term count map.
func SaveTopTerms( num uint64 ) {
    mutex.Lock()
    savedTopTerms = GetTopTermsMap( num )
    mutex.Unlock()
}

// GetTopTerms returns matching terms in the provided string that are in the saved terms map.
func GetTopTerms( text string ) []string {
    terms := ExtractTerms( text )
    topTerms := make( []string, len(terms) )
    i := 0
    for _, term := range terms {
        if savedTopTerms[term] {
            topTerms[i] = term
            i++
        }
    }
    return topTerms[0:i]
}
