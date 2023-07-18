package sql

import (
	"fmt"
	"strings"

	"git.autops.xyz/autops/base/logs"
)

var (
	eq  = "="
	neq = "<>"
	gt  = ">"
	gte = ">="
	st  = "<"
	ste = "<="
	not = "is not"
	lk  = "LIKE"
	in  = "in"
	// 这个表示参数错误
	no = "no"
	is = "is"

	ac = "ac" //asc
	dc = "dc" //desc
)

func predicate(p string) string {
	switch p {
	case "eq":
		return eq
	case "neq":
		return neq
	case "gt":
		return gt
	case "gte":
		return gte
	case "st":
		return st
	case "ste":
		return ste
	case "lk":
		return lk
	case "not":
		return not
	case "in":
		return in
	case "is":
		return is
	}
	return no
}

// BuildOrder ...
func BuildOrder(key string) string {
	field := ""
	p := ""
	query := strings.Split(key, "__")
	if len(query) == 1 {
		field, p = query[0], dc
	} else {
		field, p = query[0], query[1]
	}

	order := "DESC"
	if p == ac {
		order = "ASC"
	}
	return fmt.Sprintf("%s %s", field, order)
}

// BuildWhere ...
func BuildWhere(key string, value interface{}) (string, interface{}) {
	query := strings.Split(key, "__")
	if len(query) == 1 {
		return fmt.Sprintf("%s %s ?", key, eq), value
	} else if len(query) != 2 {
		logs.Errorf("where error: %v", key)
		return "", ""
	}
	field := query[0]
	p := query[1]
	predicate := predicate(p)
	if p == no {
		logs.Errorf("where error: %v", key)
		return "", ""
	}
	if p == "lk" {
		value = "%" + value.(string) + "%"
	}
	prefix := fmt.Sprintf("%s %s ?", field, predicate)
	if p == "in" {
		value = strings.Split(value.(string), "|")
		prefix = fmt.Sprintf("%s %s (?)", field, predicate)
	}

	return prefix, value
}
