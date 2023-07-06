package ansi

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/badele/tea-colors/internal/pkg/tools"
	"github.com/charmbracelet/lipgloss"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors

// var DEFAULTBACKGROUND = "#2D2927"
var (
	ANSILINE          = "[K"
	DEFAULTBACKGROUND = "#004427"
	DEFAULTFOREGROUND = "#FBF1C7"

	BG = lipgloss.NewStyle().
		Background(lipgloss.Color(DEFAULTBACKGROUND))
	FG = lipgloss.NewStyle().
		Foreground(lipgloss.Color(DEFAULTBACKGROUND))

	COLORNAMES = []string{
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
	}

	TEXTSTYLES = []string{
		"normal",
		"bold",
		"faint",
		"italic",
		"overline",
		"crossout",
		"underline",
		"blink",
	}
)

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

func GetANSIColorBar() string {
	nbcolors := 16

	// ###################################################################
	// # Band colors
	// ###################################################################
	bandsize := GetAllWidth() / (nbcolors / 2)
	spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2
	style := lipgloss.NewStyle().
		SetString(strings.Repeat(" ", bandsize))

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces)))
	for idx := 0; idx < nbcolors/2; idx += 1 {
		fmt.Fprintf(output, "%s", style.Background(lipgloss.Color(strconv.Itoa(idx))))

	}
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")

	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces)))
	for idx := nbcolors / 2; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s", style.Background(lipgloss.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")
	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetFullSize())))

	return output.String()
}

// Show colors panel
func GetANSI16ColorsPanel() string {
	// Show ANSI Colors
	nbcolors := 16
	// ###################################################################
	// # Colors name
	// ###################################################################
	output := bytes.NewBufferString("")

	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")
	fmt.Fprintf(output, lipgloss.JoinHorizontal(lipgloss.Center, BG.Render(tools.StrPad(tools.RIGHT, "", GetMaxStyleNameWidth())), BG.Render(tools.StrPad(tools.CENTER, "Normal", GetAllWidth()/2)), BG.Render(" "), BG.Render(tools.StrPad(tools.CENTER, "Bright", GetAllWidth()/2))))
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")

	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "ANSI", GetMaxStyleNameWidth())))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s%s", BG.Render(" "), BG.Render(tools.StrPad(tools.CENTER, strings.ToUpper(fmt.Sprintf("%02x", idx)), GetMaxColorNameWidth())))
	}
	// fmt.Fprintf(output, BG.Render(" "))
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")

	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "Color", GetMaxStyleNameWidth())))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s%s", BG.Render(" "), BG.Render(tools.StrPad(tools.CENTER, COLORNAMES[idx%(nbcolors/2)], GetMaxColorNameWidth())))
	}
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")

	// ###################################################################
	// # Colors block
	// ###################################################################

	style := lipgloss.NewStyle().
		SetString(tools.StrPad(tools.CENTER, "•••", GetMaxColorNameWidth()))

	for row := 0; row < nbcolors; row++ {
		fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, COLORNAMES[row%(nbcolors/2)], GetMaxStyleNameWidth())))
		for col := 0; col < nbcolors; col += 1 {
			fmt.Fprintf(output, "%s%s", BG.Render(" "), style.Foreground(lipgloss.Color(strconv.Itoa(row))).Background(lipgloss.Color(strconv.Itoa(col))))
		}
		fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")
	}

	return output.String()
}

// Show colors panel
func GetTextStylePanel() string {
	nbcolors := 16

	style := lipgloss.NewStyle().
		SetString(tools.StrPad(tools.CENTER, tools.StrPad(tools.CENTER, "ABC", GetMaxColorNameWidth()), GetMaxColorNameWidth()))

	output := bytes.NewBufferString("")
	for row, _ := range TEXTSTYLES {
		fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, TEXTSTYLES[row], GetMaxStyleNameWidth())))
		for col := 0; col < nbcolors; col += 1 {
			termfunc := style.Foreground(lipgloss.Color(strconv.Itoa(col))).Background(lipgloss.Color(DEFAULTBACKGROUND))
			switch TEXTSTYLES[row] {
			case "bold":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Bold(true))
			case "faint":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Faint(true))
			case "italic":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Italic(true))
			case "crossOut":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Strikethrough(true))
			case "underline":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Underline(true))
			// case "overline":
			// 	fmt.Fprintf(output, "%s%s", bg.Render(" "), termfunc.Overline())
			case "blink":
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc.Blink(true))
			default:
				fmt.Fprintf(output, "%s%s", BG.Render(" "), termfunc)
			}
		}
		fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")
	}

	return output.String()
}

func GetGrayColorsPanel() string {
	nbcolors := 256 - 232 // 24
	bandsize := GetAllWidth() / (nbcolors / 2)
	spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")
	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces)))
	for idx := 232; idx <= 232+nbcolors/2-1; idx += 1 {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color(strconv.Itoa(idx))).
			SetString(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize))

		fmt.Fprintf(output, "%s", style)
	}
	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetFullSize()-(bandsize*(nbcolors/2))-GetStartPosX()-spaces-2)))
	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")

	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces)))
	for idx := 232 + nbcolors/2; idx <= 255; idx += 1 {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color(strconv.Itoa(idx))).
			SetString(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize))

		fmt.Fprintf(output, "%s", style)
	}
	fmt.Fprintf(output, BG.Render(tools.StrPad(tools.RIGHT, "", GetFullSize()-(bandsize*(nbcolors/2))-GetStartPosX()-spaces-2)))

	fmt.Fprintf(output, BG.Render(ANSILINE)+"\n")
	fmt.Fprintf(output, BG.Render(ANSILINE))

	return output.String()
}
