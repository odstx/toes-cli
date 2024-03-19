package toes

import (
	"github.com/starjun/toes-cli/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"path"
)

func (b *ToesGen) TableToModel(tables []string) error {
	modPath := path.Join(b.RootPath, b.ModelPath)
	os.MkdirAll(modPath, 0777)
	g := gen.NewGenerator(gen.Config{
		OutPath:      modPath,
		ModelPkgPath: "models",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	db, err := gorm.Open(mysql.Open(config.GetDSN()))
	if err != nil {
		return err
	}
	g.UseDB(db)
	for _, table := range tables {
		g.GenerateModel(table)
		g.Execute()
	}
	return nil
}
