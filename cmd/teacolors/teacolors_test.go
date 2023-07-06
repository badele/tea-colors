// Tea Colors
// Copyright (C) 2023  Bruno Adelé

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"strings"
	"testing"

	"github.com/badele/tea-colors/internal/pkg/ansi"
	"github.com/badele/tea-colors/internal/pkg/tools"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestStrPad(t *testing.T) {
	left := tools.StrPad(tools.LEFT, "LEFT", 10)
	assert.Equal(t, "LEFT      ", left)

	right := tools.StrPad(tools.RIGHT, "RIGHT", 10)
	assert.Equal(t, "     RIGHT", right)

	center := tools.StrPad(tools.CENTER, "CENTER", 10)
	assert.Equal(t, "  CENTER  ", center)

	biggest := tools.StrPad(tools.CENTER, "TOO BIGGEST", 10)
	assert.Equal(t, "TOO BIGGEST", biggest)

}

func TestGetGlobalConf(t *testing.T) {
	// onelist, twolist order
	assert.Equal(t, ansi.GetMaxColorNameWidth(), 7) // magenta
	assert.Equal(t, ansi.GetMaxStyleNameWidth(), 9) // underline

	// twolist, onelist order
	assert.Equal(t, ansi.GetStartPosX(), 9)

	AllWidth := (ansi.GetMaxColorNameWidth() + 1) * 16
	assert.Equal(t, ansi.GetAllWidth(), AllWidth)
}

func TestShowANSIColorBar(t *testing.T) {

	output := ansi.GetANSIColorBar()
	lines := strings.Split(output, "\n")

	// Test two row band colors (3 newlines)
	assert.Equal(t, 2+1, len(lines))
}

func TestShowANSI16ColorsPanel(t *testing.T) {
	output := ansi.GetANSI16ColorsPanel()
	lines := strings.Split(output, "\n")

	// test colors
	assert.Contains(t, lines[4], "\x1b[30;41m")  // Black(background Red) line
	assert.Contains(t, lines[11], "\x1b[37;42")  // Silver/(background Green) line
	assert.Contains(t, lines[19], "\x1b[97;43m") // White/(background Yellow) line
}

func TestShowTextStylePanel(t *testing.T) {
	output := ansi.GetTextStylePanel()
	lines := strings.Split(output, "\n")

	// test styles
	assert.Contains(t, lines[1], "\x1b[31;1m  ABC") // Bold Red ABC
	assert.Contains(t, lines[7], "\x1b[31;5m  ABC") // Blink Red ABC
}

func TestShowGrayColorsPanel(t *testing.T) {
	output := ansi.GetGrayColorsPanel(termenv.ANSI)
	lines := strings.Split(output, "\n")

	// Test two row band colors (3 newlines)
	assert.Equal(t, 2+1, len(lines))
}

func TestOutputBar(t *testing.T) {
	output := OutputBar()
	lines := strings.Split(output, "\n")

	// All lines
	assert.Equal(t, 32+1, len(lines))
}
