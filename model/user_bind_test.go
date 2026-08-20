package model

import "testing"

func TestUpdateUserBindColumnRejectsUnknownColumn(t *testing.T) {
	if err := UpdateUserBindColumn(1, "role", "user"); err == nil {
		t.Fatal("expected unknown user bind column to be rejected")
	}
}
