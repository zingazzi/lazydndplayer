// internal/ui/view/layout_test.go
package view

import (
	"testing"
)

func TestNewLayoutCalculator(t *testing.T) {
	calc := NewLayoutCalculator()

	if calc == nil {
		t.Fatal("NewLayoutCalculator() returned nil")
	}
}

func TestCalculateLayout(t *testing.T) {
	calc := NewLayoutCalculator()

	tests := []struct {
		name           string
		width          int
		height         int
		wantValid      bool
		checkMainWidth bool
		minMainWidth   int
	}{
		{
			name:           "Standard size",
			width:          120,
			height:         40,
			wantValid:      true,
			checkMainWidth: true,
			minMainWidth:   50,
		},
		{
			name:           "Minimum size",
			width:          80,
			height:         24,
			wantValid:      true,
			checkMainWidth: true,
			minMainWidth:   30,
		},
		{
			name:      "Very small size",
			width:     50,
			height:    20,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := calc.CalculateLayout(tt.width, tt.height)

		if layout == nil {
			t.Fatal("CalculateLayout() returned nil")
		}

		if tt.checkMainWidth && layout.MainPanelWidth < tt.minMainWidth {
			t.Errorf("Expected MainPanelWidth >= %d, got %d", tt.minMainWidth, layout.MainPanelWidth)
		}

		// Check that all dimensions are non-negative
		if layout.MainPanelWidth < 0 || layout.CharStatsWidth < 0 || layout.ActionsWidth < 0 {
			t.Error("Expected all dimensions to be non-negative")
		}
		})
	}
}

func TestCalculateLayoutDimensions(t *testing.T) {
	calc := NewLayoutCalculator()

	layout := calc.CalculateLayout(120, 40)

	// Verify layout dimensions are reasonable
	if layout.TopRowHeight <= 0 {
		t.Error("TopRowHeight should be positive")
	}

	if layout.BottomHeight <= 0 {
		t.Error("BottomHeight should be positive")
	}

	// Note: We can't directly check width/height since LayoutDimensions doesn't store them
	// But we can verify that calculated dimensions are reasonable
	if layout.MainPanelWidth+layout.CharStatsWidth < 0 {
		t.Error("Main panel and character stats panel widths should be valid")
	}

	if layout.ActionsWidth+layout.DiceWidth < 0 {
		t.Error("Actions and dice panel widths should be valid")
	}
}

func TestLayoutPopupDimensions(t *testing.T) {
	calc := NewLayoutCalculator()

	layout := calc.CalculateLayout(120, 40)

	// Check popup dimensions are reasonable
	if layout.PopupSmallWidth <= 0 || layout.PopupSmallHeight <= 0 {
		t.Error("Popup small dimensions should be positive")
	}

	if layout.PopupMediumWidth <= 0 || layout.PopupMediumHeight <= 0 {
		t.Error("Popup medium dimensions should be positive")
	}

	if layout.PopupLargeWidth <= 0 || layout.PopupLargeHeight <= 0 {
		t.Error("Popup large dimensions should be positive")
	}

	// Check that popup sizes are ordered correctly
	if layout.PopupSmallWidth >= layout.PopupMediumWidth {
		t.Error("Small popup should be smaller than medium popup")
	}

	if layout.PopupMediumWidth >= layout.PopupLargeWidth {
		t.Error("Medium popup should be smaller than large popup")
	}
}
