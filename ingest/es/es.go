package es

import (
    "errors"
    "fmt"
    "strings"
    "net/http"

    "github.com/parnurzeal/gorequest"
)

func Bulk( host string, port string, index string, actions []string ) error {
    jsonLines := fmt.Sprintf( "%s\n", strings.Join( actions, "\n" ) )
	response, err := http.Post( host + ":" + port + "/" + index + "/_bulk", "application/json", strings.NewReader( jsonLines ) )
    if err != nil {
        fmt.Println( err )
        return errors.New("Bulk request failed")
    }
	response.Body.Close()
    return nil
}

func IndexExists( host string, port string, index string ) ( bool, error ) {
    resp, _, errs := gorequest.New().
		Head( host + ":" + port + "/" + index ).
		End()
    if errs != nil {
        fmt.Println( errs )
        return false, errors.New("Unable to determine if index exists")
    }
    return resp.StatusCode != 404, nil
}

func DeleteIndex( host string, port string, index string ) error {
    fmt.Println( "Clearing index '" +  index + "'" )
    _, _, errs := gorequest.New().
        Delete( host + ":" + port + "/" + index ).
        End()
    if errs != nil {
        fmt.Println( errs )
        return errors.New("Failed to delete index")
    }
    return nil
}

func CreateIndex( host string, port string, index string, body string ) error {
    fmt.Println( "Creating index '" + index + "'" )
    _, _, errs := gorequest.New().
        Put( host + ":" + port + "/" + index ).
        Send( body ).
        End()
    if errs != nil {
        fmt.Println( errs )
        return errors.New("Failed to create index")
    }
    return nil
}
