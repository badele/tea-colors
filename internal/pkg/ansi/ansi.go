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
		// "overline",
		// "crossout",
		// "underline",
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
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 0; idx < nbcolors/2; idx += 1 {
		fmt.Fprintf(output, "%s", style.Background(lipgloss.Color(strconv.Itoa(idx))))

	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := nbcolors / 2; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s", style.Background(lipgloss.Color(strconv.Itoa(idx))))
	}
	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetFullSize()))

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

	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, lipgloss.JoinHorizontal(lipgloss.Center, tools.StrPad(tools.RIGHT, "", GetMaxStyleNameWidth()), tools.StrPad(tools.CENTER, "Normal", GetAllWidth()/2), " ", tools.StrPad(tools.CENTER, "Bright", GetAllWidth()/2)))
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "ANSI", GetStartPosX()))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s%s", " ", tools.StrPad(tools.CENTER, strings.ToUpper(fmt.Sprintf("%02x", idx)), GetStartPosX()))
	}
	// fmt.Fprintf(output, " "))
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "Color", GetStartPosX()))
	for idx := 0; idx < nbcolors; idx += 1 {
		fmt.Fprintf(output, "%s%s", " ", tools.StrPad(tools.CENTER, COLORNAMES[idx%(nbcolors/2)], GetStartPosX()))
	}
	fmt.Fprintf(output, "\n")

	// ###################################################################
	// # Colors block
	// ###################################################################

	style := lipgloss.NewStyle().
		SetString(tools.StrPad(tools.CENTER, "•••", GetMaxColorNameWidth()))

	for row := 0; row < nbcolors; row++ {
		fmt.Fprintf(output, tools.StrPad(tools.RIGHT, COLORNAMES[row%(nbcolors/2)], GetStartPosX()))
		for col := 0; col < nbcolors; col += 1 {
			fmt.Fprintf(output, "%s%s", " ", style.Foreground(lipgloss.Color(strconv.Itoa(row))).Background(lipgloss.Color(strconv.Itoa(col))))
		}
		fmt.Fprintf(output, "\n")
	}

	return output.String()
}

// Show colors panel
func GetTextStylePanel() string {
	nbcolors := 16

	textcontent := tools.StrPad(tools.CENTER, tools.StrPad(tools.CENTER, "ABC", GetMaxColorNameWidth()), GetMaxColorNameWidth())

	output := bytes.NewBufferString("")
	for row, _ := range TEXTSTYLES {
		fmt.Fprintf(output, tools.StrPad(tools.RIGHT, TEXTSTYLES[row], GetStartPosX()))
		for col := 0; col < nbcolors; col += 1 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(col)))
			switch TEXTSTYLES[row] {
			case "bold":
				fmt.Fprintf(output, "%s%s", " ", style.Bold(true).Render(textcontent))
			case "faint":
				fmt.Fprintf(output, "%s%s", " ", style.Faint(true).Render(textcontent))
			case "italic":
				fmt.Fprintf(output, "%s%s", " ", style.Italic(true).Render(textcontent))
			// case "crossout":
			// 	fmt.Fprintf(output, "%s%s", " ", style.Strikethrough(true).Render(textcontent))
			// case "underline":
			// 	fmt.Fprintf(output, "%s%s", " ", style.Underline(true).Render(textcontent))
			case "blink":
				fmt.Fprintf(output, "%s%s", " ", style.Blink(true).Render(textcontent))
			default:
				fmt.Fprintf(output, "%s%s", " ", style.Render(textcontent))
			}
		}
		fmt.Fprintf(output, "\n")
	}

	return output.String()
}

func GetGrayColorsPanel() string {
	nbcolors := 256 - 232 // 24
	bandsize := GetAllWidth() / (nbcolors / 2)
	spaces := (GetAllWidth() - (bandsize * (nbcolors / 2))) / 2

	output := bytes.NewBufferString("")

	// Line 1
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 232; idx <= 232+nbcolors/2-1; idx += 1 {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color(strconv.Itoa(idx))).
			SetString(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize))

		fmt.Fprintf(output, "%s", style)
	}
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetFullSize()-(bandsize*(nbcolors/2))-GetStartPosX()-spaces-2))
	fmt.Fprintf(output, "\n")

	// Line 2
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetStartPosX()+1+spaces))
	for idx := 232 + nbcolors/2; idx <= 255; idx += 1 {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color(strconv.Itoa(idx))).
			SetString(tools.StrPad(tools.CENTER, strconv.Itoa(idx), bandsize))

		fmt.Fprintf(output, "%s", style)
	}
	fmt.Fprintf(output, tools.StrPad(tools.RIGHT, "", GetFullSize()-(bandsize*(nbcolors/2))-GetStartPosX()-spaces-2))

	fmt.Fprintf(output, "\n")

	return output.String()
}
