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

type SchemeBase16 struct {
	Scheme string `yaml:"scheme"`
	Author string `yaml:"author"`

	Base00 string `yaml:"base00"`
	Base01 string `yaml:"base01"`
	Base02 string `yaml:"base02"`
	Base03 string `yaml:"base03"`
	Base04 string `yaml:"base04"`
	Base05 string `yaml:"base05"`
	Base06 string `yaml:"base06"`
	Base07 string `yaml:"base07"`
	Base08 string `yaml:"base08"`
	Base09 string `yaml:"base09"`
	Base0A string `yaml:"base0A"`
	Base0B string `yaml:"base0B"`
	Base0C string `yaml:"base0C"`
	Base0D string `yaml:"base0D"`
	Base0E string `yaml:"base0E"`
	Base0F string `yaml:"base0F"`
}

const (
	base16_source                = "base16"
	base16_rawGithubUrl   string = "https://raw.githubusercontent.com"
	base16_mainscheme_url string = "chriskempson/base16-schemes-source/main/list.yaml"
)

func (s *SchemeBase16) getYamlList(url string) []string {
	links := []string{}

	c := colly.NewCollector()

	r, _ := regexp.Compile(`/(.*?\.yaml)$`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if e.Attr("class") == "js-navigation-open Link--primary" {
			// e.Request.Visit(e.Attr("href"))
			link := e.Attr("href")

			if match := r.MatchString(link); match == true {
				filename := strings.Replace(r.FindStringSubmatch(link)[1], "/blob/", "/", -1)
				absoluteURL := fmt.Sprintf("%s/%s", base16_rawGithubUrl, filename)

				links = append(links, absoluteURL)
			}
		}
	})

	mainURL := url
	c.Visit(mainURL)

	return links
}

func (s *SchemeBase16) GetRessouresList() []string {
	links := []string{}

	urls, err := tools.GetHttpContent(base16_rawGithubUrl + "/" + base16_mainscheme_url)
	tools.CheckError(err)

	urllist := strings.Split(urls, "\n")
	r, _ := regexp.Compile(`([^:]+): *(.*)`)
	for _, url := range urllist {
		if match := r.MatchString(url); match == true {
			groups := r.FindStringSubmatch(url)

			schemelinks := s.getYamlList(groups[2])
			links = append(links, schemelinks...)
		}
	}

	return links
}

func (s *SchemeBase16) ImportRessources(ressources []string, schemes Schemes) Schemes {
	for _, url := range ressources {

		resp, err := tools.GetHttpContent(url)
		tools.CheckError(err)

		err = yaml.Unmarshal([]byte(resp), s)
		tools.CheckError(err)

		newsscheme := Scheme{
			Source: base16_source,
			Name:   s.Scheme,
			Author: s.Author,
			Colors: []lipgloss.Color{
				lipgloss.Color(getHexaColor(s.Base00)),
				lipgloss.Color(getHexaColor(s.Base01)),
				lipgloss.Color(getHexaColor(s.Base02)),
				lipgloss.Color(getHexaColor(s.Base03)),
				lipgloss.Color(getHexaColor(s.Base04)),
				lipgloss.Color(getHexaColor(s.Base05)),
				lipgloss.Color(getHexaColor(s.Base06)),
				lipgloss.Color(getHexaColor(s.Base07)),
				lipgloss.Color(getHexaColor(s.Base08)),
				lipgloss.Color(getHexaColor(s.Base09)),
				lipgloss.Color(getHexaColor(s.Base0A)),
				lipgloss.Color(getHexaColor(s.Base0B)),
				lipgloss.Color(getHexaColor(s.Base0C)),
				lipgloss.Color(getHexaColor(s.Base0D)),
				lipgloss.Color(getHexaColor(s.Base0E)),
				lipgloss.Color(getHexaColor(s.Base0F)),
			},
		}
		schemes = append(schemes, newsscheme)
	}

	return schemes
}
