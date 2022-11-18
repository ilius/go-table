package table

import (
	"fmt"
	"io"
	"strings"
)

func (t *Table) MergeRowsVertical(
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
	getWidth := func(colI int, groupI int) uint16 {
		return t.columnWidth[t.Columns[colI].Name]
	}
	if compact {
		extra, cellWidth := t.compactCalcV(items, maxWidth, margin, groupCount)
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
			itemIdx := groupI*lineCount + lineI
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

func (t *Table) compactCalcV(
	items FormattedItemList,
	maxWidth uint16,
	margin uint16,
	groupCountInit int,
) (extra int, cellWidth []uint16) {
	colN := uint16(t.ColumnCount())
	for {
		groupCount := groupCountInit + extra + 1
		cellWidthNew := t.mergedRowsWidthV(items, groupCount)
		if cellWidthNew == nil {
			return
		}
		fmt.Printf("extra=%d, cellWidthNew=%#v\n", extra, cellWidthNew)
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
	return 0, nil
}

func (t *Table) mergedRowsWidthV(
	items FormattedItemList,
	groupCount int,
) []uint16 {
	itemN := items.Len()
	colN := t.ColumnCount()
	lineCount := (itemN-1)/groupCount + 1 // correct
	if groupCount > (itemN-1)/lineCount+1 {
		// <=> groupCount - 1 > (itemN - 1) / lineCount
		// <=> (groupCount - 1) * lineCount > itemN-1
		// so the last run of loop won't do anything other than appending
		// some zero to cellWidth array
		// and this groupCount will not give us anything
		return nil
	}
	cellWidth := make([]uint16, 0, groupCount*colN)
	for groupI := 0; groupI < groupCount; groupI++ {
		endItemIdx := (groupI + 1) * lineCount
		if endItemIdx > itemN {
			endItemIdx = itemN
		}
		for colI := 0; colI < colN; colI++ {
			mcw := uint16(0) // max cell width
			for itemIdx := groupI * lineCount; itemIdx < endItemIdx; itemIdx++ {
				cw := visualWidth(items.Get(itemIdx)[colI])
				if cw > mcw {
					mcw = cw
				}
			}
			// cellWidth[groupI*colN+colI] = mcw
			cellWidth = append(cellWidth, mcw)
		}
	}
	return cellWidth
}
