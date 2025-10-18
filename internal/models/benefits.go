// internal/models/benefits.go
package models

// BenefitType represents the type of benefit granted
type BenefitType string

const (
	BenefitAbilityScore BenefitType = "ability_score"
	BenefitProficiency  BenefitType = "proficiency"
	BenefitLanguage     BenefitType = "language"
	BenefitResistance   BenefitType = "resistance"
	BenefitSpeed        BenefitType = "speed"
	BenefitHP           BenefitType = "hp"
	BenefitSpell        BenefitType = "spell"
	BenefitSkill        BenefitType = "skill"
)

// BenefitSource represents where a benefit came from
type BenefitSource struct {
	Type string `json:"type"` // "feat", "species", "class", "background"
	Name string `json:"name"` // Name of the feat/species/etc
}

// GrantedBenefit tracks a single benefit and its source
type GrantedBenefit struct {
	Source      BenefitSource `json:"source"`
	Type        BenefitType   `json:"benefit_type"`
	Target      string        `json:"target"`      // e.g., "Strength", "Perception", "Elvish"
	Value       int           `json:"value"`       // Numeric value (e.g., +1, +10)
	Description string        `json:"description"` // Human-readable description
}

// BenefitTracker manages all granted benefits
type BenefitTracker struct {
	Benefits []GrantedBenefit `json:"benefits"`
}

// NewBenefitTracker creates a new benefit tracker
func NewBenefitTracker() *BenefitTracker {
	return &BenefitTracker{
		Benefits: []GrantedBenefit{},
	}
}

// AddBenefit records a new granted benefit
func (bt *BenefitTracker) AddBenefit(benefit GrantedBenefit) {
	bt.Benefits = append(bt.Benefits, benefit)
}

// GetBenefitsBySource returns all benefits from a specific source
func (bt *BenefitTracker) GetBenefitsBySource(sourceType, sourceName string) []GrantedBenefit {
	var benefits []GrantedBenefit
	for _, b := range bt.Benefits {
		if b.Source.Type == sourceType && b.Source.Name == sourceName {
			benefits = append(benefits, b)
		}
	}
	return benefits
}

// RemoveBenefitsBySource removes all benefits from a specific source
func (bt *BenefitTracker) RemoveBenefitsBySource(sourceType, sourceName string) []GrantedBenefit {
	removed := []GrantedBenefit{}
	filtered := []GrantedBenefit{}

	for _, b := range bt.Benefits {
		if b.Source.Type == sourceType && b.Source.Name == sourceName {
			removed = append(removed, b)
		} else {
			filtered = append(filtered, b)
		}
	}

	bt.Benefits = filtered
	return removed
}

// GetBenefitsByType returns all benefits of a specific type
func (bt *BenefitTracker) GetBenefitsByType(benefitType BenefitType) []GrantedBenefit {
	var benefits []GrantedBenefit
	for _, b := range bt.Benefits {
		if b.Type == benefitType {
			benefits = append(benefits, b)
		}
	}
	return benefits
}
