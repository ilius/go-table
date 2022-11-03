package table

type FormattedItemList interface {
	Len() int
	Get(int) []string
}

type Getter interface {
	Value(item any) (any, error)
	ValueString(colName string, item any) (string, error)
	Format(item any, value any) (string, error)
}
