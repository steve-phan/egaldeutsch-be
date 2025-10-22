package models

import "testing"

func TestUserRoleIsValid(t *testing.T) {
	if !UserRoleAdmin.IsValid() {
		t.Fatalf("admin should be valid")
	}
	if UserRole("unknown").IsValid() {
		t.Fatalf("unknown should be invalid")
	}
}

func TestLanguageLevelIsValid(t *testing.T) {
	if !LevelA1.IsValid() {
		t.Fatalf("A1 should be valid")
	}
	if LanguageLevel("X9").IsValid() {
		t.Fatalf("X9 should be invalid")
	}
}
