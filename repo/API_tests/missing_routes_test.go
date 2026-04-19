package api_tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"w2t86/internal/repository"
)

// missing_routes_test.go adds HTTP tests for routes that had no prior coverage:
//   GET    /                                     — root redirect
//   GET    /login                                — login page (form render)
//   GET    /favorites                            — favorites list index
//   DELETE /favorites/:id/items/:materialID      — remove item from list
//   GET    /moderation/items                     — moderation queue partial (HTMX)
//   GET    /dashboard/instructor                 — instructor analytics dashboard
//   POST   /admin/duplicates/merge               — merge duplicate users
//   GET    /analytics/export/distribution        — distribution CSV export

// ---------------------------------------------------------------------------
// GET /  — root redirect
// ---------------------------------------------------------------------------

func TestRoot_Unauthenticated_Redirects(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	resp := makeRequest(app, http.MethodGet, "/", "", "", "")
	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("expected redirect for GET /, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestRoot_Authenticated_Redirects(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/", "", cookie, "")
	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("expected redirect for authenticated GET /, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc == "" {
		t.Fatal("expected Location header in redirect response")
	}
}

// ---------------------------------------------------------------------------
// GET /login — login page form rendering
// ---------------------------------------------------------------------------

func TestLogin_Page_RendersForm(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	resp := makeRequest(app, http.MethodGet, "/login", "", "", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for GET /login, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
	body := readBody(resp)
	lower := strings.ToLower(body)
	if !strings.Contains(lower, "login") && !strings.Contains(lower, "username") &&
		!strings.Contains(lower, "password") {
		t.Errorf("login page should contain login form keywords; got %.200s...", body)
	}
}

func TestLogin_Page_AlreadyAuthenticated(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	// Authenticated user hitting GET /login should still get a 200 (form rendered)
	// or a redirect depending on implementation.  Both are acceptable.
	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/login", "", cookie, "")
	if resp.StatusCode == http.StatusInternalServerError {
		t.Fatalf("server error for authenticated GET /login: %s", readBody(resp))
	}
}

// ---------------------------------------------------------------------------
// GET /favorites — favorites lists index
// ---------------------------------------------------------------------------

func TestFavorites_Index_AuthenticatedReturnsOK(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/favorites", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /favorites (authenticated), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestFavorites_Index_Unauthenticated_Rejected(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	resp := makeRequest(app, http.MethodGet, "/favorites", "", "", "")
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 401/302 for unauthenticated /favorites, got %d", resp.StatusCode)
	}
}

func TestFavorites_Index_WithExistingLists(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")

	var ownerID int64
	_ = db.QueryRow(`SELECT id FROM users ORDER BY id DESC LIMIT 1`).Scan(&ownerID)
	engRepo := repository.NewEngagementRepository(db)
	_, _ = engRepo.CreateList(ownerID, "Reading List", "private")
	_, _ = engRepo.CreateList(ownerID, "Shared Picks", "public")

	resp := makeRequest(app, http.MethodGet, "/favorites", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for /favorites with lists, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

// ---------------------------------------------------------------------------
// DELETE /favorites/:id/items/:materialID — remove item from list
// ---------------------------------------------------------------------------

func TestFavorites_RemoveItem_OwnerSucceeds(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	mat := createMaterial(t, db)

	var ownerID int64
	_ = db.QueryRow(`SELECT id FROM users ORDER BY id DESC LIMIT 1`).Scan(&ownerID)
	engRepo := repository.NewEngagementRepository(db)
	list, err := engRepo.CreateList(ownerID, "Delete Test List", "private")
	if err != nil {
		t.Fatalf("create list: %v", err)
	}
	if err := engRepo.AddToList(list.ID, mat.ID); err != nil {
		t.Fatalf("add to list: %v", err)
	}

	resp := makeRequest(app, http.MethodDelete,
		fmt.Sprintf("/favorites/%d/items/%d", list.ID, mat.ID),
		"", cookie, "", htmx())

	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Fatalf("DELETE /favorites/:id/items/:materialID returned 405 — route not registered")
	}
	if resp.StatusCode/100 != 2 && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 2xx/302 for owner removing item, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}

	// Verify item was removed from the list.
	var count int
	_ = db.QueryRow(
		`SELECT COUNT(*) FROM favorites_items WHERE list_id=? AND material_id=?`,
		list.ID, mat.ID,
	).Scan(&count)
	if count != 0 {
		t.Error("expected material to be removed from favorites list after DELETE")
	}
}

