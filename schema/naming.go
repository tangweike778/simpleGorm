package schema

import (
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Replacer replacer interface like strings.Replacer
type Replacer interface {
	Replace(name string) string
}

var _ Namer = (*NamingStrategy)(nil)

// NamingStrategy tables, columns naming strategy
type NamingStrategy struct {
	TablePrefix         string
	SingularTable       bool
	NameReplacer        Replacer
	NoLowerCase         bool
	IdentifierMaxLength int
}

// TableName convert string to table name
func (ns NamingStrategy) TableName(str string) string {
	if ns.SingularTable {
		return ns.TablePrefix + ns.toDBName(str)
	}
	return ns.TablePrefix + inflection.Plural(ns.toDBName(str))
}

// SchemaName generate schema name from table name, don't guarantee it is the reverse value of TableName
func (ns NamingStrategy) SchemaName(table string) string {
	table = strings.TrimPrefix(table, ns.TablePrefix)

	if ns.SingularTable {
		return ns.toSchemaName(table)
	}
	return ns.toSchemaName(inflection.Singular(table))
}

// ColumnName convert string to column name
func (ns NamingStrategy) ColumnName(table, column string) string {
	return ns.toDBName(column)
}

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, cases.Title(language.Und).String(initialism))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func (ns NamingStrategy) toDBName(name string) string {
	if name == "" {
		return ""
	}

	if ns.NameReplacer != nil {
		tmpName := ns.NameReplacer.Replace(name)

		if tmpName == "" {
			return name
		}

		name = tmpName
	}

	if ns.NoLowerCase {
		return name
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool // upper case == true
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}
	ret := buf.String()
	return ret
}

func (ns NamingStrategy) toSchemaName(name string) string {
	result := strings.ReplaceAll(cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(name, "_", " ")), " ", "")
	for _, initialism := range commonInitialisms {
		result = regexp.MustCompile(cases.Title(language.Und, cases.NoLower).String(strings.ToLower(initialism))+"([A-Z]|$|_)").ReplaceAllString(result, initialism+"$1")
	}
	return result
}
