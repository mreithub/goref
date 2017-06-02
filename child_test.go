package goref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChild(t *testing.T) {
	g := NewGoRef()

	assert.Equal(t, g, g.GetChild()) // empty parameters -> return itself

	dbInstance := g.GetChild("backend", "db")
	assert.Equal(t, []string{"backend", "db"}, dbInstance.GetPath())

	// check ._children
	assert.Contains(t, g._children, "backend")
	assert.NotContains(t, g._children, "db")

	backendInstance := g._children["backend"]
	assert.Equal(t, backendInstance, g.GetChild("backend")) // is supposed to return the existing instance
	assert.Contains(t, backendInstance._children, "db")
	assert.Equal(t, backendInstance._children["db"], dbInstance)

	// check .parent
	assert.Equal(t, backendInstance, dbInstance.parent)
	assert.Equal(t, g, backendInstance.GetParent())
	assert.Nil(t, g.parent)

	// check .name
	assert.Equal(t, "", g.name)
	assert.Equal(t, "backend", backendInstance.name)
	assert.Equal(t, "db", dbInstance.name)
}
