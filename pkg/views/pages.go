package views

import (
	"slices"

	"github.com/a-h/templ"
)

type Page int

const (
	PageProfile Page = iota
	PageServer
	PageServerEdit
	PageContainerList
	PageContainerView
	PageContainerEdit

	totalPages // keep this at the end to count total pages.
)

type NavbarPageConfig struct {
	Name  string
	Link  templ.SafeURL
	Pages []Page
}

func (p NavbarPageConfig) IsActive(page Page) bool {
	return slices.Contains(p.Pages, page)
}

//nolint:gochecknoglobals
var (
	navbarPages = []NavbarPageConfig{
		{
			Name: "Home",
			Link: "/dashboard",
			Pages: []Page{
				PageProfile,
			},
		},
		{
			Name: "Server",
			Link: "/dashboard/server",
			Pages: []Page{
				PageServer,
				PageServerEdit,
			},
		},
		{
			Name: "Containers",
			Link: "/dashboard/containers",
			Pages: []Page{
				PageContainerList,
				PageContainerView,
				PageContainerEdit,
			},
		},
	}
)
