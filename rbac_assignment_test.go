package cato_go_sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/catonetworks/cato-go-sdk/models"
)

func TestAssignmentGeneratedRequests(t *testing.T) {
	for _, operation := range []string{"create", "update", "list", "remove", "error"} {
		t.Run(operation, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("x-api-key") != "fixture-key" || r.Header.Get("x-account-id") != "123" {
					t.Error("missing authentication/account headers")
				}
				var body struct {
					Variables struct {
						Account string                     `json:"accountID"`
						Input   map[string]json.RawMessage `json:"input"`
					} `json:"variables"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Variables.Account != "456" {
					t.Errorf("account override lost: %q", body.Variables.Account)
				}
				input := body.Variables.Input
				if _, ok := input["allowViewingUnobfuscatedData"]; ok {
					t.Error("unset preference must be omitted")
				}
				w.Header().Set("Content-Type", "application/json")
				switch operation {
				case "create":
					if string(input["groups"]) != "[]" || string(input["allowedEntities"]) != "[]" {
						t.Error("empty arrays lost")
					}
					_, _ = w.Write([]byte(`{"data":{"rbac":{"addRbacAssignments":{"items":[{"id":"a"}]}}}}`))
				case "update":
					if _, ok := input["users"]; ok {
						t.Error("create-only users sent on update")
					}
					if string(input["rbacAssignmentId"]) != `"a"` || string(input["roles"]) != "[]" {
						t.Error("update shape incorrect")
					}
					_, _ = w.Write([]byte(`{"data":{"rbac":{"updateRbacAssignment":{"assignment":{"id":"a"}}}}}`))
				case "list":
					var paging struct {
						From  int64 `json:"from"`
						Limit int64 `json:"limit"`
					}
					_ = json.Unmarshal(input["paging"], &paging)
					if paging.From != 2 || paging.Limit != 100 {
						t.Error("paging changed")
					}
					_, _ = w.Write([]byte(`{"data":{"rbac":{"rbacAssignmentList":{"items":[],"paging":{"total":0}}}}}`))
				case "remove":
					if string(input["rbacAssignmentId"]) != `"a"` {
						t.Error("missing removal ID")
					}
					_, _ = w.Write([]byte(`{"data":{"rbac":{"removeRbacAssignment":{"rbacAssignmentId":"a"}}}}`))
				default:
					_, _ = w.Write([]byte(`{"errors":[{"message":"access denied"}]}`))
				}
			}))
			defer server.Close()
			c, err := New(server.URL, "fixture-key", "123", server.Client(), nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			switch operation {
			case "create":
				result, err := c.RbacAddRbacAssignments(ctx, "456", m.AddRbacAssignmentsInput{Users: []*m.UserRefInput{{By: m.ObjectRefByID, Input: "u"}}, Groups: []*m.UsersGroupRefInput{}, Roles: []*m.RBACRoleRefInput{}, AllowedEntities: []*m.AuthorizedEntityRefInput{}, AllowAllSites: m.RBACActionNone, AllowAllUsers: m.RBACActionNone})
				if err != nil || result.GetRbac().GetAddRbacAssignments().GetItems()[0].GetID() != "a" {
					t.Fatalf("create: %v", err)
				}
			case "update":
				result, err := c.RbacUpdateRbacAssignment(ctx, "456", m.UpdateRbacAssignmentInput{RbacAssignmentID: "a", Roles: []*m.RBACRoleRefInput{}, AllowedEntities: []*m.AuthorizedEntityRefInput{}, AllowAllSites: m.RBACActionNone, AllowAllUsers: m.RBACActionNone})
				if err != nil || result.GetRbac().GetUpdateRbacAssignment().GetAssignment().GetID() != "a" {
					t.Fatalf("update: %v", err)
				}
			case "list":
				result, err := c.RbacRbacAssignmentList(ctx, "456", m.RBACAssignmentListInput{Paging: &m.PagingInput{From: 2, Limit: 100}, Sort: &m.RBACAssignmentSortInput{}})
				if err != nil || result.GetRbac().GetRbacAssignmentList().GetPaging().GetTotal() != 0 {
					t.Fatalf("list: %v", err)
				}
			case "remove":
				result, err := c.RbacRemoveRbacAssignment(ctx, "456", m.RemoveRbacAssignmentInput{RbacAssignmentID: "a"})
				if err != nil || result.GetRbac().GetRemoveRbacAssignment().GetRbacAssignmentID() != "a" {
					t.Fatalf("remove: %v", err)
				}
			default:
				if _, err := c.RbacRemoveRbacAssignment(ctx, "456", m.RemoveRbacAssignmentInput{RbacAssignmentID: "a"}); err == nil {
					t.Fatal("GraphQL errors must propagate")
				}
			}
		})
	}
}
