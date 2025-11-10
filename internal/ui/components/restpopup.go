// internal/ui/components/restpopup.go
package components

import (
	"fmt"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RestType represents the type of rest
type RestType string

const (
	ShortRestType RestType = "short"
	LongRestType  RestType = "long"
)

// RestPopup is a popup for resting (short/long rest)
type RestPopup struct {
	visible     bool
	restType    RestType
	character   *models.Character
	diceToSpend int    // For short rest: number of hit dice to spend
	confirmed   bool   // Whether the user confirmed the rest
	cancelled   bool   // Whether the user cancelled
	healing     int    // HP restored (set after rest is performed)
	roller      models.DiceRoller
}

// NewRestPopup creates a new rest popup
func NewRestPopup(char *models.Character, roller models.DiceRoller) *RestPopup {
	return &RestPopup{
		character: char,
		roller:    roller,
	}
}

// Show displays the rest popup
func (rp *RestPopup) Show(restType RestType) {
	rp.visible = true
	rp.restType = restType
	rp.diceToSpend = 0
	rp.confirmed = false
	rp.cancelled = false
	rp.healing = 0
}

// Hide hides the rest popup
func (rp *RestPopup) Hide() {
	rp.visible = false
	rp.confirmed = false
	rp.cancelled = false
}

// IsVisible returns whether the popup is visible
func (rp *RestPopup) IsVisible() bool {
	return rp.visible
}

// IsConfirmed returns whether the user confirmed the rest
func (rp *RestPopup) IsConfirmed() bool {
	return rp.confirmed
}

// IsCancelled returns whether the user cancelled
func (rp *RestPopup) IsCancelled() bool {
	return rp.cancelled
}

// GetHealing returns the HP restored
func (rp *RestPopup) GetHealing() int {
	return rp.healing
}

// GetRestType returns the type of rest
func (rp *RestPopup) GetRestType() RestType {
	return rp.restType
}

// Update handles input
func (rp *RestPopup) Update(msg tea.Msg) tea.Cmd {
	if !rp.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			rp.cancelled = true
			rp.visible = false
			return nil

		case "enter":
			// Confirm rest
			if rp.restType == ShortRestType {
				rp.healing = rp.character.PerformShortRest(rp.diceToSpend, rp.roller)
			} else {
				rp.character.PerformLongRest()
				rp.healing = 0 // Long rest fully restores HP, not tracked as "healing"
			}
			rp.confirmed = true
			rp.visible = false
			return nil

		case "up", "k":
			if rp.restType == ShortRestType && rp.diceToSpend < rp.character.GetTotalHitDice() {
				rp.diceToSpend++
			}

		case "down", "j":
			if rp.restType == ShortRestType && rp.diceToSpend > 0 {
				rp.diceToSpend--
			}

		case "left", "h":
			if rp.restType == ShortRestType && rp.diceToSpend > 0 {
				rp.diceToSpend--
			}

		case "right", "l":
			if rp.restType == ShortRestType && rp.diceToSpend < rp.character.GetTotalHitDice() {
				rp.diceToSpend++
			}
		}
	}

	return nil
}

