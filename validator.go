package prism

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Validator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type Validator struct {
	json.Validator
	pipeline *Pipeline
}

// NewValidator instantiates and returns a new query expression object.
func NewValidator(pipeline *Pipeline) *Validator {
	v := &Validator{
		pipeline: pipeline,
	}
	v.Output = make([]string, 0)
	return v
}

func (v *Validator) ValidateTileRequest(args map[string]interface{}) (*TileRequest, error) {

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

func (v *Validator) ValidateMetaRequest(args map[string]interface{}) (*MetaRequest, error) {

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

// parses the tile request JSON for the provided URI.
//
// Ex:
//     {
//         "uri": "example-uri-value0"
//     }
//
func (v *Validator) parseURI(args map[string]interface{}) (string, error) {
	val, ok := args["uri"]
	if !ok {
		return "???", fmt.Errorf("`uri` not found")
	}
	uri, ok := val.(string)
	if !ok {
		return v.FormatVal(val), fmt.Errorf("`uri` not of type `string`")
	}
	return uri, nil
}

func (v *Validator) validateURI(args map[string]interface{}, indent int) string {
	uri, err := v.parseURI(args)
	v.BufferKeyValue("uri", uri, indent, err)
	return uri
}

// parses the tile request JSON for the provided tile coordinate.
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
func (v *Validator) parseCoord(args map[string]interface{}) (*binning.TileCoord, error) {
	c, ok := args["coord"]
	coord, ok := c.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`coord` not found")
	}
	ix, ok := coord["x"]
	if !ok {
		return nil, fmt.Errorf("`coord.x` not found")
	}
	x, ok := ix.(float64)
	if !ok {
		return nil, fmt.Errorf("`coord.x` is not of type `number`")
	}
	iy, ok := coord["y"]
	if !ok {
		return nil, fmt.Errorf("`coord.y` not found")
	}
	y, ok := iy.(float64)
	if !ok {
		return nil, fmt.Errorf("`coord.y` is not of type `number`")
	}
	iz, ok := coord["z"]
	if !ok {
		return nil, fmt.Errorf("`coord.z` not found")
	}
	z, ok := iz.(float64)
	if !ok {
		return nil, fmt.Errorf("`coord.z` is not of type `number`")
	}
	return &binning.TileCoord{
		X: uint32(x),
		Y: uint32(y),
		Z: uint32(z),
	}, nil
}

func (v *Validator) validateCoord(args map[string]interface{}, indent int) *binning.TileCoord {
	coord, err := v.parseCoord(args)
	v.BufferKeyValue("coord", args, indent, err)
	return coord
}

// parses the tile request JSON for the provided tile coordinate.
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
func (v *Validator) parseTile(args map[string]interface{}) (string, map[string]interface{}, Tile, error) {
	arg, ok := args["tile"]
	if !ok {
		return "", nil, nil, fmt.Errorf("`tile` not found")
	}
	val, ok := arg.(map[string]interface{})
	if !ok {
		return "", nil, nil, fmt.Errorf("`tile` does not contain parameters")
	}
	id, params, err := v.GetIDAndParams(val)
	if err != nil {
		return id, params, nil, fmt.Errorf("unable to parse `tile` parameters")
	}
	tile, err := v.pipeline.GetTile(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, tile, nil
}

func (v *Validator) validateTile(args map[string]interface{}, indent int) Tile {
	_, _, tile, err := v.parseTile(args)
	v.BufferKeyValue("tile", args, indent, err)
	return tile
}

func (v *Validator) validateQuery(args map[string]interface{}, indent int) Query {
	val, ok := args["query"]
	if !ok {
		return nil
	}
	validated := v.validateToken(val, indent+1)
	query, err := parseQueryExpression(validated)
	if err != nil {
		return nil
	}
	return query
}

// parses the tile request JSON for the provided tile coordinate.
//
// Ex:
//     {
//         "range": {
//              "field": "age",
//              "gte": 19
//         }
//     }
//
func (v *Validator) parseQuery(args map[string]interface{}) (string, map[string]interface{}, Query, error) {
	id, params, err := v.GetIDAndParams(args)
	if err != nil {
		return id, params, nil, fmt.Errorf("unable to parse `query` parameters")
	}
	query, err := v.pipeline.GetQuery(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, query, nil
}

func (v *Validator) validateQueryToken(args map[string]interface{}, indent int) Query {
	id, params, query, err := v.parseQuery(args)
	v.BufferKeyValue(id, params, indent, err)
	return query
}

func (v *Validator) validateOperatorToken(op string, indent int) interface{} {
	if !IsBoolOperator(op) {
		v.StartError("invalid operator", indent)
		v.Buffer(fmt.Sprintf("\"%v\"", op), indent)
		v.EndError()
		return nil
	}
	v.Buffer(fmt.Sprintf("\"%s\"", op), indent)
	return op
}

func (v *Validator) validateExpressionToken(exp []interface{}, indent int) interface{} {
	// open paren
	v.Buffer("[", indent)
	// track last token to ensure next is valid
	var last interface{}
	// for each component
	for i, sub := range exp {
		// next line
		if last != nil {
			if !nextTokenIsValid(last, sub) {
				v.StartError("unexpected token", indent+1)
				v.validateToken(sub, indent+1)
				v.EndError()
				last = sub
				continue
			}
		}
		exp[i] = v.validateToken(sub, indent+1)
		last = sub
	}
	// close paren
	v.Buffer("]", indent)
	return exp
}

func (v *Validator) validateToken(arg interface{}, indent int) interface{} {
	// expression
	exp, ok := arg.([]interface{})
	if ok {
		return v.validateExpressionToken(exp, indent)
	}
	// query
	query, ok := arg.(map[string]interface{})
	if ok {
		return v.validateQueryToken(query, indent)
	}
	// operator
	op, ok := arg.(string)
	if ok {
		return v.validateOperatorToken(op, indent)
	}
	// err
	v.StartError("unrecognized symbol", indent)
	v.Buffer(fmt.Sprintf("\"%v\"", arg), indent)
	v.EndError()
	return arg
}

func getTokenType(token interface{}) string {
	_, ok := token.([]interface{})
	if ok {
		return "exp"
	}
	op, ok := token.(string)
	if ok {
		if IsBinaryOperator(op) {
			return "binary"
		} else if IsUnaryOperator(op) {
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

// parses the tile request JSON for the provided tile coordinate.
//
// Ex:
//     {
//         "meta": {
//             "default": {}
//         }
//     }
//
func (v *Validator) parseMeta(args map[string]interface{}) (string, map[string]interface{}, Meta, error) {
	arg, ok := args["meta"]
	if !ok {
		return "", nil, nil, fmt.Errorf("`meta` not found")
	}
	val, ok := arg.(map[string]interface{})
	if !ok {
		return "", nil, nil, fmt.Errorf("`meta` does not contain parameters")
	}
	id, params, err := v.GetIDAndParams(val)
	if err != nil {
		return id, params, nil, fmt.Errorf("unable to parse `meta` parameters")
	}
	meta, err := v.pipeline.GetMeta(id, params)
	if err != nil {
		return id, params, nil, err
	}
	return id, params, meta, nil
}

func (v *Validator) validateMeta(args map[string]interface{}, indent int) Meta {
	_, _, meta, err := v.parseMeta(args)
	v.BufferKeyValue("meta", args, indent, err)
	return meta
}
