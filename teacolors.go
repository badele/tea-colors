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
	"fmt"
	"strconv"
	"strings"

	"unicode/utf8"

	"github.com/muesli/termenv"
)

var base16 = []string{
	"black",
	"maroon",
	"green",
	"olive",
	"navy",
	"purple",
	"teal",
	"silver",
	"gray",
	"red",
	"lime",
	"yellow",
	"blue",
	"fuchsia",
	"aqua",
	"white",
}

var colwidth, allwidth int

type PadDirection int

const (
	LEFT PadDirection = iota
	CENTER
	RIGHT
)

// Pad the string content
func StrPad(direction PadDirection, text string, width int) string {
	output := ""
	textwith := utf8.RuneCountInString(text)
	nbspaces := width - textwith

	if nbspaces <= 0 {
		return text
	}

	switch direction {
	case RIGHT:
		output = strings.Repeat(" ", nbspaces) + text
	case LEFT:
		output = text + strings.Repeat(" ", nbspaces)
	case CENTER:
		leftspaces := int(nbspaces / 2)
		rightspace := nbspaces - leftspaces

		output = strings.Repeat(" ", leftspaces) + text + strings.Repeat(" ", rightspace)
	}

	return output
}

// Show colors panel
func ShowANSI16ColorsPanel(profile termenv.Profile, colors []string) {
	// Show ANSI Colors
	nbcolors := 16
	sep := ""

	// ###################################################################
	// # Band colors
	// ###################################################################
	bandsize := ((colwidth + 1) * nbcolors) / (nbcolors / 2)
	fmt.Printf(StrPad(RIGHT, "", colwidth))
	for idx := 0; idx < nbcolors/2; idx += 1 {
		fmt.Printf("%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(colors[idx])))
	}
	fmt.Println()

	fmt.Printf(StrPad(RIGHT, "", colwidth))
	for idx := nbcolors / 2; idx < nbcolors; idx += 1 {
		fmt.Printf("%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(colors[idx])))
	}
	fmt.Println()

	// ###################################################################
	// # Colors name
	// ###################################################################
	fmt.Printf(StrPad(RIGHT, "ANSI", colwidth))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Printf("%s %s", sep, StrPad(CENTER, strings.ToUpper(fmt.Sprintf("%02x", idx)), colwidth))
	}
	fmt.Println()

	fmt.Printf(StrPad(RIGHT, "Color", colwidth))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Printf("%s %s", sep, termenv.String(StrPad(CENTER, base16[idx], colwidth)).Foreground(profile.Color(colors[idx])))
	}
	fmt.Println()

	// ###################################################################
	// # Colors block
	// ###################################################################
	for row := 0; row < nbcolors; row += 1 {
		fmt.Printf(StrPad(RIGHT, base16[row], colwidth))
		for col := 0; col < nbcolors; col += 1 {
			fmt.Printf("%s %s", sep, termenv.String(StrPad(CENTER, "•••", colwidth)).Foreground(profile.Color(colors[row])).Background(profile.Color(colors[col])))
		}
		fmt.Println()
	}

	// Print colorname (in same color)
	fmt.Printf(StrPad(RIGHT, "Color", colwidth))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Printf("%s %s", sep, StrPad(CENTER, base16[idx], colwidth))
	}
	fmt.Println()
	fmt.Println()
}

func ShowGrayColorsPanel(profile termenv.Profile) {
	nbcolors := 255 - 232
	bandsize := ((colwidth + 1) * nbcolors) / 16

	allspaces := (allwidth - (bandsize * ((nbcolors / 2) + 1)))

	fmt.Printf(StrPad(RIGHT, "", allspaces))
	for idx := 232; idx <= 232+nbcolors/2; idx += 1 {
		fmt.Printf("%s", termenv.String(StrPad(CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("15")),
		)
	}
	fmt.Println()

	fmt.Printf(StrPad(RIGHT, "", allspaces))
	for idx := 232 + nbcolors/2 + 1; idx <= 255; idx += 1 {
		fmt.Printf("%s", termenv.String(StrPad(CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("0")),
		)
	}
	fmt.Println()
}

func main() {
	// Init termenv
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
	if err != nil {
		panic(err)
	}
	defer restoreConsole()

	// Get current colors profile
	colors := []string{}
	for i := 0; i < 16; i += 1 {
		colors = append(colors, fmt.Sprintf("%d", i))
	}

	colwidth = 0
	for _, colname := range base16 {
		size := len(colname)
		if size > colwidth {
			colwidth = size
		}
	}
	allwidth = ((colwidth + 1) * 17) - 2

	ShowANSI16ColorsPanel(termenv.ANSI, colors)
	ShowGrayColorsPanel(termenv.ANSI256)
}
