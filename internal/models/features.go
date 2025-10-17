// internal/models/features.go
package models

// RestType indicates when a feature recharges
type RestType string

const (
	ShortRest RestType = "Short Rest"
	LongRest  RestType = "Long Rest"
	Daily     RestType = "Daily"
	None      RestType = "None" // Passive features
)

// Feature represents a limited-use ability (class features, racial abilities, etc.)
type Feature struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MaxUses     int      `json:"max_uses"`      // Maximum uses per rest
	CurrentUses int      `json:"current_uses"`  // Current available uses
	RestType    RestType `json:"rest_type"`     // When it recharges
	Source      string   `json:"source"`        // e.g., "Class: Barbarian", "Species: Dragonborn"
}

// FeatureList manages character features
type FeatureList struct {
	Features []Feature `json:"features"`
}

// NewFeatureList creates an empty feature list
func NewFeatureList() *FeatureList {
	return &FeatureList{
		Features: []Feature{},
	}
}

// AddFeature adds a new feature
func (fl *FeatureList) AddFeature(feature Feature) {
	// Set current uses to max if not specified
	if feature.CurrentUses == 0 && feature.MaxUses > 0 {
		feature.CurrentUses = feature.MaxUses
	}
	fl.Features = append(fl.Features, feature)
}

// UseFeature decrements the usage count for a feature
func (fl *FeatureList) UseFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}

	feature := &fl.Features[index]
	if feature.CurrentUses > 0 {
		feature.CurrentUses--
		return true
	}
	return false
}

// RestoreFeature restores one use of a feature
func (fl *FeatureList) RestoreFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}

	feature := &fl.Features[index]
	if feature.CurrentUses < feature.MaxUses {
		feature.CurrentUses++
		return true
	}
	return false
}

// ShortRestRecover recovers features that recharge on short rest
func (fl *FeatureList) ShortRestRecover() {
	for i := range fl.Features {
		if fl.Features[i].RestType == ShortRest || fl.Features[i].RestType == Daily {
			fl.Features[i].CurrentUses = fl.Features[i].MaxUses
		}
	}
}

// LongRestRecover recovers all rechargeable features
func (fl *FeatureList) LongRestRecover() {
	for i := range fl.Features {
		if fl.Features[i].RestType != None {
			fl.Features[i].CurrentUses = fl.Features[i].MaxUses
		}
	}
}

// RemoveFeature removes a feature by index
func (fl *FeatureList) RemoveFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}
	fl.Features = append(fl.Features[:index], fl.Features[index+1:]...)
	return true
}
