package prism

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
