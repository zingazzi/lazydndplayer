// internal/ui/view/layout.go
package view

// PopupSize represents the size category for popups
type PopupSize int

const (
	PopupSizeSmall PopupSize = iota
	PopupSizeMedium
	PopupSizeLarge
	PopupSizeFull
)

// LayoutDimensions contains all calculated layout dimensions
type LayoutDimensions struct {
	// Status bar
	StatusBarHeight int

	// Top row
	MainPanelWidth  int
	CharStatsWidth  int
	TopRowHeight    int
	MainContentHeight int
	CharStatsInnerWidth  int
	CharStatsInnerHeight int

	// Bottom row
	BottomHeight      int
	ActionsWidth      int
	DiceWidth         int
	ActionsWidthRatio int
	DiceWidthRatio    int
	BottomInnerHeight int

	// Popup dimensions
	PopupSmallWidth   int
	PopupSmallHeight  int
	PopupMediumWidth  int
	PopupMediumHeight int
	PopupLargeWidth   int
	PopupLargeHeight  int
}

// LayoutCalculator calculates layout dimensions based on screen size
type LayoutCalculator struct{}

// NewLayoutCalculator creates a new layout calculator
func NewLayoutCalculator() *LayoutCalculator {
	return &LayoutCalculator{}
}

// CalculateLayout calculates all layout dimensions for the given screen size
func (lc *LayoutCalculator) CalculateLayout(width, height int) *LayoutDimensions {
	// Status bar height (4% of screen height, minimum 1)
	statusBarHeight := int(float64(height) * 0.04)
	if statusBarHeight < 1 {
		statusBarHeight = 1
	}

	// First row: Main panel (55%) + Character stats (43%)
	mainPanelWidth := int(float64(width) * 0.55)
	charStatsWidth := int(float64(width) * 0.43)

	// Calculate panel heights: top row 48%, bottom row 42%
	topRowHeight := int(float64(height) * 0.48)
	bottomHeight := int(float64(height) * 0.42)

	// Ensure minimum heights
	if topRowHeight < 10 {
		topRowHeight = 10
	}
	if bottomHeight < 8 {
		bottomHeight = 8
	}

	// Main content height accounts for tabs and spacing
	// Note: tabBarWidth and tabHeight are calculated in renderer, not here
	// We'll calculate this after getting tab height, but estimate here
	mainContentHeight := topRowHeight - 5 // border (2) + padding vertical (2) + spacing line (1)
	if mainContentHeight < 5 {
		mainContentHeight = 5
	}

	// Character stats inner dimensions
	charStatsInnerWidth := charStatsWidth - 8 // Account for border + padding
	charStatsInnerHeight := topRowHeight - 6  // Account for border (2) + vertical padding (4)

	// Bottom panels: Actions (50%) + Dice Roller (48%)
	actionsWidthRatio := int(float64(width) * 0.50)
	diceWidthRatio := int(float64(width) * 0.48)

	// Calculate inner dimensions: account for border + padding
	actionsWidth := actionsWidthRatio - 8
	diceWidth := diceWidthRatio - 8
	bottomInnerHeight := bottomHeight - 6 // Account for border (2) + vertical padding (4)

	// Popup dimensions
	popupSmallWidth, popupSmallHeight := lc.CalculatePopupDimensions(width, height, PopupSizeSmall)
	popupMediumWidth, popupMediumHeight := lc.CalculatePopupDimensions(width, height, PopupSizeMedium)
	popupLargeWidth, popupLargeHeight := lc.CalculatePopupDimensions(width, height, PopupSizeLarge)

	return &LayoutDimensions{
		StatusBarHeight:     statusBarHeight,
		MainPanelWidth:      mainPanelWidth,
		CharStatsWidth:      charStatsWidth,
		TopRowHeight:        topRowHeight,
		MainContentHeight:   mainContentHeight,
		CharStatsInnerWidth: charStatsInnerWidth,
		CharStatsInnerHeight: charStatsInnerHeight,
		BottomHeight:        bottomHeight,
		ActionsWidth:        actionsWidth,
		DiceWidth:           diceWidth,
		ActionsWidthRatio:   actionsWidthRatio,
		DiceWidthRatio:      diceWidthRatio,
		BottomInnerHeight:   bottomInnerHeight,
		PopupSmallWidth:     popupSmallWidth,
		PopupSmallHeight:    popupSmallHeight,
		PopupMediumWidth:    popupMediumWidth,
		PopupMediumHeight:   popupMediumHeight,
		PopupLargeWidth:     popupLargeWidth,
		PopupLargeHeight:    popupLargeHeight,
	}
}

// CalculatePopupDimensions calculates popup dimensions based on size category
func (lc *LayoutCalculator) CalculatePopupDimensions(width, height int, size PopupSize) (int, int) {
	const (
		PopupSmallWidthPercent  = 0.50
		PopupSmallHeightPercent = 0.60
		PopupSmallMinWidth      = 60
		PopupSmallMinHeight     = 20

		PopupMediumWidthPercent  = 0.75
		PopupMediumHeightPercent = 0.80
		PopupMediumMinWidth      = 80
		PopupMediumMinHeight     = 25

		PopupLargeWidthPercent  = 0.85
		PopupLargeHeightPercent = 0.85
		PopupLargeMinWidth      = 90
		PopupLargeMinHeight     = 30
	)

	var popupWidth, popupHeight int
	var minWidth, minHeight int

	switch size {
	case PopupSizeSmall:
		popupWidth = int(float64(width) * PopupSmallWidthPercent)
		popupHeight = int(float64(height) * PopupSmallHeightPercent)
		minWidth = PopupSmallMinWidth
		minHeight = PopupSmallMinHeight
	case PopupSizeMedium:
		popupWidth = int(float64(width) * PopupMediumWidthPercent)
		popupHeight = int(float64(height) * PopupMediumHeightPercent)
		minWidth = PopupMediumMinWidth
		minHeight = PopupMediumMinHeight
	case PopupSizeLarge:
		popupWidth = int(float64(width) * PopupLargeWidthPercent)
		popupHeight = int(float64(height) * PopupLargeHeightPercent)
		minWidth = PopupLargeMinWidth
		minHeight = PopupLargeMinHeight
	case PopupSizeFull:
		popupWidth = width
		popupHeight = height
		minWidth = 0
		minHeight = 0
	}

	// Apply minimums
	if popupWidth < minWidth {
		popupWidth = minWidth
	}
	if popupHeight < minHeight {
		popupHeight = minHeight
	}

	return popupWidth, popupHeight
}
