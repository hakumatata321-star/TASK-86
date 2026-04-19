package services_test

import (
	"database/sql"
	"testing"

	"w2t86/internal/config"
	"w2t86/internal/repository"
	"w2t86/internal/services"
	"w2t86/internal/testutil"
)

func newAdminSvc(t *testing.T) (*services.AdminService, *sql.DB) {
	t.Helper()
	db := testutil.NewTestDB(t)
	adminRepo := repository.NewAdminRepository(db)
	userRepo := repository.NewUserRepository(db)
	matRepo := repository.NewMaterialRepository(db)
	return services.NewAdminService(adminRepo, userRepo, matRepo), db
}

func newAuthSvcForAdmin(t *testing.T, db *sql.DB) *services.AuthService {
	t.Helper()
	cfg := &config.Config{SessionSecret: "test-secret"}
	return services.NewAuthService(
		repository.NewUserRepository(db),
		repository.NewSessionRepository(db),
		cfg,
	)
}

func TestAdminService_ListUsers_Empty(t *testing.T) {
	svc, _ := newAdminSvc(t)

	// The schema migration always seeds one admin user, so filter by
	// "student" role to verify no non-admin users exist in a fresh DB.
	users, err := svc.ListUsers("student", 50, 0)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 students in fresh DB, got %d", len(users))
	}
}

func TestAdminService_ListUsers_WithUsers(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	_, _ = auth.Register("au1", "au1@x.com", "TestPassword123!", "student")
	_, _ = auth.Register("au2", "au2@x.com", "TestPassword123!", "instructor")

	users, err := svc.ListUsers("", 50, 0)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) < 2 {
		t.Errorf("expected ≥2 users, got %d", len(users))
	}
}

func TestAdminService_ListUsers_RoleFilter(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	_, _ = auth.Register("rfstud", "rfs@x.com", "TestPassword123!", "student")
	_, _ = auth.Register("rfinstr", "rfi@x.com", "TestPassword123!", "instructor")

	students, err := svc.ListUsers("student", 50, 0)
	if err != nil {
		t.Fatalf("ListUsers(student): %v", err)
	}
	for _, u := range students {
		if u.Role != "student" {
			t.Errorf("expected role=student, got %q", u.Role)
		}
	}
}

func TestAdminService_ListUsers_Pagination(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	for i := 0; i < 5; i++ {
		u := string(rune('a' + i))
		_, _ = auth.Register("pgau"+u, "pgau"+u+"@x.com", "TestPassword123!", "student")
	}

	page1, err := svc.ListUsers("", 2, 0)
	if err != nil {
		t.Fatalf("ListUsers page1: %v", err)
	}
	if len(page1) != 2 {
		t.Errorf("expected 2 users on page 1, got %d", len(page1))
	}
}

func TestAdminService_MergeUsers_Success(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	primary, err := auth.Register("mrgprim", "mp@x.com", "TestPassword123!", "student")
	if err != nil {
		t.Fatalf("register primary: %v", err)
	}
	dup, err := auth.Register("mrgdup", "md@x.com", "TestPassword123!", "student")
	if err != nil {
		t.Fatalf("register duplicate: %v", err)
	}

	if err := svc.MergeUsers(primary.ID, dup.ID, primary.ID); err != nil {
		t.Fatalf("MergeUsers: %v", err)
	}

	// Primary user must still exist.
	var primaryCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE id=?`, primary.ID).Scan(&primaryCount); err != nil {
		t.Fatalf("query primary: %v", err)
	}
	if primaryCount != 1 {
		t.Errorf("primary user should still exist after merge, count=%d", primaryCount)
	}
}

func TestAdminService_SetAndGetCustomField(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	user, err := auth.Register("cftest", "cf@x.com", "TestPassword123!", "student")
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := svc.SetCustomField("user", user.ID, "department", "Engineering",
		false, nil, user.ID, "test"); err != nil {
		t.Fatalf("SetCustomField: %v", err)
	}

	fields, err := svc.GetCustomFields("user", user.ID, nil)
	if err != nil {
		t.Fatalf("GetCustomFields: %v", err)
	}

	found := false
	for _, f := range fields {
		if f.FieldName == "department" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'department' field; got %+v", fields)
	}
}

func TestAdminService_DeleteCustomField(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	user, err := auth.Register("delftest", "delf@x.com", "TestPassword123!", "student")
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	_ = svc.SetCustomField("user", user.ID, "badge", "gold", false, nil, user.ID, "")

	if err := svc.DeleteCustomField("user", user.ID, "badge", user.ID, "test"); err != nil {
		t.Fatalf("DeleteCustomField: %v", err)
	}

	fields, _ := svc.GetCustomFields("user", user.ID, nil)
	for _, f := range fields {
		if f.FieldName == "badge" {
			t.Error("expected 'badge' field to be deleted")
		}
	}
}

func TestAdminService_GetCustomFields_Empty(t *testing.T) {
	svc, db := newAdminSvc(t)
	auth := newAuthSvcForAdmin(t, db)

	user, _ := auth.Register("cfe", "cfe@x.com", "TestPassword123!", "student")

	fields, err := svc.GetCustomFields("user", user.ID, nil)
	if err != nil {
		t.Fatalf("GetCustomFields empty: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(fields))
	}
}
