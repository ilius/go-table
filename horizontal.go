package table

import (
	"fmt"
	"io"
	"strings"
)

func (t *Table) MergeRowsHorizontal(
	out io.Writer,
	items FormattedItemList,
	maxWidthArg int,
	sep string,
	compact bool,
) {
	margin := uint16(len(sep))
	colN := t.ColumnCount()
	if maxWidthArg > MAX_WIDTH {
		maxWidthArg = MAX_WIDTH
	}
	maxWidth := uint16(maxWidthArg)
	groupCount := int((maxWidth + margin) / (t.TableWidth(margin) + margin))
	if groupCount < 1 {
		groupCount = 1
	}
	getWidth := func(colI int, _ int) uint16 {
		return t.columnWidth[t.Columns[colI].Name]
	}
	if compact {
		extra, cellWidth := t.compactCalcH(items, maxWidth, margin, groupCount)
		if extra > 0 {
			groupCount += extra
			getWidth = func(colI int, groupI int) uint16 {
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
		_, err := fmt.Fprintln(out, strings.Join(line, sep))
		if err != nil {
			panic(err)
		}
	}
}

func (t *Table) compactCalcH(
	items FormattedItemList,
	maxWidth uint16,
	margin uint16,
	groupCountInit int,
) (extra int, cellWidth []uint16) {
	colN := uint16(t.ColumnCount())
	for {
		groupCount := groupCountInit + extra + 1
		cellWidthNew := t.mergedRowsWidthH(items, groupCount)
		totalWidth := margin*(uint16(groupCount)-1) + innerMargin*(colN-1)*uint16(groupCount)
		for _, w := range cellWidthNew {
			totalWidth += w
		}
		if totalWidth > maxWidth {
			return
		}
		cellWidth = cellWidthNew
		extra += 1
	}
}

func (t *Table) mergedRowsWidthH(
	items FormattedItemList,
	groupCount int,
) []uint16 {
	itemN := items.Len()
	colN := t.ColumnCount()
	width := make([]uint16, 0, groupCount*colN)
	for groupI := 0; groupI < groupCount; groupI++ {
		for colI := 0; colI < colN; colI++ {
			mw := uint16(0)
			for itemIdx := int(groupI); itemIdx < itemN; itemIdx += int(groupCount) {
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
