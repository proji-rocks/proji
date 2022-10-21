package text

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type (
	// TablePrinter is a printer for tables.
	TablePrinter interface {
		SetRowPrefix(prefix string)
		SetRowSuffix(suffix string)
		SetCellPadding(padding int)
		SetHeaderUnderlineSymbol(separator string)
		SetSink(w io.Writer)
		AddHeaderColumns(columns ...any)
		AddRow(values ...any)
		Render() error
		Reset()
	}

	tablePrinter struct {
		w                     io.Writer
		columns               []any
		columnsWidths         []int
		rows                  [][]any
		rowPrefix             string
		rowSuffix             string
		cellPadding           int
		headerUnderlineSymbol string
		mu                    sync.Mutex
	}
)

const (
	defaultRowPrefix             = " " // Prefix for each row.
	defaultRowSuffix             = ""  // Suffix for each row.
	defaultCellPadding           = 2   // Additional padding for each cell.
	defaultHeaderUnderlineSymbol = "-" // Symbol used to underline the header row.
)

// NewTablePrinter returns a new TablePrinter. The sink is set to os.Stdout by default. The horizontal separator is set
// to "".
func NewTablePrinter() TablePrinter {
	return &tablePrinter{
		w:                     os.Stdout,
		columns:               make([]any, 0),
		columnsWidths:         make([]int, 0),
		rows:                  make([][]any, 0),
		rowPrefix:             defaultRowPrefix,
		rowSuffix:             defaultRowSuffix,
		cellPadding:           defaultCellPadding,
		headerUnderlineSymbol: defaultHeaderUnderlineSymbol,
		mu:                    sync.Mutex{},
	}
}

// SetRowPrefix sets the prefix for each row. The default is " ".
func (p *tablePrinter) SetRowPrefix(prefix string) {
	p.rowPrefix = prefix
}

// SetRowSuffix sets the suffix for each row. The default is "".
func (p *tablePrinter) SetRowSuffix(suffix string) {
	p.rowSuffix = suffix
}

// SetCellPadding sets the padding for each cell. The default is 2.
func (p *tablePrinter) SetCellPadding(padding int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// We need this to recalculate the columns widths.
	previousPadding := p.cellPadding

	// Set the new padding.
	if padding < 0 {
		padding = 0
	}

	p.cellPadding = padding

	// Recalculate the columns widths.
	for idx, width := range p.columnsWidths {
		p.columnsWidths[idx] = width - previousPadding + padding
	}
}

// Reset clears the table.
func (p *tablePrinter) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.columns = make([]any, 0)
	p.columnsWidths = make([]int, 0)
	p.rows = make([][]any, 0)
}

// SetHeaderUnderlineSymbol sets the symbol used to underline the header row. The default is "". Spaces and new lines
// will be trimmed. If the symbol is empty, an empty bar will be rendered. The function expects only a single symbol.
func (p *tablePrinter) SetHeaderUnderlineSymbol(symbol string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Remove unwanted characters.
	symbol = strings.Trim(symbol, "\n")

	// Symbol should be a single character.
	if len(symbol) > 1 {
		symbol = symbol[:1]
	}

	// Set the symbol.
	p.headerUnderlineSymbol = symbol
}

// SetSink sets the sink for the table printer. The default is os.Stdout. It will not set a nil sink.
func (p *tablePrinter) SetSink(w io.Writer) {
	if w != nil {
		p.w = w
	}
}

func prepareEntry(entry any) string {
	// If entry is nil, set it to an empty string.
	if entry == nil {
		return ""
	}

	// Check if entry is a pointer.
	if reflect.TypeOf(entry).Kind() == reflect.Ptr {
		value := reflect.ValueOf(entry)
		if value.IsNil() || value.IsZero() {
			return "" // If entry is nil or zero value, set it to an empty string.
		}

		// If entry is not nil, set it to the value of the pointer.
		entry = value.Elem().Interface()
	}

	// Get string representation of entry and trim of unwanted characters.
	_entry := fmt.Sprintf("%v", entry)
	_entry = strings.Trim(_entry, " \n")

	return _entry
}

// calcColumnsWidths calculates the maximum width of each column. The function expects the columns to be passed as a
// slice of any type.
func (p *tablePrinter) calcColumnsWidths(entries ...any) {
	for idx, entry := range entries {
		preppedEntry := prepareEntry(entry)

		// Update the entry.
		entries[idx] = preppedEntry

		// If idx is out of bounds, expand the columnsWidths slice.
		if idx >= len(p.columnsWidths) {
			p.columnsWidths = append(p.columnsWidths, 0)
		}

		// Calculate the entries length and only add cell padding if the entry is not the last one and the normalized
		// row suffix is not empty.
		entryLength := len(preppedEntry)
		if idx < len(entries)-1 {
			entryLength += p.cellPadding
		}

		// Check if a new longest entry is found for each column. If so, update the max length.
		if p.columnsWidths[idx] < entryLength {
			p.columnsWidths[idx] = entryLength
		}
	}
}

// AddHeaderColumns adds the header columns to the table. An empty list of columns will be skipped.
func (p *tablePrinter) AddHeaderColumns(columns ...any) {
	if len(columns) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.columnsWidths == nil || len(p.columnsWidths) == 0 {
		p.columnsWidths = make([]int, len(columns))
	}

	p.columns = columns
	p.calcColumnsWidths(columns...)
}

// AddRow adds a row to the table. An empty list of values will be skipped.
func (p *tablePrinter) AddRow(values ...any) {
	if len(values) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.columnsWidths == nil || len(p.columnsWidths) == 0 {
		p.columnsWidths = make([]int, len(values))
	}

	// Make sure that the number of columns is the same for all rows.
	numMissingValues := len(p.columns) - len(values)
	if numMissingValues > 0 {
		values = append(values, make([]any, numMissingValues)...)
	} else if numMissingValues < 0 {
		p.columns = append(p.columns, make([]any, numMissingValues*-1)...)
	}

	// Check for new longest entry in each column.
	p.calcColumnsWidths(values...)

	// Add the rows
	p.rows = append(p.rows, values)
}

// calcTableWidth calculates the total width of the table.
func (p *tablePrinter) calcTableWidth() (totalWidth int) {
	// Add width of each column.
	for _, width := range p.columnsWidths {
		totalWidth += width
	}

	return totalWidth
}

// Render writes the table to the sink. The function expects the table to be populated with a header and rows.
func (p *tablePrinter) Render() error {
	if len(p.columns) == 0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Build row format string with respect to the longest entry in each column.
	rowFormat := p.rowPrefix
	for _, width := range p.columnsWidths {
		widthInt := strconv.Itoa(width)
		rowFormat += "%-" + widthInt + "v"
	}
	rowFormat += p.rowSuffix + "\n"

	// Begin table with header.
	var table string
	table += fmt.Sprintf(rowFormat, p.columns...)

	// Add separator row.
	totalWidth := p.calcTableWidth()
	table += p.rowPrefix + strings.Repeat(p.headerUnderlineSymbol, totalWidth) + p.rowSuffix + "\n"

	// Add row by row.
	for _, row := range p.rows {
		table += fmt.Sprintf(rowFormat, row...)
	}

	// Print the table.
	_, err := fmt.Fprint(p.w, table)

	return err
}
