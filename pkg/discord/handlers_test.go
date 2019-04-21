package discord

import (
	"github.com/quorauk/cholibot/pkg/models"
	"github.com/quorauk/dmux"
	"github.com/jinzhu/configor"
	"testing"
)

func TestGetRole(t *testing.T) {
	var mockRoles []dmux.Role
	var mockRole = dmux.MockRole{RoleID: "1234", RoleName: "Test Role"}
	mockRoles = append(mockRoles, mockRole)
	mockGuild := dmux.MockGuild{SessionRoles: mockRoles}
	mockSession := dmux.MockSession{Guild: mockGuild}
	role, err := GetRole(mockSession, mockGuild, "Test Role")
	if err != nil {
		t.Errorf("GetRole threw an error, " + err.Error())
	}
	if role.Name() != "Test Role" {
		t.Errorf("Incorrect Role returned")
	}
}
