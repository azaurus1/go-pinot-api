package model

type GetSchemaResponse []string

func (g *GetSchemaResponse) ForEachSchema(apply func(schemaName string)) {
	for _, v := range *g {
		apply(v)
	}
}
