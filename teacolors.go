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
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"unicode/utf8"

	"github.com/muesli/termenv"
)

var COLORNAMES = []string{
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

var TEXTSTYLES = []string{
	"Normal",
	"Bold",
	"Faint",
	"Italic",
	"Overline",
	"CrossOut",
	"Underline",
	"Blink",
}

type PadDirection int

const (
	LEFT PadDirection = iota
	CENTER
	RIGHT
)

type GlobalConfig struct {
	maxcolnamewidth int
	startposx       int
	allwidth        int
}

var globalconfig GlobalConfig

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

func ShowANSIColorBar(profile termenv.Profile) string {
	nbcolors := 16

	// ###################################################################
	// # Band colors
	// ###################################################################
	bandsize := globalconfig.allwidth / (nbcolors / 2)
	spaces := (globalconfig.allwidth - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, StrPad(RIGHT, "", globalconfig.startposx+spaces))
	for idx := 0; idx < nbcolors/2; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, StrPad(RIGHT, "", globalconfig.startposx+spaces))
	for idx := nbcolors / 2; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, "\n")

	return output.String()
}

// Show colors panel
func ShowANSI16ColorsPanel(profile termenv.Profile) string {
	// Show ANSI Colors
	nbcolors := 16
	sep := ""

	// ###################################################################
	// # Colors name
	// ###################################################################
	// Print colorname (in same color)
	output := bytes.NewBufferString("")
	fmt.Fprintf(output, StrPad(RIGHT, "ANSI", globalconfig.startposx))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s %s", sep, StrPad(CENTER, strings.ToUpper(fmt.Sprintf("%02x", idx)), globalconfig.maxcolnamewidth))
	}
	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, StrPad(RIGHT, "Color", globalconfig.startposx))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s %s", sep, StrPad(CENTER, COLORNAMES[idx], globalconfig.maxcolnamewidth))
	}
	fmt.Fprintf(output, "\n")

	// ###################################################################
	// # Colors block
	// ###################################################################
	for row, _ := range COLORNAMES {
		fmt.Fprintf(output, StrPad(RIGHT, COLORNAMES[row], globalconfig.startposx))
		for col := 0; col < nbcolors; col += 1 {
			fmt.Fprintf(output, "%s %s", sep, termenv.String(StrPad(CENTER, "•••", globalconfig.maxcolnamewidth)).Foreground(profile.Color(strconv.Itoa(row))).Background(profile.Color(strconv.Itoa(col))))
		}
		fmt.Fprintf(output, "\n")
	}

	return output.String()
}

// Show colors panel
func ShowTextStylePanel(profile termenv.Profile) string {
	nbcolors := 16
	sep := ""

	output := bytes.NewBufferString("")
	for row, _ := range TEXTSTYLES {
		fmt.Fprintf(output, StrPad(RIGHT, TEXTSTYLES[row], globalconfig.startposx))
		for col := 0; col < nbcolors; col += 1 {
			termfunc := termenv.String(StrPad(CENTER, "ABC", globalconfig.maxcolnamewidth)).Foreground(profile.Color(strconv.Itoa(col)))
			switch TEXTSTYLES[row] {
			case "Bold":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Bold())
			case "Faint":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Faint())
			case "Italic":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Italic())
			case "CrossOut":
				fmt.Fprintf(output, "%s %s", sep, termfunc.CrossOut())
			case "Underline":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Underline())
			case "Overline":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Overline())
			case "Blink":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Blink())
			default:
				fmt.Fprintf(output, "%s %s", sep, termfunc)
			}
		}
		fmt.Fprintf(output, "\n")

	}

	return output.String()
}

func ShowGrayColorsPanel(profile termenv.Profile) string {
	nbcolors := 255 - 232
	bandsize := globalconfig.allwidth / ((nbcolors + 1) / 2)
	spaces := (globalconfig.allwidth - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, StrPad(RIGHT, "", globalconfig.startposx+spaces-2))
	for idx := 232; idx <= 232+nbcolors/2; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(StrPad(CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("15")),
		)
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, StrPad(RIGHT, "", globalconfig.startposx+spaces-2))
	for idx := 232 + nbcolors/2 + 1; idx <= 255; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(StrPad(CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("0")),
		)
	}
	fmt.Fprintf(output, "\n")

	return output.String()
}

func GetGlobalConf(colornames []string, textstyles []string) GlobalConfig {
	globalconfig := GlobalConfig{}

	globalconfig.maxcolnamewidth = 0
	for _, colname := range colornames {
		size := len(colname)
		if size > globalconfig.maxcolnamewidth {
			globalconfig.maxcolnamewidth = size
		}
	}

	globalconfig.startposx = 0
	for _, colname := range textstyles {
		size := len(colname)
		if size > globalconfig.startposx {
			globalconfig.startposx = size
		}
	}

	if globalconfig.startposx < globalconfig.maxcolnamewidth {
		globalconfig.startposx = globalconfig.maxcolnamewidth
	}

	globalconfig.allwidth = ((globalconfig.maxcolnamewidth + 1) * 16) - 2

	return globalconfig

}

func OutputBar() string {
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
	if err != nil {
		panic(err)
	}
	defer restoreConsole()

	globalconfig = GetGlobalConf(COLORNAMES, TEXTSTYLES)
	output := ShowANSIColorBar(termenv.ANSI)
	output += ShowANSI16ColorsPanel(termenv.ANSI)
	output += ShowTextStylePanel(termenv.ANSI)
	output += ShowGrayColorsPanel(termenv.ANSI256)

	return output
}

func main() {
	fmt.Print(OutputBar())
}
