// tests/ui/services/class_service_test.go
package services_test

import (
	"strings"
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/services"
)

func TestNewClassService(t *testing.T) {
	service := services.NewClassService()
	if service == nil {
		t.Fatal("NewClassService() returned nil")
	}
}

func TestClassServiceRequiresSkillChoice(t *testing.T) {
	service := services.NewClassService()

	// Test class that requires skill choices (Rogue)
	requires := service.RequiresSkillChoice("Rogue")
	if !requires {
		t.Error("Expected Rogue to require skill choices")
	}

	// Test class that doesn't exist
	requires = service.RequiresSkillChoice("NonExistentClass")
	if requires {
		t.Error("Expected non-existent class to not require skill choices")
	}
}

func TestClassServiceGetSkillChoices(t *testing.T) {
	service := services.NewClassService()

	// Test Rogue skill choices
	from, choose := service.GetSkillChoices("Rogue")
	if choose == 0 {
		t.Error("Expected Rogue to have skill choices")
	}
	if len(from) == 0 {
		t.Error("Expected skill list to be non-empty")
	}
	if choose > len(from) {
		t.Errorf("Expected choose count (%d) to be <= available skills (%d)", choose, len(from))
	}
}

func TestClassServiceApplyClass(t *testing.T) {
	service := services.NewClassService()

	char := models.NewCharacter()
	char.Name = "Test Character"

	msg, err := service.ApplyClass(char, "Fighter")
	if err != nil {
		t.Fatalf("Expected no error applying class, got %v", err)
	}

	if !strings.Contains(msg, "Fighter") {
		t.Errorf("Expected message to contain class name, got %s", msg)
	}

	// Check that character has the class
	if !char.HasClass("Fighter") {
		t.Error("Expected character to have Fighter class")
	}
}

func TestClassServiceApplyClassInvalid(t *testing.T) {
	service := services.NewClassService()

	char := models.NewCharacter()

	_, err := service.ApplyClass(char, "InvalidClass")
	if err == nil {
		t.Error("Expected error when applying invalid class")
	}
}

func TestClassServiceApplySubclass(t *testing.T) {
	service := services.NewClassService()

	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 1},
	}

	msg, err := service.ApplySubclass(char, "Champion")
	if err != nil {
		t.Fatalf("Expected no error applying subclass, got %v", err)
	}

	if !strings.Contains(msg, "Champion") {
		t.Errorf("Expected message to contain subclass name, got %s", msg)
	}

	if char.Classes[0].Subclass != "Champion" {
		t.Errorf("Expected subclass to be 'Champion', got %s", char.Classes[0].Subclass)
	}
}

func TestClassServiceApplySubclassNoClasses(t *testing.T) {
	service := services.NewClassService()

	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{}

	_, err := service.ApplySubclass(char, "Champion")
	if err == nil {
		t.Error("Expected error when character has no classes")
	}
}

func TestClassServiceGetSubclasses(t *testing.T) {
	service := services.NewClassService()

	subclasses, err := service.GetSubclasses("Fighter", 1)
	if err != nil {
		t.Fatalf("Expected no error getting subclasses, got %v", err)
	}

	if len(subclasses) == 0 {
		t.Error("Expected Fighter to have subclasses at level 1")
	}

	// Test invalid class
	_, err = service.GetSubclasses("InvalidClass", 1)
	if err == nil {
		t.Error("Expected error when getting subclasses for invalid class")
	}
}
