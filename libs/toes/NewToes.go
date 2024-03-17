package toes

import (
	"database/sql"
	"github.com/starjun/toes-cli/config"
	"log"
	"os"
	"path"
)

type ToesGen struct {
	PackageName    string `json:"packageName"`    // 主项目名称
	RootPath       string `json:"rootPath"`       // 项目根路径
	ModelPath      string `json:"modelPath"`      // model 路径
	ControllerPath string `json:"controllerPath"` // controller 路径
	RouterPath     string `json:"routerPath"`     // router 路径
	RequestPath    string `json:"requestPath"`    // request 路径
	middlewarePath string `json:"middlewarePath"` // middleware 路径
	libsPath       string `json:"libsPath"`       // libs 路径
	sysinfoPath    string `json:"sysinfoPath"`    // sysinfo 路径
	utilsPath      string `json:"utilsPath"`      // utils 路径

	assetsPath string `json:"assetsPath"` // assets 路径
	buildPath  string `json:"buildPath"`  // bulid 路径
	cmdPath    string `json:"cmdPath"`    // cmd 路径
	confPath   string `json:"confPath"`   // conf 路径
	docsPath   string `json:"docsPath"`   // docs 路径
	globalPath string `json:"globalPath"` // global 路径
	jobsPath   string `json:"jobsPath"`   // jobs 路径
	toolsPath  string `json:"toolsPath"`  // tools 路径
	webPath    string `json:"webPath"`    // web 路径

	Dsn      string `json:"dsn"` // dsn 数据库连接字符串
	Driver   string `json:"driver"`
	Database *sql.DB
	Tables   []string
}

func NewToesGen(packageName, rootPath string) (*ToesGen, error) {
	gen := &ToesGen{
		PackageName:    packageName,
		RootPath:       rootPath,
		ModelPath:      "/internal/models",
		ControllerPath: "/internal/controller",
		RouterPath:     "/internal/routers",
		RequestPath:    "/internal/request",
		middlewarePath: "/internal/middleware",
		libsPath:       "/internal/libs",
		sysinfoPath:    "/internal/sysinfo",
		utilsPath:      "/internal/utils",
		assetsPath:     "assets",
		buildPath:      "build",
		cmdPath:        "cmd",
		confPath:       "conf",
		docsPath:       "docs",
		globalPath:     "global",
		jobsPath:       "jobs",
		toolsPath:      "tools",
		webPath:        "web",
		Driver:         "mysql",
		Dsn:            config.GetDSN(),
		Database:       nil,
	}

	if config.Table != "" {
		gen.SetTables([]string{config.Table})
	}
	db, err := sql.Open(gen.Driver, gen.Dsn)
	if err != nil {
		log.Println(err)

		return nil, err
	}
	gen.Database = db

	return gen, nil
}

func (b *ToesGen) SetPackageName(packageName string) *ToesGen {
	b.PackageName = packageName
	return b
}

func (b *ToesGen) SetRootPath(rootPath string) *ToesGen {
	b.RouterPath = rootPath
	return b
}

func (b *ToesGen) SetModelsPath(modelsPath string) *ToesGen {
	b.ModelPath = modelsPath
	return b
}

func (b *ToesGen) SetTables(tables []string) *ToesGen {
	if len(tables) > 0 {
		b.Tables = append(b.Tables, tables...)
	}
	return b
}

func (b *ToesGen) SetCtrlPath(ctrlPath string) *ToesGen {
	b.ControllerPath = ctrlPath
	return b
}

func (b *ToesGen) SetDriver(driver string) *ToesGen {
	b.Driver = driver
	return b
}

func (b *ToesGen) SetDSN(dsn string) *ToesGen {
	b.Dsn = dsn
	return b
}

// GetTables 获取所有表名
func (b *ToesGen) GetTables() ([]string, error) {
	var tables []string
	rows, err := b.Database.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("Could not show Tables: %s", err)
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Could not show Tables: %s", err)
			return tables, err
		}
		tables = append(tables, name)
	}
	b.Tables = tables

	return tables, nil
}

func (b *ToesGen) GenOthers() error {
	tmpPaths := []string{}
	//	1: assets 目录
	tmpPaths = append(tmpPaths, b.assetsPath)
	// 2: build 目录
	tmpPaths = append(tmpPaths, b.buildPath)
	// 3: cmd 目录
	tmpPaths = append(tmpPaths, b.cmdPath)
	// 4: conf 目录
	tmpPaths = append(tmpPaths, b.confPath)
	// 5: docs 目录
	tmpPaths = append(tmpPaths, b.docsPath)
	// 6: global 目录
	tmpPaths = append(tmpPaths, b.globalPath)
	// 7: jobs 目录
	tmpPaths = append(tmpPaths, b.jobsPath)
	// 8: tools 目录
	tmpPaths = append(tmpPaths, b.toolsPath)
	// 9: web 目录
	tmpPaths = append(tmpPaths, b.webPath)

	for i := 0; i < len(tmpPaths); i++ {
		tmp := path.Join(b.RootPath, tmpPaths[i])
		err := os.MkdirAll(tmp, 0777)
		if err != nil {

			log.Println(tmpPaths[i], err)
			return err
		}
	}
	return nil
}
