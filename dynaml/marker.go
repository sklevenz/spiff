package dynaml

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/spiff/yaml"
)

const (
	TEMPORARY = "&temporary"
	TEMPLATE  = "&template"
)

type MarkerExpr struct {
	list []string
	expr Expression
}

func (e MarkerExpr) String() string {
	if e.expr != nil {
		return fmt.Sprintf("%s (%s)", strings.Join(e.list, " "), e.expr)
	}
	return fmt.Sprintf("%s", strings.Join(e.list, " "))
}

func (e MarkerExpr) Evaluate(binding Binding, locally bool) (interface{}, EvaluationInfo, bool) {
	info := DefaultInfo()
	for _, m := range e.list {
		switch m {
		case TEMPLATE:
			info.Issue = yaml.NewIssue("&template only usable as marker for templates")
			return nil, info, false
		case TEMPORARY:
			info.Temporary = true
		}
	}
	if e.expr != nil {
		result, infoe, ok := e.expr.Evaluate(binding, locally)
		infoe = infoe.Join(info)
		return result, infoe, ok
	}
	return nil, info, true
}

func (e MarkerExpr) setExpression(expr Expression) MarkerExpr {
	e.expr = expr
	return e
}

func (e MarkerExpr) Has(t string) bool {
	for _, v := range e.list {
		if v == t {
			return true
		}
	}
	return false
}

func (e MarkerExpr) add(m string) MarkerExpr {
	e.list = append(e.list, m)
	return e
}

func newMarkerExpr(m string) MarkerExpr {
	return MarkerExpr{list: []string{m}}
}