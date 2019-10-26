package item_test

import (
	"testing"

	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/stretchr/testify/assert"
)

func TestNewStatus(t *testing.T) {
	statusExp := &item.Status{
		ID:        99,
		Title:     "active",
		IsDefault: false,
		Comment:   "This project is active.",
	}

	statusAct := item.NewStatus(99, "active", "This project is active.", false)
	assert.Equal(t, statusExp, statusAct)
}
