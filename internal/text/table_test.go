package text

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

func TestNewTablePrinter(t *testing.T) {
	t.Parallel()

	table := NewTablePrinter()
	if table == nil {
		t.Fatal("NewTablePrinter() returned nil")
	}

	// Type assertion so that we can access private fields.
	_table, ok := table.(*tablePrinter)
	if !ok {
		t.Fatal("NewTablePrinter() returned wrong type")
	}

	// Test sink
	if _table.w != os.Stdout {
		t.Fatal("NewTablePrinter() did not set correct writer")
	}

	// Should never set nil sink
	_table.SetSink(nil)
	if _table.w == nil {
		t.Fatal("SetSink() set nil sink")
	}

	// Forcefully setting sink to nil, so that we can confirm SetSink()'s expected behaviour.
	_table.w = nil
	_table.SetSink(os.Stdout)
	if _table.w != os.Stdout {
		t.Fatalf("SetSink() expected to set %T as sink, but got %T", os.Stdout, _table.w)
	}

	// Test setting columns
	if _table.columns == nil {
		t.Fatal("NewTablePrinter() did not initialize columns")
	}
	if _table.columnsWidths == nil {
		t.Fatal("NewTablePrinter() did not initialize columnsWidths")
	}

	_table.AddHeaderColumns("a", "b", "c")
	if len(_table.columns) != 3 {
		t.Fatalf("AddHeaderColumns() expected to add 3 columns, but got %d", len(_table.columns))
	}
	if _table.columns[0] != "a" || _table.columns[1] != "b" || _table.columns[2] != "c" {
		t.Fatalf("AddHeaderColumns() expected to add columns in order, but got %v", _table.columns)
	}

	// Test setting rows
	if _table.rows == nil {
		t.Fatal("NewTablePrinter() did not initialize rows")
	}

	_table.AddRow("1", "2", "3")
	if len(_table.rows) != 1 {
		t.Fatalf("AddRow() expected to add 1 row, but got %d", len(_table.rows))
	} else {
		if len(_table.rows[0]) != 3 {
			t.Fatalf("AddRow() expected to add 3 columns, but got %d", len(_table.rows[0]))
		}
	}
	if _table.rows[0][0] != "1" || _table.rows[0][1] != "2" || _table.rows[0][2] != "3" {
		t.Fatalf("AddRow() expected to add rows in order, but got %v", _table.rows)
	}

	// Test header underline symbol
	if _table.headerUnderlineSymbol != defaultHeaderUnderlineSymbol {
		t.Fatalf("headerUnderlineSymbol expected to be %s, but got %s", defaultHeaderUnderlineSymbol, _table.headerUnderlineSymbol)
	}
	_table.SetHeaderUnderlineSymbol("-")
	if _table.headerUnderlineSymbol != "-" {
		t.Fatalf("SetHeaderUnderlineSymbol() expected to set '-' as symbol, but got %v", _table.headerUnderlineSymbol)
	}

	// Try to set a symbol that is too long
	_table.SetHeaderUnderlineSymbol("^-^")
	if _table.headerUnderlineSymbol != "^" {
		t.Fatalf("SetHeaderUnderlineSymbol() expected to set '^' as symbol, but got %v", _table.headerUnderlineSymbol)
	}
}

