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

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestStrPad(t *testing.T) {
	left := StrPad(LEFT, "LEFT", 10)
	assert.Equal(t, "LEFT      ", left)

	right := StrPad(RIGHT, "RIGHT", 10)
	assert.Equal(t, "     RIGHT", right)

	center := StrPad(CENTER, "CENTER", 10)
	assert.Equal(t, "  CENTER  ", center)

	biggest := StrPad(CENTER, "TOO BIGGEST", 10)
	assert.Equal(t, "TOO BIGGEST", biggest)

}

func TestGetGlobalConf(t *testing.T) {
	onelist := []string{
		"one",
		"five",
		"two",
	}

	twolist := []string{
		"one",
		"eleven",
		"two",
	}

	// onelist, twolist order
	globalconfig = GetGlobalConf(onelist, twolist)
	assert.Equal(t, globalconfig.maxcolnamewidth, 4)
	assert.Equal(t, globalconfig.startposx, 6)

	// twolist, onelist order
	globalconfig = GetGlobalConf(twolist, onelist)
	assert.Equal(t, globalconfig.maxcolnamewidth, 6)
	assert.Equal(t, globalconfig.startposx, 6)

	allwidth := ((globalconfig.maxcolnamewidth + 1) * 16) - 2
	assert.Equal(t, globalconfig.allwidth, allwidth)
}

func TestShowANSIColorBar(t *testing.T) {

	globalconfig = GetGlobalConf(COLORNAMES, TEXTSTYLES)
	output := ShowANSIColorBar(termenv.ANSI)
	lines := strings.Split(output, "\n")

	// Test two row band colors (3 newlines)
	assert.Equal(t, 2+1, len(lines))
}

func TestShowANSI16ColorsPanel(t *testing.T) {
	globalconfig = GetGlobalConf(COLORNAMES, TEXTSTYLES)
	output := ShowANSI16ColorsPanel(termenv.ANSI)
	lines := strings.Split(output, "\n")

	// test colors
	assert.Contains(t, lines[2], "\x1b[30;41m")  // Black(background Red) line
	assert.Contains(t, lines[9], "\x1b[37;42")   // Silver/(background Green) line
	assert.Contains(t, lines[17], "\x1b[97;43m") // White/(background Yellow) line
}

func TestShowTextStylePanel(t *testing.T) {
	globalconfig = GetGlobalConf(COLORNAMES, TEXTSTYLES)
	output := ShowTextStylePanel(termenv.ANSI)
	lines := strings.Split(output, "\n")

	// test styles
	assert.Contains(t, lines[1], "\x1b[31;1m  ABC") // Bold Red ABC
	assert.Contains(t, lines[7], "\x1b[31;5m  ABC") // Blink Red ABC
}

func TestShowGrayColorsPanel(t *testing.T) {
	globalconfig = GetGlobalConf(COLORNAMES, TEXTSTYLES)
	output := ShowGrayColorsPanel(termenv.ANSI)
	lines := strings.Split(output, "\n")

	// Test two row band colors (3 newlines)
	assert.Equal(t, 2+1, len(lines))
}

func TestOutputBar(t *testing.T) {
	output := OutputBar()
	lines := strings.Split(output, "\n")

	// All lines
	assert.Equal(t, 30+1, len(lines))
}
