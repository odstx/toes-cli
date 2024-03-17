package pkg

import (
	"errors"
	"strings"
)

// GormRule
// contains icontains 包含
// in
// gt gte  大于 大于等于
// lt lte  小于 小于等于.
type GormRule struct {
	Opt        string        `json:"opt"`
	ReStrList  []interface{} `json:"reStrList"`
	Rev        bool          `json:"rev"`
	Lcon       string        `json:"lcon"`
	MaLocation string        `json:"maLocation"`
}

type QueryConfigRequest struct {
	Query   []*GormRule `form:"query" json:"query"`
	Fields  []string    `form:"fields" json:"fields"`
	SortBy  []string    `form:"sortBy" json:"sortBy"`
	Order   []string    `form:"order" json:"order"`
	Limit   int         `form:"limit" json:"limit"`
	Offset  int         `form:"offset" json:"offset"`
	Deleted int8        `form:"deleted" json:"deleted"`
}

func (p *QueryConfigRequest) Check() error {
	for k, v := range p.Query {
		if strings.TrimSpace(v.Opt) == "=" {
			p.Query[k].Opt = "exact"
		}
	}
	for _, val := range p.Query {
		if len(val.ReStrList) == 0 {
			return errors.New("query param error")
		}
		if strings.TrimSpace(val.Lcon) == "" || strings.TrimSpace(val.MaLocation) == "" || strings.TrimSpace(val.Opt) == "" {
			return errors.New("query param error")
		}
	}
	for k, v := range p.Query {
		if strings.ToLower(v.Opt) != "in" {
			p.Query[k].ReStrList = p.Query[k].ReStrList[0:1]
		}
	}

	return nil
}
