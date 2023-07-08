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
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/badele/tea-colors/internal/pkg/ansi"
	"github.com/badele/tea-colors/internal/pkg/scheme"
	importer "github.com/badele/tea-colors/internal/pkg/scheme/importer/iterm2"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
/////////////////////////////////////////////////
// Model / View
/////////////////////////////////////////////////
*/

var docStyle = lipgloss.NewStyle().
	Margin(1, 2).
	BorderForeground(lipgloss.Color("8")).
	BorderStyle(lipgloss.NormalBorder())

type SELECTEDFOCUS int

const (
	FOCUSLIST SELECTEDFOCUS = iota
	FOCUSPREVIEW
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type modelpreview struct {
	content    string
	ready      bool
	termwidth  int
	termheight int
	viewport   viewport.Model
	listview   list.Model
	focusState SELECTEDFOCUS

	previous_state  SELECTEDFOCUS
	previous_offset int
}

func (m modelpreview) Init() tea.Cmd {
	return nil
}

var (
	textColor = lipgloss.AdaptiveColor{Light: "0", Dark: "4"}
	lineColor = lipgloss.AdaptiveColor{Light: "0", Dark: "8"}

	focusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Padding(0, 1).
			MarginRight(1).
			Render

	unfocusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Padding(0, 1).
			MarginRight(1).
			Render

	verticalspace = lipgloss.NewStyle().
			BorderForeground(lineColor).
			MarginRight(2).
			Height(8).
			Width(2).
			Render

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func (m modelpreview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab":
			m.focusState = (m.focusState + 1) % 2
		}

	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := footerHeight + 1

		if !m.ready {
			m.termwidth = msg.Width
			m.termheight = msg.Height

			h, v := docStyle.GetFrameSize()
			m.listview.SetSize(m.termwidth-ansi.GetFullSize()-h, msg.Height-v+2)

			m.viewport = viewport.New(ansi.GetFullSize(), msg.Height-verticalMarginHeight-1)
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	// m.viewport, cmd = m.viewport.Update(msg)
	// cmds = append(cmds, cmd)

	if m.focusState != m.previous_state {
		if m.focusState == FOCUSLIST {
			m.previous_offset = m.viewport.YOffset
		} else {
			m.viewport.SetContent(m.content)
			m.viewport.SetYOffset(m.previous_offset)
		}

		m.previous_state = m.focusState
	}

	if m.focusState == FOCUSLIST {
		m.listview, cmd = m.listview.Update(msg)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m modelpreview) GetMenuText(text string, focused bool, width int) string {
	style := unfocusStyle
	if focused {
		style = focusStyle
	}

	menutext := fmt.Sprintf("│ %s ├", text)
	size := len(menutext)

	menutext += strings.Repeat("─", width-size+2)

	return style(menutext)
}

func (m modelpreview) GetFocusedMenu() string {
	previewmenu := m.GetMenuText("Preview theme", m.focusState == FOCUSPREVIEW, ansi.GetFullSize())
	listmenu := m.GetMenuText("Select colorscheme", m.focusState == FOCUSLIST, m.termwidth-ansi.GetFullSize())
	return lipgloss.JoinHorizontal(lipgloss.Left, previewmenu, listmenu)
}

func (m modelpreview) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	preview := fmt.Sprintf("%s\n%s", m.viewport.View(), m.footerView())
	focusmenu := m.GetFocusedMenu()

	previwbloc := lipgloss.JoinHorizontal(lipgloss.Center, preview, verticalspace(" "), m.listview.View())

	return lipgloss.JoinVertical(lipgloss.Top, focusmenu, "", previwbloc)
}

func (m modelpreview) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func TitleGenerator(text string, boxwidth int) string {
	output := bytes.NewBufferString("")

	size := len(text) + 2
	leftspaces := (boxwidth - size) / 2
	rightspaces := boxwidth - leftspaces - size

	lineStyle := lipgloss.NewStyle().Foreground(lineColor)

	textStyle := lipgloss.NewStyle().
		Bold(true)

	leftline := lineStyle.Render(strings.Repeat("─", leftspaces))
	rightline := lineStyle.Render(strings.Repeat("─", rightspaces))

	title := fmt.Sprintf("\n%s%s%s\n\n", leftline, textStyle.Render(" "+text+" "), rightline)

	fmt.Fprintf(output, title)

	return output.String()
}

func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getAnsiContentFile(filename string) string {
	lines, err := readLinesFromFile(filename)

	if err != nil || len(lines) < 2 {
		return ""
	}

	output := TitleGenerator(lines[0], ansi.GetFullSize())
	for _, line := range lines[1:] {
		output += fmt.Sprintf("%s\n", line)
	}

	output = strings.Replace(output, "\t", "   ", -1)

	return output
}

func OutputBar() string {
	output := ansi.GetANSIColorBar()
	output += ansi.GetANSI16ColorsPanel()
	output += ansi.GetTextStylePanel()
	output += ansi.GetGrayColorsPanel()

	output += getAnsiContentFile("samples/duf.ans")
	output += getAnsiContentFile("samples/exa.ans")
	output += getAnsiContentFile("samples/git-status.ans")
	output += getAnsiContentFile("samples/man.ans")

	// output += "\n\n"

	// style := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("#448844")).
	// 	Foreground(lipgloss.Color("#884488"))

	// output = strings.Replace(output, "44m", fmt.Sprintf("48;2;%d;%d;%dm", r/256, g/256, b/256), -1)
	// output = strings.Replace(output, "[91", fmt.Sprintf("[38;2;%d;%d;%d", r/256, g/256, b/256), -1)
	// output = strings.Replace(output, "91m", fmt.Sprintf("38;2;%d;%d;%dm", r/256, g/256, b/256), -1)
	// output = strings.Replace(output, "91;", fmt.Sprintf("38;2;%d;%d;%d;", r/256, g/256, b/256), -1)

	// you must replace ( [0m or [m ) by [;{BG}
	// output = strings.Replace(output, "{BG}", fmt.Sprintf("48;2;%d;%d;%dm", r/256, g/256, b/256), -1)

	return output
}

func imports() {
	// Declar interface type
	var schemeImporter scheme.Importer

	newsschemes := scheme.Schemes{}

	// Wezterm
	schemeImporter = &importer.SchemeWezterm{}
	newsschemes = scheme.Import(newsschemes, schemeImporter)

	// // Base16
	// schemeImporter = &importer.SchemeBase16{}
	// newsschemes = scheme.Import(newsschemes, schemeImporter)

	// // Gogh
	// schemeImporter = &importer.SchemeGogh{}
	// newsschemes = scheme.Import(newsschemes, schemeImporter)

	// newsschemes.Read("/tmp/scheme.txt")

	sort.Sort(newsschemes)

	newsschemes.Write("/tmp/scheme.txt")

}

func GetItemBarColor(colors []string) string {
	bar := ""

	// If not colors, generate ANSI color bar
	if colors == nil {
		for i := 0; i < 16; i++ {
			colors = append(colors, fmt.Sprintf("%d", i))
		}
	}

	for i := 0; i < 16; i++ {
		colstyle := lipgloss.NewStyle().
			SetString(" ").
			Background(lipgloss.Color(colors[i]))

		bar += fmt.Sprint(colstyle)
	}

	return bar
}

func main() {

	// imports()
	// return

	// test := importer.SchemeBase16{}
	// links := test.GetRessouresList()
	// fmt.Println(links)
	// return
	// fmt.Print(OutputBar())
	// return

	// Load some text for our viewport
	content := strings.Repeat(OutputBar()+"\n", 5)

	newsschemes := scheme.Schemes{}
	newsschemes.Read("schemes.txt")
	items := []list.Item{
		item{title: "Current", desc: GetItemBarColor(nil)},
	}

	for _, scheme := range newsschemes {
		colors := []string{}
		colors = append(colors, scheme.Normal...)
		colors = append(colors, scheme.Brights...)

		items = append(items, item{title: scheme.Name, desc: GetItemBarColor(colors)})
	}

	m := modelpreview{
		content:  string(content),
		listview: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.listview.SetShowTitle(false)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