func TestFavorites_RemoveItem_RouteRegistered(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	mat := createMaterial(t, db)
	cookie := loginAs(t, app, db, "student")

	resp := makeRequest(app, http.MethodDelete,
		fmt.Sprintf("/favorites/999999/items/%d", mat.ID),
		"", cookie, "")
	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Fatalf("DELETE /favorites/:id/items/:materialID returned 405 — route not registered")
	}
}

// ---------------------------------------------------------------------------
// GET /moderation/items — moderation queue partial (HTMX)
// ---------------------------------------------------------------------------

func TestModeration_Items_ModeratorAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "moderator")
	resp := makeRequest(app, http.MethodGet, "/moderation/items", "", cookie, "", htmx())
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /moderation/items (moderator), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestModeration_Items_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet, "/moderation/items", "", cookie, "", htmx())
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /moderation/items (admin), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestModeration_Items_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/moderation/items", "", cookie, "")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on /moderation/items, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// GET /dashboard/instructor — instructor analytics dashboard
// ---------------------------------------------------------------------------

func TestDashboard_Instructor_InstructorAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "instructor")
	resp := makeRequest(app, http.MethodGet, "/dashboard/instructor", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /dashboard/instructor (instructor), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestDashboard_Instructor_ManagerAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "manager")
	resp := makeRequest(app, http.MethodGet, "/dashboard/instructor", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /dashboard/instructor (manager), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestDashboard_Instructor_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet, "/dashboard/instructor", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /dashboard/instructor (admin), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestDashboard_Instructor_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/dashboard/instructor", "", cookie, "")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on /dashboard/instructor, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// POST /admin/duplicates/merge — merge duplicate users
// ---------------------------------------------------------------------------

func TestDuplicates_Merge_AdminSucceeds(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	primary := createUser(t, db, "student")
	duplicate := createUser(t, db, "student")
	cookie := loginAs(t, app, db, "admin")

	body := fmt.Sprintf("primary_id=%d&duplicate_id=%d", primary.ID, duplicate.ID)
	resp := makeRequest(app, http.MethodPost, "/admin/duplicates/merge",
		body, cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Fatalf("POST /admin/duplicates/merge returned 405 — route not registered")
	}
	if resp.StatusCode/100 != 2 && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 2xx/302 for admin merge, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestDuplicates_Merge_SameIDReturnsBadRequest(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	user := createUser(t, db, "student")
	cookie := loginAs(t, app, db, "admin")

	body := fmt.Sprintf("primary_id=%d&duplicate_id=%d", user.ID, user.ID)
	resp := makeRequest(app, http.MethodPost, "/admin/duplicates/merge",
		body, cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when merging same user ID, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestDuplicates_Merge_MissingIDReturnsBadRequest(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodPost, "/admin/duplicates/merge",
		"", cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing IDs, got %d", resp.StatusCode)
	}
}

func TestDuplicates_Merge_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodPost, "/admin/duplicates/merge",
		"primary_id=1&duplicate_id=2", cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on merge, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// GET /analytics/export/distribution — distribution CSV export
// ---------------------------------------------------------------------------

func TestAnalytics_ExportDistribution_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet, "/analytics/export/distribution", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for GET /analytics/export/distribution (admin), got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAnalytics_ExportDistribution_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/analytics/export/distribution", "", cookie, "")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on /analytics/export/distribution, got %d",
			resp.StatusCode)
	}
}

func TestAnalytics_ExportDistribution_Unauthenticated(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	resp := makeRequest(app, http.MethodGet, "/analytics/export/distribution", "", "", "")
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 401/302 for unauthenticated export, got %d", resp.StatusCode)
	}
}
