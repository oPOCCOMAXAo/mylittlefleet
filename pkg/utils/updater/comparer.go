package updater

type Comparer struct {
	changed bool
}

func NewComparer() *Comparer {
	return &Comparer{}
}

func (cmp *Comparer) IsChanged() bool {
	return cmp.changed
}

func (cmp *Comparer) IsEqual() bool {
	return !cmp.changed
}

func (cmp *Comparer) SetChanged() {
	cmp.changed = true
}

func CompareValues[T comparable](cmp *Comparer, old T, newValue T) {
	if old == newValue {
		return
	}

	cmp.changed = true
}
