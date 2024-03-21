package toes

import (
	"github.com/starjun/toes-cli/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var dbDmlTpl = `




func {{StructName}}Create(ctx context.Context, s {{StructName}}) error {
	return global.DB.Create(&s).Error
}

func {{StructName}}Get(ctx context.Context, id int64) (obj {{StructName}}, resault *gorm.DB) {
	resault = global.DB.Where("ID=?", id).Find(&obj)
	return
}

func {{StructName}}Delete(ctx context.Context, id int64, Unscoped bool) error {
	var err error
	if Unscoped {
		// true 使用硬删除
		err = global.DB.Unscoped().Delete(&{{StructName}}{}, "Id=?", id).Error
	} else {
		// 默认软删除
		err = global.DB.Where("Id = ?", id).Delete(&{{StructName}}{}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

// 插入 & 更新数据
func {{StructName}}Update(ctx context.Context, s {{StructName}}) error {
	return global.DB.Save(&s).Error
}

// 仅更新部分字段 如果动态实现
func AccountUpdateExt(ctx context.Context, s {{StructName}}, args ...interface{}) error {
	//return global.DB.Debug().Model(&s).Select(args).Updates(s).Error
	return global.DB.Model(s).Select(args).Updates(s).Error
}

func {{StructName}}ListExt(ctx context.Context, reqParam *QueryConfigRequest) (count int64, ret []{{StructName}}, err error) {
	dbObj := global.DB
	if reqParam.Deleted == 2 {
		// 2 表示查询已删除
		dbObj = dbObj.Unscoped()
	}
	reqParam.MakeGormDbByQueryConfig(dbObj)
	err = dbObj.
		Offset(reqParam.Offset).
		Limit(defaultLimit(reqParam.Limit)).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return count, ret, err
}

// 还需要增加一个常规的 List 函数


`

func (b *ToesGen) TableToModels(tables []string) error {
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
		time.Sleep(1000 * time.Microsecond)
		dbDmlTpl = strings.Replace(dbDmlTpl, "{{at}}", "`", -1)
		dbDmlTpl = strings.Replace(dbDmlTpl, "{{StructName}}", TbNameToStructName(table), -1)
		fPath := modPath + "/" + table + ".gen.go"
		WriteModelsExt(fPath, dbDmlTpl)
	}
	return nil
}

func TbNameToStructName(_tbname string) (filename string) {
	// avoid test file
	filename = _tbname
	for strings.HasSuffix(filename, "_test") {
		pos := strings.LastIndex(filename, "_")
		filename = filename[:pos] + filename[pos+1:]
	}
	return
}

func WriteModelsExt(_path, _msg string) error {
	f, err := os.OpenFile(_path, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return err
	}
	// 关闭文件
	defer f.Close()
	// 字符串写入
	_, err = f.WriteString(_msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
