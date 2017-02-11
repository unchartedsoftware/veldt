package veldt

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	missing     = "???"
	expType     = iota
	queryType   = iota
	binaryType  = iota
	unaryType   = iota
	invalidType = iota
)

// validator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type validator struct {
	json.Validator
	pipeline *Pipeline
}

func newValidator(pipeline *Pipeline) *validator {
	v := &validator{
		pipeline: pipeline,
	}
	return v
}

func (v *validator) validateTileRequest(args map[string]interface{}) (*TileRequest, error) {

	v.StartObject()

	req := &TileRequest{}

	// validate URI
	req.URI = v.validateURI(args)

	// validate coord
	req.Coord = v.validateCoord(args)

	// validate tile
	req.Tile = v.validateTile(args)

	// validate query
	req.Query = v.validateQuery(args)

	v.EndObject()

	// check for any errors
	err := v.Error()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (v *validator) validateMetaRequest(args map[string]interface{}) (*MetaRequest, error) {

	v.StartObject()

	req := &MetaRequest{}

	// validate URI
	req.URI = v.validateURI(args)

	// validate meta
	req.Meta = v.validateMeta(args)

	v.EndObject()

	// check for any errors
	err := v.Error()
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Parses the tile request JSON for the provided URI.
//
// Ex:
//     {
//         "uri": "example-uri-value0"
//     }
//
func (v *validator) parseURI(args map[string]interface{}) (string, error) {
	val, ok := args["uri"]
	if !ok {
		return missing, fmt.Errorf("`uri` not found")
	}
	uri, ok := val.(string)
	if !ok {
		return fmt.Sprintf("%v", val), fmt.Errorf("`uri` not of type `string`")
	}
	return uri, nil
}

func (v *validator) validateURI(args map[string]interface{}) string {
	uri, err := v.parseURI(args)
	v.BufferKeyValue("uri", uri, err)
	return uri
}

// Parses the tile request JSON for the provided tile coordinate.
//
// Ex:
//     {
//         "coord": {
//             "z": 4,
//             "x": 12,
//             "y": 3,
//         }
//     }
//
func (v *validator) parseCoord(args map[string]interface{}) (interface{}, *binning.TileCoord, error) {
	c, ok := args["coord"]
	if !ok {
		return nil, nil, fmt.Errorf("`coord` not found")
	}
	coord, ok := c.(map[string]interface{})
	if !ok {
		return c, nil, fmt.Errorf("`coord` is not of correct type")
	}
	ix, ok := coord["x"]
	if !ok {
		return coord, nil, fmt.Errorf("`coord.x` not found")
	}
	x, ok := ix.(float64)
	if !ok {
		return coord, nil, fmt.Errorf("`coord.x` is not of type `number`")
	}
	iy, ok := coord["y"]
	if !ok {
		return coord, nil, fmt.Errorf("`coord.y` not found")
	}
	y, ok := iy.(float64)
	if !ok {
		return coord, nil, fmt.Errorf("`coord.y` is not of type `number`")
	}
	iz, ok := coord["z"]
	if !ok {
		return coord, nil, fmt.Errorf("`coord.z` not found")
	}
	z, ok := iz.(float64)
	if !ok {
		return coord, nil, fmt.Errorf("`coord.z` is not of type `number`")
	}
	return coord, &binning.TileCoord{
		X: uint32(x),
		Y: uint32(y),
		Z: uint32(z),
	}, nil
}

func (v *validator) validateCoord(args map[string]interface{}) *binning.TileCoord {
	params, coord, err := v.parseCoord(args)
	if params != nil {
		v.BufferKeyValue("coord", params, err)
	} else {
		v.BufferKeyValue("coord", missing, err)
	}
	return coord
}

// Parses the tile request JSON for the provided tile type and parameters.
//
// Ex:
//     {
//         "tile": {
//             "heatmap": {
//                  "xField": "pixel.x",
//                  "yField": "pixel.y",
//                  "left": 0,
//                  "right": 4294967296,
//                  "bottom": 0,
//                  "top": 4294967296,
//                  "resolution": 256
//             }
//         }
//     }
//
func (v *validator) parseTile(args map[string]interface{}) (string, interface{}, Tile, error) {
	id, params, err := v.GetIDAndParams(args)
	if err != nil {
		return id, params, nil, err
	}
	tile, err := v.pipeline.GetTile(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, tile, nil
}

func (v *validator) validateTile(args map[string]interface{}) Tile {
	// check if the tile key exists
	arg, ok := args["tile"]
	if !ok {
		v.BufferKeyValue("tile", missing, fmt.Errorf("`tile` not found"))
		return nil
	}

	// check if the tile value is an object
	val, ok := arg.(map[string]interface{})
	if !ok {
		v.BufferKeyValue("tile", arg, fmt.Errorf("`tile` is not of correct type"))
		return nil
	}

	// check if tile is correct
	v.StartSubObject("tile")
	id, params, tile, err := v.parseTile(val)
	if id == "" {
		id = missing
		params = missing
	}
	v.BufferKeyValue(id, params, err)
	v.EndObject()
	return tile
}

// Parses the meta request JSON for the provided meta type and parameters.
//
// Ex:
//     {
//         "meta": {
//             "default": {}
//         }
//     }
//
func (v *validator) parseMeta(args map[string]interface{}) (string, interface{}, Meta, error) {
	id, params, err := v.GetIDAndParams(args)
	if err != nil {
		return id, params, nil, err
	}
	tile, err := v.pipeline.GetMeta(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, tile, nil
}

func (v *validator) validateMeta(args map[string]interface{}) Meta {
	// check if the meta key exists
	arg, ok := args["meta"]
	if !ok {
		v.BufferKeyValue("meta", missing, fmt.Errorf("`meta` not found"))
		return nil
	}

	// check if the meta value is an object
	val, ok := arg.(map[string]interface{})
	if !ok {
		v.BufferKeyValue("meta", arg, fmt.Errorf("`meta` is not of correct type"))
		return nil
	}

	// check if meta is correct
	v.StartSubObject("meta")
	id, params, meta, err := v.parseMeta(val)
	if id == "" {
		id = missing
		params = missing
	}
	v.BufferKeyValue(id, params, err)
	v.EndObject()
	return meta
}

func (v *validator) validateQuery(args map[string]interface{}) Query {
	val, ok := args["query"]
	if !ok {
		return nil
	}
	// nil query is valid
	if val == nil {
		return nil
	}
	// validate the query
	v.StartObject()
	validated := v.validateToken(val, true)
	v.EndObject()
	// parse the expression
	query, err := newExpressionParser(v.pipeline).Parse(validated)
	if err != nil {
		return nil
	}
	return query
}

// Parses the query request JSON for the provided query expression.
//
// Ex:
//     {
//         "range": {
//              "field": "age",
//              "gte": 19
//         }
//     }
//
func (v *validator) parseQuery(args map[string]interface{}) (string, interface{}, Query, error) {
	id, params, err := v.GetIDAndParams(args)
	if err != nil {
		return id, params, nil, err
	}
	query, err := v.pipeline.GetQuery(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, query, nil
}

func (v *validator) validateQueryToken(args map[string]interface{}, first bool) Query {
	id, params, query, err := v.parseQuery(args)
	if id == "" {
		id = missing
		params = missing
	}
	if first {
		v.StartSubObject("query")
		v.BufferKeyValue(id, params, err)
		v.EndObject()
	} else {
		v.BufferKeyValue(id, params, err)
	}
	return query
}

func isValidBinaryOperator(op string) bool {
	return op == And || op == Or
}

func isValidUnaryOperator(op string) bool {
	return op == Not
}

func isValidBoolOperator(op string) bool {
	return isValidBinaryOperator(op) || isValidUnaryOperator(op)
}

func (v *validator) validateOperatorToken(op string) interface{} {
	if !isValidBoolOperator(op) {
		v.BufferValue(op, fmt.Errorf("invalid operator"))
		return nil
	}
	v.BufferValue(op, nil)
	return op
}

func (v *validator) validateExpressionToken(exp []interface{}, first bool) interface{} {
	// open paren
	if first {
		v.StartSubArray("query")
	} else {
		v.StartArray()
	}
	// track last token to ensure next is valid
	var last interface{}
	// for each component
	for i, current := range exp {
		// next line
		if !isTokenValid(last, current) {
			v.StartError("unexpected token")
			v.validateToken(current, false)
			v.EndError()
			last = current
			continue
		}
		exp[i] = v.validateToken(current, false)
		last = current
	}
	// close paren
	v.EndArray()
	return exp
}

func (v *validator) validateToken(arg interface{}, first bool) interface{} {
	// expression
	exp, ok := arg.([]interface{})
	if ok {
		return v.validateExpressionToken(exp, first)
	}
	// query
	query, ok := arg.(map[string]interface{})
	if ok {
		return v.validateQueryToken(query, first)
	}
	// operator
	op, ok := arg.(string)
	if ok {
		return v.validateOperatorToken(op)
	}
	// err
	if first {
		v.BufferKeyValue("query", fmt.Sprintf("%v", arg), fmt.Errorf("`query` is not of correct type"))
	} else {
		v.BufferValue(arg, fmt.Errorf("unrecognized symbol"))
	}
	return arg
}

func getTokenType(token interface{}) int {
	_, ok := token.([]interface{})
	if ok {
		return expType
	}
	op, ok := token.(string)
	if ok {
		if isValidBinaryOperator(op) {
			return binaryType
		} else if isValidUnaryOperator(op) {
			return unaryType
		} else {
			return invalidType
		}
	}
	_, ok = token.(map[string]interface{})
	if ok {
		return queryType
	}
	return invalidType
}

func isTokenValid(c interface{}, n interface{}) bool {
	if c == nil {
		return firstTokenIsValid(n)
	}
	return nextTokenIsValid(c, n)
}

func nextTokenIsValid(c interface{}, n interface{}) bool {
	current := getTokenType(c)
	next := getTokenType(n)
	if current == invalidType || next == invalidType {
		// NOTE: consider unrecognized tokens as valid to allow the parsing to
		// continue correctly
		return true
	}
	switch current {
	case expType:
		return next == binaryType
	case queryType:
		return next == binaryType
	case binaryType:
		return next == unaryType || next == queryType || next == expType
	case unaryType:
		return next == queryType || next == expType
	}
	return false
}

func firstTokenIsValid(n interface{}) bool {
	next := getTokenType(n)
	switch next {
	case expType:
		return true
	case queryType:
		return true
	case unaryType:
		return true
	}
	return false
}
