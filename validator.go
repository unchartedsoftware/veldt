package veldt

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	missing = "???"
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

	req := &TileRequest{}

	v.Buffer("{", 0)

	// validate URI
	req.URI = v.validateURI(args, 1)

	// validate coord
	req.Coord = v.validateCoord(args, 1)

	// validate tile
	req.Tile = v.validateTile(args, 1)

	// validate query
	req.Query = v.validateQuery(args, 1)

	v.Buffer("}", 0)

	// check for any errors
	err := v.Error()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (v *validator) validateMetaRequest(args map[string]interface{}) (*MetaRequest, error) {

	req := &MetaRequest{}

	v.Buffer("{", 0)

	// validate URI
	req.URI = v.validateURI(args, 1)

	// validate tile
	req.Meta = v.validateMeta(args, 1)

	v.Buffer("}", 0)

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

func (v *validator) validateURI(args map[string]interface{}, indent int) string {
	uri, err := v.parseURI(args)
	v.BufferKeyValue("uri", uri, indent, err)
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

func (v *validator) validateCoord(args map[string]interface{}, indent int) *binning.TileCoord {
	params, coord, err := v.parseCoord(args)
	if params != nil {
		v.BufferKeyValue("coord", params, indent, err)
	} else {
		v.BufferKeyValue("coord", missing, indent, err)
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

func (v *validator) validateTile(args map[string]interface{}, indent int) Tile {
	// check if the tile key exists
	arg, ok := args["tile"]
	if !ok {
		v.BufferKeyValue("tile", missing, indent, fmt.Errorf("`tile` not found"))
		return nil
	}

	// check if the tile value is an object
	val, ok := arg.(map[string]interface{})
	if !ok {
		v.BufferKeyValue("tile", arg, indent, fmt.Errorf("`tile` is not of correct type"))
		return nil
	}

	// check if tile is correct
	v.Buffer(`"tile": {`, indent)
	id, params, tile, err := v.parseTile(val)
	if id == "" {
		id = missing
		params = missing
	}
	v.BufferKeyValue(id, params, indent+1, err)
	v.Buffer("}", indent)
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

func (v *validator) validateMeta(args map[string]interface{}, indent int) Meta {
	// check if the meta key exists
	arg, ok := args["meta"]
	if !ok {
		v.BufferKeyValue("meta", missing, indent, fmt.Errorf("`meta` not found"))
		return nil
	}

	// check if the meta value is an object
	val, ok := arg.(map[string]interface{})
	if !ok {
		v.BufferKeyValue("meta", arg, indent, fmt.Errorf("`meta` is not of correct type"))
		return nil
	}

	// check if meta is correct
	v.Buffer(`"meta": {`, indent)
	id, params, meta, err := v.parseMeta(val)
	if id == "" {
		id = missing
		params = missing
	}
	v.BufferKeyValue(id, params, indent+1, err)
	v.Buffer("}", indent)
	return meta
}

func (v *validator) validateQuery(args map[string]interface{}, indent int) Query {
	val, ok := args["query"]
	if !ok {
		return nil
	}
	// nil query is valid
	if val == nil {
		return nil
	}
	// validate the query
	v.Buffer("{", indent)
	validated := v.validateToken(val, indent+1, true)
	v.Buffer("}", indent)
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

func (v *validator) validateQueryToken(args map[string]interface{}, indent int, first bool) Query {
	id, params, query, err := v.parseQuery(args)
	if id == "" {
		id = missing
		params = missing
	}
	if first {
		v.Buffer(`"query": {`, indent)
		v.BufferKeyValue(id, params, indent+1, err)
		v.Buffer("}", indent)
	} else {
		v.BufferKeyValue(id, params, indent, err)
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

func (v *validator) validateOperatorToken(op string, indent int) interface{} {
	if !isValidBoolOperator(op) {
		v.StartError("invalid operator", indent)
		v.Buffer(fmt.Sprintf(`"%v"`, op), indent)
		v.EndError()
		return nil
	}
	v.Buffer(fmt.Sprintf(`"%s"`, op), indent)
	return op
}

func (v *validator) validateExpressionToken(exp []interface{}, indent int, first bool) interface{} {
	// open paren
	if first {
		v.Buffer(`"query": [`, indent)
	} else {
		v.Buffer("[", indent)
	}
	// track last token to ensure next is valid
	var last interface{}
	// for each component
	for i, current := range exp {
		// next line
		if !isTokenValid(last, current) {
			v.StartError("unexpected token", indent+1)
			v.validateToken(current, indent+1, false)
			v.EndError()
			last = current
			continue
		}
		exp[i] = v.validateToken(current, indent+1, false)
		last = current
	}
	// close paren
	v.Buffer("]", indent)
	return exp
}

func (v *validator) validateToken(arg interface{}, indent int, first bool) interface{} {
	// expression
	exp, ok := arg.([]interface{})
	if ok {
		return v.validateExpressionToken(exp, indent, first)
	}
	// query
	query, ok := arg.(map[string]interface{})
	if ok {
		return v.validateQueryToken(query, indent, first)
	}
	// operator
	op, ok := arg.(string)
	if ok {
		return v.validateOperatorToken(op, indent)
	}
	// err
	if first {
		v.BufferKeyValue("query", fmt.Sprintf("%v", arg), indent, fmt.Errorf("`query` is not of correct type"))
	} else {
		v.StartError("unrecognized symbol", indent)
		v.Buffer(fmt.Sprintf("%v", arg), indent)
		v.EndError()
	}
	return arg
}

func getTokenType(token interface{}) string {
	_, ok := token.([]interface{})
	if ok {
		return "exp"
	}
	op, ok := token.(string)
	if ok {
		if isValidBinaryOperator(op) {
			return "binary"
		} else if isValidUnaryOperator(op) {
			return "unary"
		} else {
			return "unrecognized"
		}
	}
	_, ok = token.(map[string]interface{})
	if ok {
		return "query"
	}
	return "unrecognized"
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
	if current == "unrecognized" || next == "unrecognized" {
		// NOTE: consider unrecognized tokens as valid to allow the parsing to
		// continue correctly
		return true
	}
	switch current {
	case "exp":
		return next == "binary"
	case "query":
		return next == "binary"
	case "binary":
		return next == "unary" || next == "query" || next == "exp"
	case "unary":
		return next == "query" || next == "exp"
	}
	return false
}

func firstTokenIsValid(n interface{}) bool {
	next := getTokenType(n)
	switch next {
	case "exp":
		return true
	case "query":
		return true
	case "unary":
		return true
	}
	return false
}
