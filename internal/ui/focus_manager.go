package ui

type Focusable interface {
	Focus()
	Blur()
}

type FocusManager struct {
	currentIndex int
	//maxIndex     int
	components []Focusable
}

func (f *FocusManager) Next() {
	f.components[f.currentIndex].Blur()
	f.currentIndex = (f.currentIndex + 1) % len(f.components)
	f.components[f.currentIndex].Focus()
}

func (f *FocusManager) Prev() {
	f.components[f.currentIndex].Blur()
	f.currentIndex = (f.currentIndex - 1) % len(f.components)
	f.components[f.currentIndex].Focus()
}
