package scheme

import (
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/o1egl/fwencoder"
)

type Scheme struct {
	Name   string
	Source string
	Colors []lipgloss.Color
	Author string
}

func getHexaColor(color string) string {
	r, _ := regexp.Compile(`#?([0-9A-F]{6})`)
	if match := r.MatchString(strings.ToUpper(color)); match == true {
		groups := r.FindStringSubmatch(strings.ToUpper(color))
		return "#" + groups[1]
	}

	return ""
}

type Importer interface {
	GetRessouresList() []string
	ImportRessources(ressources []string, schemes Schemes) Schemes
}

func Import(schemes Schemes, i Importer) Schemes {
	res := i.GetRessouresList()
	schemes = i.ImportRessources(res, schemes)

	return schemes
}

type Schemes []Scheme

func (schemes *Schemes) Read(filename string) {
	if _, err := os.Stat(filename); err == nil {

		f, _ := os.Open(filename)
		defer f.Close()

		// var datas Bands
		err := fwencoder.UnmarshalReader(f, schemes)

		if err != nil {
			panic(err)
		}
	}
}

func (schemes *Schemes) Write(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	_ = fwencoder.MarshalWriter(f, schemes)
}

// Sorting
func (schemes Schemes) Len() int      { return len(schemes) }
func (schemes Schemes) Swap(i, j int) { schemes[i], schemes[j] = schemes[j], schemes[i] }
func (schemes Schemes) Less(i, j int) bool {
	return schemes[i].Name < schemes[j].Name
}
