package ansi

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/badele/tea-colors/internal/pkg/tools"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
var COLORNAMES = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

var TEXTSTYLES = []string{
	"normal",
	"bold",
	"faint",
	"italic",
	"overline",
	"crossout",
	"underline",
	"blink",
}

type GlobalConfig struct {
	MaxColorNameWidth int
	StartPosX         int
	AllWidth          int
	FullSize          int
}

// Get max ANSI color names width
func GetMaxColorNameWidth() int {
	value := 0
	for _, colname := range COLORNAMES {
		size := len(colname)
		if size > value {
			value = size
		}
	}

	return value
}

func GetMaxStyleNameWidth() int {
	// Get max ANSI style names width
	value := 0
	for _, colname := range TEXTSTYLES {
		size := len(colname)
		if size > value {
			value = size
		}
	}

	return value
}

func GetStartPosX() int {
	value := GetMaxColorNameWidth()
	if GetMaxStyleNameWidth() > GetMaxColorNameWidth() {
		value = GetMaxStyleNameWidth()
	}

	return value
}

func GetAllWidth() int {
	return ((GetMaxColorNameWidth() + 1) * 16) // Ignore last space
}

func GetFullSize() int {
	return GetStartPosX() + 1 + GetAllWidth()
}

// func GetGlobalConf(colornames []string, textstyles []string) GlobalConfig {
// 	globalconfig := GlobalConfig{}

// 	// // Get max ANSI color names width
// 	// GetMaxColorNameWidth() = 0
// 	// for _, colname := range colornames {
// 	// 	size := len(colname)
// 	// 	if size > GetMaxColorNameWidth() {
// 	// 		GetMaxColorNameWidth() = size
// 	// 	}
// 	// }

// 	// // Get max ANSI style names width
// 	// GetMaxStyleNameWidth() = 0
// 	// for _, colname := range textstyles {
// 	// 	size := len(colname)
// 	// 	if size > GetMaxStyleNameWidth() {
// 	// 		GetMaxStyleNameWidth() = size
// 	// 	}
// 	// }

// 	// Get max widh (MaxColorNameWidth, MaxStyleNameWidth)
// 	// GetStartPosX() = GetMaxColorNameWidth()
// 	// if GetMaxStyleNameWidth() > GetMaxColorNameWidth() {
// 	// 	GetStartPosX() = GetMaxStyleNameWidth()
// 	// }

// 	// Box color preview width
// 	// GetAllWidth() = ((GetMaxColorNameWidth() + 1) * 16) // Ignore last space

// 	// Max preview size
// 	// globalconfig.FullSize = GetStartPosX() + 1 + GetAllWidth()

// 	return globalconfig
// }

func GetANSIColorBar(profile termenv.Profile) string {
	nbcolors := 16

	// ###################################################################
	// # Band colors
	// ###################################################################
	bandsize := GetAllWidth() / (nbcolors / 2)
	spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 0; idx < nbcolors/2; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := nbcolors / 2; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(strings.Repeat(" ", bandsize)).Background(profile.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, "\n")

	return output.String()
}

// Show colors panel
func GetANSI16ColorsPanel(profile termenv.Profile) string {
	// Show ANSI Colors
	nbcolors := 16
	sep := ""

	// ###################################################################
	// # Colors name
	// ###################################################################
	output := bytes.NewBufferString("")

	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, lipgloss.JoinHorizontal(lipgloss.Center, tools.StrPad(tools.RIGHT, "", GetMaxStyleNameWidth()), tools.StrPad(tools.CENTER, "Normal", GetAllWidth()/2), " ", tools.StrPad(tools.CENTER, "Bright", GetAllWidth()/2)))
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "ANSI", GetMaxStyleNameWidth()))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s %s", sep, tools.StrPad(tools.CENTER, strings.ToUpper(fmt.Sprintf("%02x", idx)), GetMaxColorNameWidth()))
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "Color", GetMaxStyleNameWidth()))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s %s", sep, tools.StrPad(tools.CENTER, COLORNAMES[idx%(nbcolors/2)], GetMaxColorNameWidth()))
	}
	fmt.Fprintf(output, "\n")

	// ###################################################################
	// # Colors block
	// ###################################################################

	for row := 0; row < nbcolors; row++ {
		fmt.Fprintf(output, tools.StrPad(tools.RIGHT, COLORNAMES[row%(nbcolors/2)], GetMaxStyleNameWidth()))
		for col := 0; col < nbcolors; col += 1 {
			fmt.Fprintf(output, "%s %s", sep, termenv.String(tools.StrPad(tools.CENTER, "•••", GetMaxColorNameWidth())).Foreground(profile.Color(strconv.Itoa(row))).Background(profile.Color(strconv.Itoa(col))))
		}
		fmt.Fprintf(output, "\n")
	}

	return output.String()
}

// Show colors panel
func GetTextStylePanel(profile termenv.Profile) string {
	nbcolors := 16
	sep := ""

	output := bytes.NewBufferString("")
	for row, _ := range TEXTSTYLES {
		fmt.Fprintf(output, tools.StrPad(tools.RIGHT, TEXTSTYLES[row], GetMaxStyleNameWidth()))
		for col := 0; col < nbcolors; col += 1 {
			termfunc := termenv.String(tools.StrPad(tools.CENTER, "ABC", GetMaxColorNameWidth())).Foreground(profile.Color(strconv.Itoa(col)))
			switch TEXTSTYLES[row] {
			case "bold":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Bold())
			case "faint":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Faint())
			case "italic":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Italic())
			case "crossOut":
				fmt.Fprintf(output, "%s %s", sep, termfunc.CrossOut())
			case "underline":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Underline())
			case "overline":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Overline())
			case "blink":
				fmt.Fprintf(output, "%s %s", sep, termfunc.Blink())
			default:
				fmt.Fprintf(output, "%s %s", sep, termfunc)
			}
		}
		fmt.Fprintf(output, "\n")

	}

	return output.String()
}

func GetGrayColorsPanel(profile termenv.Profile) string {
	nbcolors := 256 - 232
	// bandsize := GetAllWidth() / ((nbcolors + 1) / 2)
	// spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2
	bandsize := GetAllWidth() / (nbcolors / 2)
	spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 232; idx <= 232+nbcolors/2-1; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("15")),
		)
	}
	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 232 + nbcolors/2; idx <= 255; idx += 1 {
		fmt.Fprintf(output, "%s", termenv.String(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize)).
			Background(profile.Color(strconv.Itoa(idx))).
			Foreground(profile.Color("0")),
		)
	}
	fmt.Fprintf(output, "\n")

	return output.String()
}
