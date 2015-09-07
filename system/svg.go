// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Options are a collection of settings for rendering systems as SVG.
type Options struct {
	Depth     int     // Depth of recursion.
	Thickness float64 // Stroke thickness.
	Color     string  // Stroke color.
	Precision int     // Float precision (number of decimal places).
}

// SVG renders the system as a curve, returning a string of SVG data.
func (s *System) SVG(opts *Options) string {
	segments := s.render(s.min + opts.Depth)
	view := calcViewBox(segments, opts)
	thickness := opts.Thickness * view.w / stepFactor

	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		"<svg xmlns='http://www.w3.org/2000/svg' viewBox='%s %s %s %s'>"+
			"<defs><style>polyline { fill: none; stroke-linecap: square; "+
			"stroke-width: %s; stroke: %s; }</style></defs>",
		precise(view.x, opts.Precision),
		precise(view.y, opts.Precision),
		precise(view.w, opts.Precision),
		precise(view.h, opts.Precision),
		precise(thickness, opts.Precision),
		opts.Color,
	)

	for _, points := range segments {
		buf.WriteString("<polyline points='")
		for _, pt := range points {
			x := precise(pt.x, opts.Precision)
			y := precise(pt.y, opts.Precision)
			fmt.Fprintf(&buf, "%s,%s ", x, y)
		}
		buf.WriteString("'/>")
	}
	buf.WriteString("</svg>")

	return buf.String()
}

// A viewBox describes the boundaries of an SVG image.
type viewBox struct {
	x, y, w, h float64
}

// paddingFactor is used to add additional padding to the viewBox.
const paddingFactor = 1.5

// calcViewBox returns a viewBox large enough to contain all the points.
func calcViewBox(segments [][]vector, opts *Options) viewBox {
	var xMin, xMax, yMin, yMax float64
	for _, points := range segments {
		for _, pt := range points {
			if pt.x < xMin {
				xMin = pt.x
			} else if pt.x > xMax {
				xMax = pt.x
			}
			if pt.y < yMin {
				yMin = pt.y
			} else if pt.y > yMax {
				yMax = pt.y
			}
		}
	}

	edge := paddingFactor * opts.Thickness / 2
	xMin -= edge
	yMin -= edge
	xMax += edge
	yMax += edge

	width := xMax - xMin
	height := yMax - yMin

	return viewBox{xMin, yMin, width, height}
}

// precise returns a float formatted as a string with the given number of
// decimal places, omitting trailing zeros.
func precise(f float64, precision int) string {
	str := strconv.FormatFloat(f, 'f', precision, 64)
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")
	return str
}
