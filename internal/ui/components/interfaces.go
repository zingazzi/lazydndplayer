// internal/ui/components/interfaces.go
package components

// Component defines the base interface that all UI components must implement
type Component interface {
	IsVisible() bool
}

// Selector defines the interface for selector components (lists that can be navigated)
type Selector interface {
	Component
	Next()
	Prev()
	Hide()
}

// Popup defines the interface for popup components (overlays that can be closed)
type Popup interface {
	Component
	Hide()
}

// DetailPopup defines the interface for detail popup components (popups with titles)
type DetailPopup interface {
	Popup
	GetTitle() string
}

// PageNavigationSelector extends Selector with page navigation capabilities
type PageNavigationSelector interface {
	Selector
	PageUp()
	PageDown()
}

// CategoryNavigationSelector extends Selector with category navigation capabilities
type CategoryNavigationSelector interface {
	Selector
	NextCategory()
	PrevCategory()
}
