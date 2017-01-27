package veldt

// GenerateTile generates a tile for the provided pipeline ID and JSON request.
func GenerateTile(id string, args map[string]interface{}) error {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return err
	}
	req, err := pipeline.NewTileRequest(args)
	if err != nil {
		return err
	}
	return pipeline.GenerateTile(req)
}

// GetTileFromStore retrieves a tile from the store for the provided pipeline ID
// and JSON request.
func GetTileFromStore(id string, args map[string]interface{}) ([]byte, error) {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return nil, err
	}
	req, err := pipeline.NewTileRequest(args)
	if err != nil {
		return nil, err
	}
	return pipeline.GetTileFromStore(req)
}

// GenerateMeta generates meta data for the provided pipeline ID and JSON
// request.
func GenerateMeta(id string, args map[string]interface{}) error {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return err
	}
	req, err := pipeline.NewMetaRequest(args)
	if err != nil {
		return err
	}
	return pipeline.GenerateMeta(req)
}

// GetMetaFromStore retrieves metadata from the store for the provided pipeline
// ID and JSON request.
func GetMetaFromStore(id string, args map[string]interface{}) ([]byte, error) {
	pipeline, err := GetPipeline(id)
	if err != nil {
		return nil, err
	}
	req, err := pipeline.NewMetaRequest(args)
	if err != nil {
		return nil, err
	}
	return pipeline.GetMetaFromStore(req)
}
