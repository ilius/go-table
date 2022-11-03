package table

import (
	"fmt"
	"io"
	"strings"
)

func (t *Table) MergeRowsHorizontal(
	out io.Writer,
	items FormattedItemList,
	maxWidth int,
	sep string,
	compact bool,
) {
	margin := len(sep)
	colN := t.ColumnCount()
	groupCount := (maxWidth + margin) / (t.TableWidth(margin) + margin)
	if groupCount < 1 {
		groupCount = 1
	}
	getWidth := func(colI int, groupI int) int {
		return t.columnWidth[t.Columns[colI].Name]
	}
	if compact {
		extra, cellWidth := t.compactCalcH(items, maxWidth, margin, groupCount)
		if extra > 0 {
			groupCount += extra
			getWidth = func(colI int, groupI int) int {
				return cellWidth[groupI*colN+colI]
			}
		}
	}
	itemN := items.Len()
	lineCount := (itemN-1)/groupCount + 1
	for lineI := 0; lineI < lineCount; lineI++ {
		line := make([]string, 0, groupCount)
		for groupI := 0; groupI < groupCount; groupI++ {
			itemIdx := lineI*groupCount + groupI
			if itemIdx >= itemN {
				line = append(line, "")
				break
			}
			item := items.Get(itemIdx)
			cell := make([]string, colN)
			for colI, col := range t.Columns {
				al := col.Alignment
				if al == nil {
					al = AlignmentLeft
				}
				cell[colI] = al(
					item[colI],
					getWidth(colI, groupI),
				)
			}
			line = append(line, strings.Join(cell, innerSep))
		}
		fmt.Fprintln(out, strings.Join(line, sep))
	}
}

func (t *Table) compactCalcH(
	items FormattedItemList,
	maxWidth int,
	margin int,
	groupCountInit int,
) (extra int, cellWidth []int) {
	colN := t.ColumnCount()
	for {
		groupCount := groupCountInit + extra + 1
		cellWidthNew := t.mergedRowsWidthH(items, groupCount)
		totalWidth := margin*(groupCount-1) + innerMargin*(colN-1)*groupCount
		for _, w := range cellWidthNew {
			totalWidth += w
		}
		if totalWidth > maxWidth {
			return
		}
		cellWidth = cellWidthNew
		extra += 1
	}
	return 0, nil
}

func (t *Table) mergedRowsWidthH(
	items FormattedItemList,
	groupCount int,
) []int {
	itemN := items.Len()
	colN := t.ColumnCount()
	width := make([]int, 0, groupCount*colN)
	for groupI := 0; groupI < groupCount; groupI++ {
		for colI := 0; colI < colN; colI++ {
			mw := 0
			for itemIdx := groupI; itemIdx < itemN; itemIdx += groupCount {
				w := visualWidth(items.Get(itemIdx)[colI])
				if w > mw {
					mw = w
				}
			}
			width = append(width, mw)
		}
	}
	return width
}
