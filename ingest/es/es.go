package es

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/util/log"
)

// Error represents the error node of an elasticsearch json response.
type Error struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// ItemIndex represents the index node of an elasticsearch item attribute.
type ItemIndex struct {
	Error Error `json:"error"`
}

// Item represents an item element of an elasticsearch items array.
type Item struct {
	Index ItemIndex `json:"index"`
}

// Response represents the json response from a valid elasticsearch action.
// {
//     "took": 5,
//     "errors": true,
//     "items": [
//         {
//             "index": {
//                 "_index": "isil_twitter",
//                 "_type": "datum",
//                 "_id": "287366349146181633",
//                 "status": 400,
//                 "error": {
//                     "type": "mapper_parsing_exception",
//                     "reason": "failed to parse",
//           		   "caused_by": {
//                         "type": "json_parse_exception",
//                         "reason": "..."
//                     }
//                 }
//             }
//         }
//     ]
// }
type Response struct {
	Errors bool   `json:"errors"`
	Items  []Item `json:"items"`
}

// BadRequest represents the json response from an invalid elasticsearch action.
// {
//     "error": {
// 	        "root_cause": [
// 	            {
// 	        	    "type": "action_request_validation_exception",
// 	            	"reason": "Validation Failed: 1: no requests added;"
// 	            }
// 	        ],
// 	        "type": "action_request_validation_exception",
// 	        "reason": "Validation Failed: 1: no requests added;"
// 	   },
//     "status": 400
// }
type BadRequest struct {
	Status uint  `json:"status"`
	Error  Error `json:"error"`
}

func getResponseString(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func parseResponse(r *http.Response) error {
	respText, err := getResponseString(r)
	if err != nil {
		return err
	}
	// first check if the request itself was malformed
	if r.StatusCode >= 400 && r.StatusCode < 500 {
		badReq := &BadRequest{}
		err := json.Unmarshal([]byte(respText), &badReq)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %d: %s, %s", r.StatusCode, badReq.Error.Type, badReq.Error.Reason)
	}
	// then check elasticsearch response json
	// unmarshal payload
	esResp := &Response{}
	err = json.Unmarshal([]byte(respText), &esResp)
	if err != nil {
		return err
	}
	if esResp.Errors {
		item := esResp.Items[0]
		log.Error(respText)
		return errors.New(item.Index.Error.Type + ": " + item.Index.Error.Reason)
	}
	return nil
}

// Bulk sends a bulk request to elasticsearch with the provided payload.
func Bulk(host string, port string, index string, datatype string, actions []string) error {
	jsonLines := fmt.Sprintf("%s\n", strings.Join(actions, "\n"))
	resp, err := http.Post(host+":"+port+"/"+index+"/"+datatype+"/_bulk", "application/json", strings.NewReader(jsonLines))
	if err != nil {
		return err
	}
	err = parseResponse(resp)
	if err != nil {
		return err
	}
	return nil
}

// IndexExists returns whether or not the provided index exists in elasticsearch.
func IndexExists(host string, port string, index string) (bool, error) {
	resp, _, errs := gorequest.New().
		Head(host + ":" + port + "/" + index).
		End()
	if errs != nil {
		log.Error(errs)
		return false, errors.New("Unable to determine if index exists")
	}
	return resp.StatusCode != 404, nil
}

// DeleteIndex deletes the provided index in elasticsearch.
func DeleteIndex(host string, port string, index string) error {
	log.Debug("Clearing index '" + index + "'")
	_, _, errs := gorequest.New().
		Delete(host + ":" + port + "/" + index).
		End()
	if errs != nil {
		log.Error(errs)
		return errors.New("Failed to delete index")
	}
	return nil
}

// CreateIndex creates the provided index in elasticsearch.
func CreateIndex(host string, port string, index string, body string) error {
	log.Debug("Creating index '" + index + "'")
	_, _, errs := gorequest.New().
		Put(host + ":" + port + "/" + index).
		Send(body).
		End()
	if errs != nil {
		log.Error(errs)
		return errors.New("Failed to create index")
	}
	return nil
}

// PrepareIndex will ensure the provided index exists, and will optionally clear it.
func PrepareIndex(host string, port string, index string, mappings string, clearExisting bool) error {
	// check if index exists
	indexExists, err := IndexExists(host, port, index)
	if err != nil {
		return err
	}
	// if index exists
	if indexExists && clearExisting {
		err = DeleteIndex(host, port, index)
		if err != nil {
			return err
		}
	}
	// if index does not exist at this point
	if !indexExists || clearExisting {
		err = CreateIndex(
			host,
			port,
			index,
			`{
                "mappings": `+mappings+`
            }`)
		if err != nil {
			return err
		}
	}
	return nil
}
