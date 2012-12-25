package mysql

import "encoding/xml"

type Database struct {
	XMLName xml.Name `xml:"database"`
	Name string `xml:"name,attr"`
	Charset string `xml:"charset,attr"`
	Collation string `xml:"collation,attr"`
	Tables []Table `xml:"tables>table,empty"`
}

type Table struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Engine string `xml:"engine,attr"`
	Collation string `xml:"collation,attr"`
	Columns []Column `xml:"columns>column"`
	Indexes []Index `xml:"indexes>index,omitempty"`
	ForeignKeys []ForeignKey `xml:"foreign_keys>foreign_key,omitempty"`
	Rows []Row `xml:"rows>row,omitempty"`
	TruncateTable bool `xml:"truncate,attr"`
}

type Column struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Length int64 `xml:"length,attr"`
	FullType string `xml:"fulltype,attr"`
	Null bool `xml:"null,attr"`
	Default string `xml:"default,attr"`
	Charset string `xml:"charset,attr"`
	Collation string `xml:"collation,attr"`
	Extra string `xml:"extra,attr"`
}

type Index struct {
	Name string `xml:"name,attr"`
	Columns []string `xml:"columns>column"`
	Unique bool `xml:"unique,attr"`
	Collation string `xml:"collation,attr"`
	Null bool `xml:"null,attr"`
	Type string `xml:"type,attr"`
}

type ForeignKey struct {
	Name string `xml:"name,attr"`
	Columns []ColumnReference `xml:"columns>column"`
	Schema string `xml:"schema,attr"`
	Table string `xml:"table,attr"`
	OnUpdate string `xml:"on_update,attr"`
	OnDelete string `xml:"on_delete,attr"`
}

type ColumnReference struct {
	Referencer string `xml:"referencer,attr"`
	Referenced string `xml:"referenced,attr"`
}

type Row struct {
	Fields []Field `xml:"field"`
}

type Field struct {
	Name string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
