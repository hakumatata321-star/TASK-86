package services_test

import (
	"database/sql"
	"testing"

	"w2t86/internal/models"
	"w2t86/internal/repository"
	"w2t86/internal/services"
	"w2t86/internal/testutil"
)

// insertModerationFixture inserts a user, material, collapsed comment, and
// a report; returns (userID, commentID).
func insertModerationFixture(t *testing.T, db *sql.DB, suffix string) (int64, int64) {
	t.Helper()
	matRepo := repository.NewMaterialRepository(db)

	var userID int64
	if err := db.QueryRow(
		`INSERT INTO users (username, email, password_hash, role) VALUES (?,?,?,?) RETURNING id`,
		"modu"+suffix, "mod"+suffix+"@x.com", "hash", "student",
	).Scan(&userID); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	mat, err := matRepo.Create(&models.Material{
		Title: "Mod Book", TotalQty: 1, AvailableQty: 1, Price: 1.0, Status: "active",
	})
	if err != nil {
		t.Fatalf("create material: %v", err)
	}
	var commentID int64
	if err := db.QueryRow(
		`INSERT INTO comments (user_id, material_id, body, status) VALUES (?,?,'Body','collapsed') RETURNING id`,
		userID, mat.ID,
	).Scan(&commentID); err != nil {
		t.Fatalf("insert comment: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO comment_reports (comment_id, reported_by, reason) VALUES (?,?,'spam')`,
		commentID, userID,
	); err != nil {
		t.Fatalf("insert report: %v", err)
	}
	return userID, commentID
}

func TestModerationService_GetQueue_Empty(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := services.NewModerationService(repository.NewModerationRepository(db))

	items, err := svc.GetQueue(10, 0)
	if err != nil {
		t.Fatalf("GetQueue: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected empty queue, got %d items", len(items))
	}
}

func TestModerationService_CountQueue_Zero(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := services.NewModerationService(repository.NewModerationRepository(db))

	n, err := svc.CountQueue()
	if err != nil {
		t.Fatalf("CountQueue: %v", err)
	}
	if n != 0 {
		t.Errorf("expected count=0, got %d", n)
	}
}

func TestModerationService_GetQueue_WithItem(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := services.NewModerationService(repository.NewModerationRepository(db))

	_, _ = insertModerationFixture(t, db, "q1")

	items, err := svc.GetQueue(10, 0)
	if err != nil {
		t.Fatalf("GetQueue: %v", err)
	}
	if len(items) == 0 {
		t.Error("expected at least 1 item in moderation queue")
	}

	count, err := svc.CountQueue()
	if err != nil {
		t.Fatalf("CountQueue: %v", err)
	}
	if count == 0 {
		t.Error("expected CountQueue > 0")
	}
}

func TestModerationService_ApproveComment(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := services.NewModerationService(repository.NewModerationRepository(db))

	userID, commentID := insertModerationFixture(t, db, "appr")

	if err := svc.ApproveComment(commentID, userID); err != nil {
		t.Fatalf("ApproveComment: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM comments WHERE id=?`, commentID).Scan(&status); err != nil {
		t.Fatalf("query comment: %v", err)
	}
	if status != "visible" {
		t.Errorf("expected status=visible after approve, got %q", status)
	}
}

func TestModerationService_RemoveComment(t *testing.T) {
	db := testutil.NewTestDB(t)
	svc := services.NewModerationService(repository.NewModerationRepository(db))

	userID, commentID := insertModerationFixture(t, db, "rem")

	if err := svc.RemoveComment(commentID, userID); err != nil {
		t.Fatalf("RemoveComment: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM comments WHERE id=?`, commentID).Scan(&status); err != nil {
		t.Fatalf("query comment: %v", err)
	}
	if status != "removed" {
		t.Errorf("expected status=removed after remove, got %q", status)
	}
}
