package discord

import (
	"github.com/heshoots/dmux"
	"testing"
)

func TestGetRole(t *testing.T) {
	var mockRoles []dmux.Role
	var mockRole = dmux.MockRole{RoleID: "1234", RoleName: "Test Role"}
	mockRoles = append(mockRoles, mockRole)
	mockSession := dmux.MockSession{SessionRoles: mockRoles}
	role, err := GetRole(mockSession, "", "Test Role")
	if err != nil {
		t.Errorf("GetRole threw an error, " + err.Error())
	}
	if role.Name() != "Test Role" {
		t.Errorf("Incorrect Role returned")
	}
}
