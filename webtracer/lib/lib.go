package lib

import (
	"regexp"
	"time"
)

// RenderItem groups the files associated with a single render.
type RenderItem struct {
	Render, Thumb, Scene string
	Date                 time.Time
}

// RenderFileRe is the file name format used.
var RenderFileRe = regexp.MustCompile(`render-(\d+)(-thumb)?.(jpg|jpeg|png|json)`)
