package discord

import (
	"testing"
)

func TestGetRole(t *testing.T) {
	var mockRoles []Role
	var mockRole = MockRole{id: "1234", name: "Test Role"}
	mockRoles = append(mockRoles, mockRole)
	mockSession := MockSession{roles: mockRoles}
	role, err := GetRole(mockSession, "", "Test Role")
	if err != nil {
		t.Errorf("GetRole threw an error, " + err.Error())
	}
	if role.Name() != "Test Role" {
		t.Errorf("Incorrect Role returned")
	}
}
