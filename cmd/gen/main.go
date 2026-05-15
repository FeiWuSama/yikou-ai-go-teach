package main

import (
	"gorm.io/gen"
	"yikou-ai-go-teach/config"
	"yikou-ai-go-teach/internal/dal"
)

func main() {
	initConfig := config.InitConfig()
	db := dal.InitDB(initConfig)
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/dal/query",
		ModelPkgPath: "model",
		Mode: gen.WithoutContext |
			gen.WithDefaultQuery |
			gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable(
		gen.FieldJSONTag("id", "id,string"),
	)...)

	g.Execute()
}
