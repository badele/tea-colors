package importer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/badele/tea-colors/internal/pkg/scheme"
	"github.com/badele/tea-colors/internal/pkg/tools"
	"github.com/gocolly/colly"
)

// Wezterm
type SchemeWezterm struct {
	Scheme string

	Colors struct {
		Ansi             []string `toml:"ansi"`
		Brights          []string `toml:"brights"`
		Foreground       string   `toml:"foreground"`
		Background       string   `toml:"background"`
		CursorForeground string   `toml:"cursor_fg"`
		CursorBg         string   `toml:"cursor_bg"`
		SelectionFg      string   `toml:"selection_fg"`
		SelectionBg      string   `toml:"selection_bg"`
	}
}

const (
	wezterm_source                = "wezterm"
	wezterm_rawGithubUrl   string = "https://raw.githubusercontent.com"
	wezterm_mainscheme_url string = "https://github.com/mbadolato/iTerm2-Color-Schemes/tree/master/wezterm"
)

func (s *SchemeWezterm) GetRessouresList() []string {
	links := []string{}

	c := colly.NewCollector()

	r, _ := regexp.Compile(`/(.*?\.toml)$`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if e.Attr("class") == "js-navigation-open Link--primary" {
			// e.Request.Visit(e.Attr("href"))
			link := e.Attr("href")

			if match := r.MatchString(link); match == true {
				filename := strings.Replace(r.FindStringSubmatch(link)[1], "/blob/", "/", -1)
				absoluteURL := fmt.Sprintf("%s/%s", wezterm_rawGithubUrl, filename)

				links = append(links, absoluteURL)
			}
		}
	})

	mainURL := wezterm_mainscheme_url
	c.Visit(mainURL)

	return links
}

func (s *SchemeWezterm) ImportRessources(ressources []string, schemes scheme.Schemes) scheme.Schemes {
	var schemename string

	for _, url := range ressources {
		importscheme := SchemeWezterm{}

		// Get content
		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		// Unmarshal toml
		_, err = toml.Decode(resp, &importscheme)
		tools.CheckError(err)

		// Get scheme name from URL
		r, _ := regexp.Compile(`.*/([^./]*)\.(.+)$`)

		if match := r.MatchString(url); match == true {
			groups := r.FindStringSubmatch(url)
			schemename = groups[1]
			schemename = strings.Replace(schemename, "%20", " ", -1)
			schemename = strings.Replace(schemename, "%2B", "+", -1)
		}

		newsscheme := scheme.Scheme{
			Source:      "iterm2",
			Application: "wezterm",
			Name:        schemename,
			Author:      "Undefined",
			Foreground:  scheme.GetHexaColor(importscheme.Colors.Foreground),
			Background:  scheme.GetHexaColor(importscheme.Colors.Background),
			Normal:      strings.Split(strings.ToUpper(strings.Join(importscheme.Colors.Ansi, ",")), ","),
			Brights:     strings.Split(strings.ToUpper(strings.Join(importscheme.Colors.Brights, ",")), ","),
			CursorFg:    scheme.GetHexaColor(importscheme.Colors.CursorForeground),
			CursorBg:    scheme.GetHexaColor(importscheme.Colors.CursorBg),
			SelectionFg: scheme.GetHexaColor(importscheme.Colors.SelectionFg),
			SelectionBg: scheme.GetHexaColor(importscheme.Colors.SelectionBg),
		}
		schemes = append(schemes, newsscheme)
	}

	return schemes
}
