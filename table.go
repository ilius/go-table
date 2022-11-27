package table

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	innerSep    = " "
	innerMargin = uint16(len(innerSep))
	MAX_WIDTH   = 65535 // 2^16 - 1
)

type Column struct {
	Type      reflect.Type
	Getter    Getter
	Alignment Alignment
	Name      string
	Title     string
}

type TableSpec struct {
	ColumnByName map[string]*Column
	TimeFormat   string
	Columns      []*Column
}

func (t *TableSpec) HasColumn(colName string) bool {
	return t.ColumnByName[colName] != nil
}

func (t *TableSpec) AddColumn(col *Column) {
	t.Columns = append(t.Columns, col)
	t.ColumnByName[col.Name] = col
}

func (t *TableSpec) ColumnCount() int {
	return len(t.Columns)
}

type Table struct {
	*TableSpec
	columnWidth map[string]uint16
	// Data        []any
}

func NewTableSpec() *TableSpec {
	return &TableSpec{
		Columns:      []*Column{},
		ColumnByName: map[string]*Column{},
	}
}

func NewTable(spec *TableSpec) *Table {
	if spec == nil {
		spec = NewTableSpec()
	}
	return &Table{
		TableSpec:   spec,
		columnWidth: map[string]uint16{},
	}
}

func (t *Table) UpdateWidth(widthByColumn map[string]uint16) {
	for colName, width := range widthByColumn {
		if width > t.columnWidth[colName] {
			t.columnWidth[colName] = width
		}
	}
}

func (t *Table) Width(colName string) uint16 {
	return t.columnWidth[colName]
}

// FormatItemBasic formats item for non-tabular formats like json and csv
// using col.Getter.ValueString
func (t *Table) FormatItemBasic(item any, sep string) (string, error) {
	str := ""
	for index, col := range t.Columns {
		colStr, err := col.Getter.ValueString(col.Name, item)
		if err != nil {
			return "", err
		}
		if index > 0 {
			str += sep
		}
		str += colStr
	}
	return str, nil
}

func (t *Table) FormatItem(item any) ([]string, error) {
	cw := t.columnWidth
	formatted := make([]string, t.ColumnCount())
	for i, col := range t.Columns {
		value, err := col.Getter.Value(item)
		if err != nil {
			return nil, err
		}
		//if reflect.TypeOf(value) != col.Type {
		//	fmt.Fprintf(os.Stderr, "invalid type %T for column %v, must be %v\n", value, col.Name, col.Type)
		//}
		valueFormatted, err := col.Getter.Format(item, value)
		if err != nil {
			return nil, err
		}
		formatted[i] = valueFormatted
		// even if col.Alignment == nil, we may need it for MergeRows* funcs
		width := visualWidth(valueFormatted)
		if width > cw[col.Name] {
			cw[col.Name] = width
		}
	}
	return formatted, nil
}

func (t *Table) AlignFormattedItem(formatted []string) ([]string, error) {
	if len(formatted) != t.ColumnCount() {
		return nil, fmt.Errorf("bad number of columns: %d, must be %d", len(formatted), t.ColumnCount())
	}
	for i, col := range t.Columns {
		if col.Alignment == nil {
			continue
		}
		formatted[i] = col.Alignment(
			formatted[i],
			t.Width(col.Name),
		)
	}
	return formatted, nil
}

func (t *Table) padColumnHeader(col *Column) string {
	width := t.Width(col.Name)
	value := col.Title
	if width == 0 {
		return value
	}
	if len(value) > int(width) {
		return strings.Repeat(" ", int(width))
	}
	return AlignmentCenter(value, width)
}

func (t *Table) FormatHeader(sep string) string {
	str := ""
	for _, col := range t.Columns {
		str += t.padColumnHeader(col) + sep
	}
	str += "\n"
	return str
}

func (t *Table) TableWidth(margin uint16) uint16 {
	width := (uint16(t.ColumnCount()) - 1) * margin
	for _, col := range t.Columns {
		width += t.Width(col.Name)
	}
	return width
}
