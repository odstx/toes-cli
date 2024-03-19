package toes

import (
	"errors"
	"github.com/starjun/gotools/v2"
	"log"
	"os"
	"path"
	"strings"
)

var job1Tpl = `
package jobs

import "log"

type Job01 struct {
	Test string
	Cnt  int
}

func (g *Job01) Run() {
	log.Println("Hello, ", g.Test, g.Cnt)
	g.Cnt++
}
`

func (b *ToesGen) Job1FileGen() error {
	filepath := path.Join(b.jobsPath, "jobdemo.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
		return e
	}
	job1Tpl = strings.Replace(job1Tpl, "{{at}}", "`", -1)
	f.WriteString(job1Tpl)
	f.Close()
	return nil
}
