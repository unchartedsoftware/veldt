package query

/*

[
	{
		"equals": {
			"field": "surname",
			"value": "bethune"
		}
	},
	"AND"
	{
		"in": {
			"field": "hashtags",
			"values": ["dank", "420", "nugz"]
		}
	},
	"AND"
	[
		"NOT",
		{
			"exists": {
				"field": "location"
			}
		},
		"OR",
		{
			"range": {
				"field": "data",
				"lt": 123234532452
			}
		}
	]
]

*/

// func ParseBinaryOp(op interface{}) (string, error) {
// 	opString, ok := op.(string)
// 	if !ok {
// 		return "", fmt.Errorf("BinaryExpression operator `%v` is not of type `string`",
// 			op)
// 	}
// 	switch opString {
// 	case And:
// 		return And, nil
// 	case Or:
// 		return Or, nil
// 	default:
// 		return "", fmt.Errorf("BinarryExpression operator `%s` is not recognized",
// 			opString)
// 	}
// }
//
// func ParseExpression(expression []interface{}) {
// 	switch len(expression) {
// 	case 1:
// 		// single query
//
// 	case 2:
// 		// unary boolean expression
//
// 	case 3:
// 		// binary boolean expression
// 	}
// }
//
// func Parse(bytes []byte]) (Query, error) {
//
// 	// Get bytes.
// 	bytes := []byte(text)
//
// 	// Unmarshal JSON to Result struct.
// 	var query []interface{}
// 	json.Unmarshal(bytes, &result)
//
// }
