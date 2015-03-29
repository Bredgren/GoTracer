package gotracer

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParser(t *testing.T) {
	scene := NewSceneFromFile("scene/example.json")
	t.Log("Parsed scene:", spew.Sdump(scene))
}
