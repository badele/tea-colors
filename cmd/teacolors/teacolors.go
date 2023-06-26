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

	"github.com/muesli/termenv"
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
	content string
	ready   bool

	termwidth  int
	termheight int
	viewport   viewport.Model
	list       list.Model
	focusState SELECTEDFOCUS
}

func (m modelpreview) Init() tea.Cmd {
	return nil
}

var (
	textColor = lipgloss.AdaptiveColor{Light: "0", Dark: "4"}
	lineColor = lipgloss.AdaptiveColor{Light: "0", Dark: "8"}

	focusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("5")).
			Padding(0, 1).
			MarginRight(1).
			Render

	unfocusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("0")).
			Padding(0, 1).
			MarginRight(1).
			Render

	lineStyle = lipgloss.NewStyle().
			Foreground(lineColor).
			Render
	textStyle = lipgloss.NewStyle().Foreground(textColor).
			Bold(true).
			Render

	verticalspace = lipgloss.NewStyle().
		// Border(lipgloss.NormalBorder(), false, true, false, false).
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
		// cmds []tea.Cmd
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
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight + 1

		//m.list.SetSize(msg.Width-ansi.GetFullSize(), msg.Height-v)

		if !m.ready {
			m.termwidth = msg.Width
			m.termheight = msg.Height

			h, v := docStyle.GetFrameSize()
			m.list.SetSize(m.termwidth-ansi.GetFullSize()-h, msg.Height-v)

			m.viewport = viewport.New(ansi.GetFullSize(), msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			// m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		// if useHighPerformanceRenderer {
		// 	// Render (or re-render) the whole viewport. Necessary both to
		// 	// initialize the viewport and when the window is resized.
		// 	//
		// 	// This is needed for high-performance rendering only.
		// 	cmds = append(cmds, viewport.Sync(m.viewport))
		// }
	}

	// Handle keyboard and mouse events in the viewport
	// m.viewport, cmd = m.viewport.Update(msg)
	// cmds = append(cmds, cmd)

	m.list, cmd = m.list.Update(msg)
	return m, cmd

	// return m, tea.Batch(cmds...)
}

func (m modelpreview) GetFocusMenu() string {

	liststyle := focusStyle
	previewstyle := unfocusStyle
	if m.focusState == FOCUSPREVIEW {
		liststyle = unfocusStyle
		previewstyle = focusStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, previewstyle("Preview theme"), liststyle("Select colorscheme"))
}

func (m modelpreview) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	preview := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	focusmenu := m.GetFocusMenu()

	previwbloc := lipgloss.JoinHorizontal(lipgloss.Center, preview, verticalspace(""), m.list.View())

	return lipgloss.JoinVertical(lipgloss.Top, focusmenu, previwbloc)
}

func (m modelpreview) headerView() string {
	title := titleStyle.Render("Mr. Pager & Mme. pagers")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
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

// func (m modellist) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return m, tea.Quit
// 		}
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)
// 	}

// 	var cmd tea.Cmd
// 	m.list, cmd = m.list.Update(msg)
// 	return m, cmd
// }

// func (m modellist) View() string {
// 	return docStyle.Render(m.list.View())
// }

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

// func WordWrap(text string, width int) string {
// 	output := bytes.NewBufferString("")

// 	posx := 0
// 	for _, character := range text {
// 		if character == '\n' {
// 			posx = -1
// 		} else {
// 			if (posx % width) == width-1 {
// 				posx = 0
// 				fmt.Fprintf(output, "\n")
// 			}
// 		}

// 		fmt.Fprintf(output, "%c", character)
// 		posx += 1
// 	}

// 	return output.String()
// }

func TitleGenerator(text string, boxwidth int) string {
	output := bytes.NewBufferString("")

	size := len(text)
	nbspaces := (boxwidth - size) / 2

	separatorline := lineStyle(strings.Repeat("─", nbspaces))

	fmt.Fprintln(output, "")
	fmt.Fprintf(output, "%s %s %s\n", separatorline, textStyle(text), separatorline)
	fmt.Fprintln(output, "")

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

	return output
}

func OutputBar() string {
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
	if err != nil {
		panic(err)
	}
	defer restoreConsole()

	output := ansi.GetANSIColorBar(termenv.ANSI)
	output += ansi.GetANSI16ColorsPanel(termenv.ANSI)
	output += ansi.GetTextStylePanel(termenv.ANSI)
	output += ansi.GetGrayColorsPanel(termenv.ANSI256)

	output += getAnsiContentFile("samples/duf.ans")
	// output += OutputAnsiContentFile("samples/exa.ans")
	// output += OutputAnsiContentFile("samples/git-status.ans")
	// output += OutputAnsiContentFile("samples/man.ans")

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
		item{title: "Current", desc: "[40m [0m[41m [0m[42m [0m[43m [0m[44m [0m[45m [0m[46m [0m[47m [0m[100m [0m[101m [0m[102m [0m[103m [0m[104m [0m[105m [0m[106m [0m[107m [0m"},
	}

	for _, scheme := range newsschemes {
		cols := []lipgloss.Style{}
		for i := 0; i < 8; i++ {
			colstyle := lipgloss.NewStyle().
				SetString(" ").
				Background(lipgloss.Color(scheme.Normal[i]))

			cols = append(cols, colstyle)
		}
		for i := 0; i < 8; i++ {
			colstyle := lipgloss.NewStyle().
				SetString(" ").
				Background(lipgloss.Color(scheme.Brights[i]))

			cols = append(cols, colstyle)
		}

		items = append(items, item{title: scheme.Name, desc: fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s", cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7], cols[8], cols[9], cols[10], cols[11], cols[12], cols[13], cols[14], cols[15])})
	}

	m := modelpreview{
		content: string(content),
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.SetShowTitle(false)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
