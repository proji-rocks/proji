package item

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatus(t *testing.T) {
	statusExp := &Status{
		ID:        99,
		Title:     "active",
		IsDefault: false,
		Comment:   "This project is active.",
	}

	statusAct := NewStatus(99, "active", "This project is active.", false)
	assert.Equal(t, statusExp, statusAct)
}
