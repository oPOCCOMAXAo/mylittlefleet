package views

type Page int

const (
	PageProfile Page = iota
	PageServer
	PageServerEdit
)

func (p Page) OneOf(pages ...Page) bool {
	for _, page := range pages {
		if p == page {
			return true
		}
	}

	return false
}
