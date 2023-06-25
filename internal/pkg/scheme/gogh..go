package scheme

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/badele/tea-colors/internal/pkg/tools"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-yaml/yaml"
	"github.com/gocolly/colly"
)

type SchemeGogh16 struct {
	Scheme string `yaml:"name"`

	Color01 string `yaml:"color_01"`
	Color02 string `yaml:"color_02"`
	Color03 string `yaml:"color_03"`
	Color04 string `yaml:"color_04"`
	Color05 string `yaml:"color_05"`
	Color06 string `yaml:"color_06"`
	Color07 string `yaml:"color_07"`
	Color08 string `yaml:"color_08"`
	Color09 string `yaml:"color_09"`
	Color10 string `yaml:"color_10"`
	Color11 string `yaml:"color_11"`
	Color12 string `yaml:"color_12"`
	Color13 string `yaml:"color_13"`
	Color14 string `yaml:"color_14"`
	Color15 string `yaml:"color_15"`
	Color16 string `yaml:"color_16"`

	Background string `yaml:"background"`
	Foreground string `yaml:"foreground"`
	Cursor     string `yaml:"cursor"`
}

const (
	gogh_source                = "gogh"
	gogh_rawGithubUrl   string = "https://raw.githubusercontent.com"
	gogh_mainscheme_url string = "https://github.com/Gogh-Co/Gogh/tree/master/themes"
)

// func getYamlList(url string) []string {
// 	links := []string{}

// 	c := colly.NewCollector()

// 	r, _ := regexp.Compile(`/(.*?\.yaml)$`)
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

// 		if e.Attr("class") == "js-navigation-open Link--primary" {
// 			// e.Request.Visit(e.Attr("href"))
// 			link := e.Attr("href")

// 			if match := r.MatchString(link); match == true {
// 				filename := strings.Replace(r.FindStringSubmatch(link)[1], "/blob/", "/", -1)
// 				absoluteURL := fmt.Sprintf("%s/%s", gogh_rawGithubUrl, filename)

// 				links = append(links, absoluteURL)
// 			}
// 		}
// 	})

// 	mainURL := url
// 	c.Visit(mainURL)

// 	return links
// }

func (s *SchemeGogh16) GetRessouresList() []string {
	links := []string{}

	c := colly.NewCollector()

	r, _ := regexp.Compile(`/(.*?\.yml)$`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if e.Attr("class") == "js-navigation-open Link--primary" {
			// e.Request.Visit(e.Attr("href"))
			link := e.Attr("href")

			if match := r.MatchString(link); match == true {
				filename := strings.Replace(r.FindStringSubmatch(link)[1], "/blob/", "/", -1)
				absoluteURL := fmt.Sprintf("%s/%s", gogh_rawGithubUrl, filename)

				links = append(links, absoluteURL)
			}
		}
	})

	mainURL := gogh_mainscheme_url
	c.Visit(mainURL)

	return links
}

func (s *SchemeGogh16) ImportRessources(ressources []string, schemes Schemes) Schemes {
	for _, url := range ressources {

		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		err = yaml.Unmarshal([]byte(resp), s)
		tools.CheckError(err)

		newsscheme := Scheme{
			Source: gogh_source,
			Name:   s.Scheme,
			Author: "Undefined",
			Colors: []lipgloss.Color{
				lipgloss.Color(s.Color01),
				lipgloss.Color(s.Color02),
				lipgloss.Color(s.Color03),
				lipgloss.Color(s.Color04),
				lipgloss.Color(s.Color05),
				lipgloss.Color(s.Color06),
				lipgloss.Color(s.Color07),
				lipgloss.Color(s.Color08),
				lipgloss.Color(s.Color09),
				lipgloss.Color(s.Color10),
				lipgloss.Color(s.Color11),
				lipgloss.Color(s.Color12),
				lipgloss.Color(s.Color13),
				lipgloss.Color(s.Color14),
				lipgloss.Color(s.Color15),
				lipgloss.Color(s.Color16),
			},
		}
		schemes = append(schemes, newsscheme)
	}

	return schemes
}
