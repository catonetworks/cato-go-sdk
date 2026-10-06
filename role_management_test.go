package cato_go_sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

func TestRoleManagementOperations(t *testing.T) {
	role := `{"id":"42","name":"Test role","description":"","predefined":false,"isUsedOnExternalAccess":true,"accountType":"REGULAR","permission":[{"resource":"Sites","action":"VIEW"}]}`
	for _, leaf := range []string{"createRole", "updateRole", "deleteRole", "role", "roleList", "permissionCatalog"} {
		t.Run(leaf, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request struct {
					Query     string                     `json:"query"`
					Variables map[string]json.RawMessage `json:"variables"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
					return
				}
				if string(request.Variables["accountId"]) != `"123"` {
					t.Errorf("account scope: %s", request.Variables["accountId"])
				}
				if !strings.Contains(request.Query, "rbac(accountId: $accountId)") || !strings.Contains(request.Query, "roleManagement") {
					t.Errorf("wrong namespace: %s", request.Query)
				}
				payload := role
				switch leaf {
				case "createRole", "updateRole":
					var input cato_models.RoleManagementUpdateRoleInput
					if err := json.Unmarshal(request.Variables["input"], &input); err != nil {
						t.Error(err)
					}
					if input.Name != "Test role" || input.Description == nil || *input.Description != "" || len(input.Permission) != 1 || input.Permission[0].Resource != "Sites" || input.Permission[0].Action != cato_models.RoleManagementGrantActionView {
						t.Errorf("unexpected mutation input: %+v", input)
					}
					if leaf == "updateRole" && input.ID != "42" {
						t.Error("missing update ID")
					}
					payload = `{"role":` + role + `}`
				case "deleteRole":
					if string(request.Variables["input"]) != `{"id":"42"}` {
						t.Errorf("delete input: %s", request.Variables["input"])
					}
					payload = `{"id":"42"}`
				case "role":
					if string(request.Variables["id"]) != `"42"` {
						t.Error("missing lookup ID")
					}
				case "roleList":
					var input cato_models.RoleManagementRoleListInput
					if err := json.Unmarshal(request.Variables["input"], &input); err != nil {
						t.Error(err)
					}
					if input.Paging == nil || input.Paging.From != 100 || input.Paging.Limit != 100 || input.Sort == nil || input.Sort.Name.Direction != cato_models.SortOrderAsc || input.Filter == nil || input.Filter.Name.Eq == nil || *input.Filter.Name.Eq != "Test role" {
						t.Errorf("list input: %+v", input)
					}
					payload = `{"items":[` + role + `],"paging":{"total":101}}`
				case "permissionCatalog":
					payload = `{"resource":[{"resource":"Sites","supportedAction":["VIEW","EDIT"]}]}`
				}
				w.Header().Set("Content-Type", "application/json")
				marker := ""
				if leaf == "role" {
					marker = `"__typename":"RoleManagementQueries",`
				}
				_, _ = w.Write([]byte(`{"data":{"rbac":{"roleManagement":{` + marker + `"` + leaf + `":` + payload + `}}}}`))
			}))
			defer srv.Close()
			client := NewClient(srv.Client(), srv.URL, nil)
			ctx := context.Background()
			description := ""
			name := "Test role"
			permissions := []*cato_models.RoleManagementPermissionInput{{Resource: "Sites", Action: cato_models.RoleManagementGrantActionView}}
			switch leaf {
			case "createRole":
				v, e := client.RbacRoleManagementCreateRole(ctx, "123", cato_models.RoleManagementCreateRoleInput{Name: name, Description: &description, Permission: permissions})
				if e != nil || v.GetRbac().GetRoleManagement().GetCreateRole().GetRole().GetID() != "42" {
					t.Fatalf("response: %v %v", v, e)
				}
			case "updateRole":
				v, e := client.RbacRoleManagementUpdateRole(ctx, "123", cato_models.RoleManagementUpdateRoleInput{ID: "42", Name: name, Description: &description, Permission: permissions})
				if e != nil || v.GetRbac().GetRoleManagement().GetUpdateRole().GetRole().GetID() != "42" {
					t.Fatalf("response: %v %v", v, e)
				}
			case "deleteRole":
				v, e := client.RbacRoleManagementDeleteRole(ctx, "123", cato_models.RoleManagementDeleteRoleInput{ID: "42"})
				if e != nil || v.GetRbac().GetRoleManagement().GetDeleteRole().GetID() != "42" {
					t.Fatalf("response: %v %v", v, e)
				}
			case "role":
				v, e := client.RbacRoleManagementRole(ctx, "123", "42")
				if e != nil || v.GetRbac().GetRoleManagement().GetRole().GetPermission()[0].GetResource() != "Sites" {
					t.Fatalf("response: %v %v", v, e)
				}
			case "roleList":
				v, e := client.RbacRoleManagementRoleList(ctx, "123", cato_models.RoleManagementRoleListInput{Paging: &cato_models.PagingInput{From: 100, Limit: 100}, Sort: &cato_models.RoleManagementRoleSortInput{Name: &cato_models.SortOrderInput{Direction: cato_models.SortOrderAsc}}, Filter: &cato_models.RoleManagementRoleFilterInput{Name: &cato_models.StringFilterInput{Eq: &name}}})
				if e != nil || v.GetRbac().GetRoleManagement().GetRoleList().GetPaging().GetTotal() != 101 {
					t.Fatalf("response: %v %v", v, e)
				}
			case "permissionCatalog":
				v, e := client.RbacRoleManagementPermissionCatalog(ctx, "123")
				if e != nil || len(v.GetRbac().GetRoleManagement().GetPermissionCatalog().GetResource()[0].GetSupportedAction()) != 2 {
					t.Fatalf("response: %v %v", v, e)
				}
			}
		})
	}
}

func TestRoleManagementNullablePayloadAndGraphQLError(t *testing.T) {
	for _, body := range []string{`{"data":{"rbac":{"roleManagement":{"__typename":"RoleManagementQueries","role":null}}}}`, `{"errors":[{"message":"Forbidden"}],"data":null}`} {
		t.Run(body, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(body))
			}))
			defer srv.Close()
			value, err := NewClient(srv.Client(), srv.URL, nil).RbacRoleManagementRole(context.Background(), "123", "42")
			if strings.Contains(body, "errors") {
				if err == nil {
					t.Fatal("HTTP 200 GraphQL error must fail")
				}
				return
			}
			if err != nil || value.GetRbac().GetRoleManagement().GetRole() != nil {
				t.Fatalf("null role: %v %v", value, err)
			}
		})
	}
}
