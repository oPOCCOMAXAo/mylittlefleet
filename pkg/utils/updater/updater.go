package updater

type Updater struct {
	changed bool
}

func New() *Updater {
	return &Updater{}
}

func (upd *Updater) IsChanged() bool {
	return upd.changed
}

func (upd *Updater) AppendChanged(newChanged bool) {
	upd.changed = upd.changed || newChanged
}

func (upd *Updater) SetChanged() {
	upd.changed = true
}

func SetValue[T comparable](upd *Updater, dst *T, newValue T) {
	if *dst == newValue {
		return
	}

	*dst = newValue

	upd.changed = true
}