// View renders the rest popup
func (rp *RestPopup) View() string {
	if !rp.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(60)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	valueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	var content strings.Builder

	if rp.restType == ShortRestType {
		content.WriteString(titleStyle.Render("SHORT REST") + "\n\n")

		content.WriteString(labelStyle.Render("Take a short rest (at least 1 hour) to:") + "\n")
		content.WriteString(dimStyle.Render("• Restore short rest features") + "\n")
		content.WriteString(dimStyle.Render("• Spend hit dice to restore HP") + "\n")
		content.WriteString(dimStyle.Render("• Restore 1 Psi Die (Psi Warrior)") + "\n")
		content.WriteString(dimStyle.Render("• Restore all Superiority Dice (Battle Master)") + "\n")
		content.WriteString(dimStyle.Render("• Restore 1 Warrior Die (Zealot)") + "\n\n")

		availableHD := rp.character.GetTotalHitDice()
		content.WriteString(labelStyle.Render("Hit Dice Available: ") +
			valueStyle.Render(fmt.Sprintf("%d/%d", availableHD, rp.character.HitDice.Max)) + "\n\n")

		// Show hit dice by class
		if len(rp.character.HitDiceByClass) > 0 {
			content.WriteString(dimStyle.Render("Your hit dice:") + "\n")
			for className, pool := range rp.character.HitDiceByClass {
				content.WriteString(dimStyle.Render(fmt.Sprintf("  %s: %dd%d (%d/%d available)",
					className, pool.Count, pool.DieSize, pool.Count, pool.MaxDice)) + "\n")
			}
			content.WriteString("\n")
		}

		content.WriteString(labelStyle.Render("Hit Dice to Spend: ") +
			valueStyle.Render(fmt.Sprintf("%d", rp.diceToSpend)) + "\n")
		content.WriteString(dimStyle.Render("(Use ↑/↓ or h/j/k/l to adjust)") + "\n\n")

		// Show estimated healing
		if rp.diceToSpend > 0 {
			conMod := rp.character.AbilityScores.GetModifier("Constitution")
			minHealing := rp.diceToSpend * (1 + conMod)
			if minHealing < rp.diceToSpend {
				minHealing = rp.diceToSpend // At least 1 HP per die
			}

			// Estimate average (assume average roll is half die size)
			avgHealing := 0
			diceSpent := 0
			for _, pool := range rp.character.HitDiceByClass {
				if diceSpent >= rp.diceToSpend {
					break
				}
				diceToSpend := rp.diceToSpend - diceSpent
				if diceToSpend > pool.Count {
					diceToSpend = pool.Count
				}
				avgRoll := (pool.DieSize + 1) / 2
				avgHealing += diceToSpend * (avgRoll + conMod)
				diceSpent += diceToSpend
			}

			content.WriteString(dimStyle.Render(fmt.Sprintf("Estimated HP restored: ~%d (min %d)", avgHealing, minHealing)) + "\n")
			content.WriteString(dimStyle.Render(fmt.Sprintf("Current HP: %d/%d", rp.character.CurrentHP, rp.character.MaxHP)) + "\n\n")
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		content.WriteString(successStyle.Render("Enter: Confirm") + "  " +
			errorStyle.Render("Esc: Cancel"))

	} else { // Long Rest
		content.WriteString(titleStyle.Render("LONG REST") + "\n\n")

		content.WriteString(labelStyle.Render("Take a long rest (at least 8 hours) to restore:") + "\n")
		content.WriteString(dimStyle.Render("• All HP restored") + "\n")
		content.WriteString(dimStyle.Render("• Half of hit dice restored (minimum 1)") + "\n")
		content.WriteString(dimStyle.Render("• All features restored") + "\n")
		content.WriteString(dimStyle.Render("• All spell slots restored") + "\n")
		content.WriteString(dimStyle.Render("• All Psi Dice restored (Psi Warrior)") + "\n")
		content.WriteString(dimStyle.Render("• All Superiority Dice restored (Battle Master)") + "\n")
		content.WriteString(dimStyle.Render("• All Warrior Dice restored (Zealot)") + "\n")
		if rp.character.IsBeastMaster() {
			content.WriteString(dimStyle.Render("• Change beast companion (Beast Master)") + "\n")
		}
		content.WriteString("\n")

		content.WriteString(labelStyle.Render("Current Status:") + "\n")
		content.WriteString(dimStyle.Render(fmt.Sprintf("  HP: %d/%d", rp.character.CurrentHP, rp.character.MaxHP)) + "\n")
		content.WriteString(dimStyle.Render(fmt.Sprintf("  Hit Dice: %d/%d", rp.character.HitDice.Current, rp.character.HitDice.Max)) + "\n")

		// Count features that will be restored
		featuresToRestore := 0
		for _, feature := range rp.character.Features.Features {
			if feature.RestType != models.None && feature.CurrentUses < feature.MaxUses {
				featuresToRestore++
			}
		}
		if featuresToRestore > 0 {
			content.WriteString(dimStyle.Render(fmt.Sprintf("  Features to restore: %d", featuresToRestore)) + "\n")
		}

		content.WriteString("\n")
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		content.WriteString(successStyle.Render("Enter: Confirm") + "  " +
			errorStyle.Render("Esc: Cancel"))
		if rp.character.IsBeastMaster() {
			content.WriteString("\n" + dimStyle.Render("c: Change companion"))
		}
	}

	return lipgloss.Place(
		80, 30,
		lipgloss.Center, lipgloss.Center,
		contentStyle.Render(content.String()),
	)
}
