package table

import (
	"fmt"
	"reflect"
	"strings"
)

const reset = "\x1b[0m"
const innerSep = " "
const innerMargin = len(innerSep)

type Column struct {
	Type      reflect.Type
	Getter    Getter
	Name      string
	Title     string
	Alignment Alignment
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
	columnWidth map[string]int
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
		columnWidth: map[string]int{},
	}
}

func (t *Table) UpdateWidth(widthByColumn map[string]int) {
	for colName, width := range widthByColumn {
		if width > t.columnWidth[colName] {
			t.columnWidth[colName] = width
		}
	}
}

func (t *Table) Width(colName string) int {
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
	if len(value) > width {
		return strings.Repeat(" ", width)
	}
	return AlignmentCenter(value, width)
}

func (t *Table) FormatHeader(color string) string {
	str := color
	for _, col := range t.Columns {
		str += t.padColumnHeader(col) + " "
	}
	str += reset + "\n"
	return str
}

func (t *Table) TableWidth(margin int) int {
	width := (t.ColumnCount() - 1) * margin
	for _, col := range t.Columns {
		width += t.Width(col.Name)
	}
	return width
}
