package prism

import (
	"encoding/json"
)

func copyJSON(obj map[string]interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = json.Unmarshal(bytes, &copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}

func GenerateTile(id string, args map[string]interface{}) error {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return err
	}
	// params are modified in place, so copy it
	copied, err := copyJSON(args)
	if err != nil {
		return err
	}
	req, err := pipeline.NewTileRequest(copied)
	if err != nil {
		return err
	}
	return pipeline.GenerateTile(req)
}

func GetTileFromStore(id string, args map[string]interface{}) ([]byte, error) {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return nil, err
	}
	// params are modified in place, so copy it
	copied, err := copyJSON(args)
	if err != nil {
		return nil, err
	}
	req, err := pipeline.NewTileRequest(copied)
	if err != nil {
		return nil, err
	}
	return pipeline.GetTileFromStore(req)
}

func GenerateMeta(id string, args map[string]interface{}) error {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return err
	}
	// params are modified in place, so copy it
	copied, err := copyJSON(args)
	if err != nil {
		return err
	}
	req, err := pipeline.NewMetaRequest(copied)
	if err != nil {
		return err
	}
	return pipeline.GenerateMeta(req)
}

func GetMetaFromStore(id string, args map[string]interface{}) ([]byte, error) {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return nil, err
	}
	// params are modified in place, so copy it
	copied, err := copyJSON(args)
	if err != nil {
		return nil, err
	}
	req, err := pipeline.NewMetaRequest(copied)
	if err != nil {
		return nil, err
	}
	return pipeline.GetMetaFromStore(req)
}
