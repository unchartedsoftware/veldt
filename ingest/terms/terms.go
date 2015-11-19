package terms

import (
    "sync"
)

var mutex = sync.Mutex{}
var termCounts = make( map[string]uint64 )
var savedTopTerms map[string]bool

func AddTerms( text string ) {
    terms := ExtractTerms( text )
    mutex.Lock()
    for _, term := range terms {
        termCounts[ term ]++
    }
    mutex.Unlock()
}

func SaveTopTerms( num uint64 ) {
    mutex.Lock()
    savedTopTerms = GetTopTermsMap( num )
    mutex.Unlock()
}

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