func TestTablePrinter__calculations(t *testing.T) {
	t.Parallel()

	table := NewTablePrinter()
	if table == nil {
		t.Fatal("NewTablePrinter() returned nil")
	}

	// Type assertion so that we can access private fields.
	_table, ok := table.(*tablePrinter)
	if !ok {
		t.Fatal("NewTablePrinter() returned wrong type")
	}

	// Test setting columns
	columns := []any{"a", "bb", "ccc", nil}
	wantWidths := []int{3, 4, 5, 0} // Added length of cell padding to all but the last column.

	_table.AddHeaderColumns(columns...)
	_table.calcColumnsWidths(columns...)

	if len(_table.columns) != 4 {
		t.Fatalf("AddHeaderColumns() expected to add %d columns, but got %d", len(columns), len(_table.columns))
	}

	for i, column := range _table.columns {
		if wantWidths[i] != _table.columnsWidths[i] {
			t.Fatalf("calcColumnsWidths() expected to set %d as width for %q, but got %d",
				wantWidths[i], column, _table.columnsWidths[i])
		}
	}

	// Writing this very explicitly here to avoid confusion. The total width gets automatically calculated. Per column,
	// the width results out of the length of its value + the column padding. To get the total width, we add up all the
	// column widths except for the last column. See the documentation for TablePrinter.calcTableWidth() for more details.
	wantTableWidth := (1 + defaultCellPadding) + (2 + defaultCellPadding) + (3 + defaultCellPadding) + 0
	gotTableWidth := _table.calcTableWidth()

	if gotTableWidth != wantTableWidth {
		t.Fatalf("calcTableWidth() expected to return %d, but got %d", wantTableWidth, gotTableWidth)
	}

	// Check resetting the table
	_table.Reset()
	if len(_table.columns) != 0 || len(_table.columnsWidths) != 0 || len(_table.rows) != 0 {
		t.Fatalf("Reset() did not reset the table correctly")
	}

	// Check more uncommon types
	upstream := "https://github.com/nikoksr/proji"
	createdAt := time.Now()
	pkg := domain.Package{
		Label:       "t",
		Name:        "Test",
		Description: nil,
		UpstreamURL: &upstream,
		CreatedAt:   createdAt,
	}

	_table.AddRow(pkg.Label, pkg.Name, pkg.Description, pkg.UpstreamURL, pkg.CreatedAt)
	if len(_table.rows) != 1 {
		t.Fatalf("AddRow() expected to add 1 row, but got %d", len(_table.rows))
	}
	if len(_table.rows[0]) != 5 {
		t.Fatalf("AddRow() expected to add 5 columns, but got %d", len(_table.rows[0]))
	}
	if _table.rows[0][0] != pkg.Label {
		t.Errorf("Wrong value for label: %v", _table.rows[0][0])
	}
	if _table.rows[0][1] != pkg.Name {
		t.Errorf("Wrong value for name: %v", _table.rows[0][1])
	}
	if _table.rows[0][2] != "" {
		t.Errorf("Wrong value for description: %v", _table.rows[0][2])
	}
	if _table.rows[0][3] != *pkg.UpstreamURL {
		t.Errorf("Wrong value for upstream: %v", _table.rows[0][2])
	}
	if _table.rows[0][4] != pkg.CreatedAt.String() {
		t.Errorf("Wrong value for createdAt: %v", _table.rows[0][3])
	}

	// Edge cases
	_table.Reset()

	// No columns given
	_table.AddHeaderColumns()
	if len(_table.columns) != 0 {
		t.Fatalf("AddHeaderColumns() expected to add 0 columns, but got %d", len(_table.columns))
	}

	// No rows given
	_table.AddRow()
	if len(_table.rows) != 0 {
		t.Fatalf("AddRow() expected to add 0 rows, but got %d", len(_table.rows))
	}

	// More columns than rows; this should add one empty cell per row
	_table.Reset()
	_table.AddHeaderColumns("a", "b", "c")
	_table.AddRow("1", "2")

	if len(_table.rows) != 1 {
		t.Fatalf("AddRow() expected to add 1 row, but got %d", len(_table.rows))
	}
	if len(_table.rows[0]) != 3 {
		t.Fatalf("AddRow() expected to add 3 columns, but got %d", len(_table.rows[0]))
	}
	if _table.rows[0][2] != "" {
		t.Errorf("Wrong value for column c: \"%v\" expected: \"\"", _table.rows[0][2])
	}

	// More rows than columns; this should add one header column without a title
	_table.Reset()
	_table.AddHeaderColumns("a")
	_table.AddRow("1", "2", "3")

	if len(_table.rows) != 1 {
		t.Fatalf("AddRow() expected to add 1 row, but got %d", len(_table.rows))
	}
	if len(_table.rows[0]) != 3 {
		t.Fatalf("AddRow() expected to create 3 cell row, but got %d", len(_table.rows[0]))
	}
	if len(_table.columns) != 3 {
		t.Fatalf("AddRow() expected to expand columns to 3, but got %d", len(_table.columns))
	}
	if _table.columns[1] != nil || _table.columns[2] != nil {
		t.Fatalf("AddRow() expected to create empty column headers, but created %v", _table.columns)
	}
}

func TestTablePrinter__Render(t *testing.T) {
	t.Parallel()

	// Create table
	_table := NewTablePrinter()

	// Set bytes.Buffer as sink, so we can check the output
	got := new(bytes.Buffer)
	_table.SetSink(got)

	// Style settings
	_table.SetHeaderUnderlineSymbol("=")
	_table.SetRowPrefix("")
	_table.SetRowSuffix("")
	_table.SetCellPadding(1)

	// Fill table
	_table.AddHeaderColumns("a", "b", "c")
	_table.AddRow("1", "2", "3")
	_table.AddRow("4", "5", "6")
	_table.AddRow("7", "8", "9")

	// Render the first table
	err := _table.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	// Verify the first table is rendered correctly
	want := `a b c
=====
1 2 3
4 5 6
7 8 9
`

	if !cmp.Equal(want, got.String()) {
		t.Fatalf("Render() rendered unequal tables:\n\nwant:\n%s\n\ngot:%s\n", want, got.String())
	}

	// Reset the table data and set new style settings
	got.Reset()
	_table.SetHeaderUnderlineSymbol("-")
	_table.SetRowPrefix("| ")
	_table.SetRowSuffix(" |")
	_table.SetCellPadding(3)

	// Render the second table
	err = _table.Render()
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	// Verify the second table is rendered correctly
	want = `| a   b   c   |
| ----------- |
| 1   2   3   |
| 4   5   6   |
| 7   8   9   |
`

	if !cmp.Equal(want, got.String()) {
		t.Fatalf("Render() rendered unequal tables:\nwant:\n%s\ngot:\n%s\n", want, got.String())
	}
}
