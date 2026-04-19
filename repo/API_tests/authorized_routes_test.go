package api_tests

import (
	"fmt"
	"net/http"
	"testing"

	"w2t86/internal/repository"
)

// authorized_routes_test.go adds authorized-user tests for routes that previously
// had only rejection-only (403/401) tests:
//
//   GET    /admin/orders                  — clerk can list orders
//   POST   /analytics/map/compute         — admin can trigger grid computation
//   GET    /analytics/map/buffer          — admin can query buffer
//   DELETE /admin/users/:id/fields/:name  — admin can delete legacy user field

// ---------------------------------------------------------------------------
// GET /admin/orders — clerk-authorized
// ---------------------------------------------------------------------------

func TestAdminOrders_ClerkAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "clerk")
	resp := makeRequest(app, http.MethodGet, "/admin/orders", "", cookie, "")
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		t.Fatalf("clerk should be allowed on GET /admin/orders, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for clerk on /admin/orders, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAdminOrders_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet, "/admin/orders", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for admin on /admin/orders, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAdminOrders_WithOrders_ClerkSeesAll(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	student := createUser(t, db, "student")
	mat := createMaterial(t, db)
	orderRepo := repository.NewOrderRepository(db)
	_, err := orderRepo.Create(student.ID, []repository.OrderItemInput{
		{MaterialID: mat.ID, Qty: 1},
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	cookie := loginAs(t, app, db, "clerk")
	resp := makeRequest(app, http.MethodGet, "/admin/orders", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for clerk seeing orders, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

// ---------------------------------------------------------------------------
// POST /analytics/map/compute — admin-authorized grid computation
// ---------------------------------------------------------------------------

func TestAnalytics_ComputeGrid_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodPost, "/analytics/map/compute",
		"layer=&metric=count&grid_size_km=10",
		cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Fatalf("POST /analytics/map/compute returned 405 — route not registered")
	}
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for admin POST /analytics/map/compute, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAnalytics_ComputeGrid_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodPost, "/analytics/map/compute",
		"layer=&metric=count&grid_size_km=10",
		cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on POST /analytics/map/compute, got %d",
			resp.StatusCode)
	}
}

func TestAnalytics_ComputeGrid_DefaultParams(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	// Verify the endpoint works with no body (defaults applied by handler).
	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodPost, "/analytics/map/compute",
		"", cookie, "application/x-www-form-urlencoded")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for compute with defaults, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

// ---------------------------------------------------------------------------
// GET /analytics/map/buffer — admin-authorized buffer query
// ---------------------------------------------------------------------------

func TestAnalytics_BufferQuery_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet,
		"/analytics/map/buffer?lat=39.5&lng=-98.35&radius_km=50",
		"", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for admin on /analytics/map/buffer, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAnalytics_BufferQuery_DefaultParams(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	// Without explicit lat/lng the handler defaults to 0,0.
	cookie := loginAs(t, app, db, "admin")
	resp := makeRequest(app, http.MethodGet, "/analytics/map/buffer", "", cookie, "")
	if resp.StatusCode/100 != 2 {
		t.Fatalf("expected 2xx for /analytics/map/buffer with default params, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}
}

func TestAnalytics_BufferQuery_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	cookie := loginAs(t, app, db, "student")
	resp := makeRequest(app, http.MethodGet, "/analytics/map/buffer", "", cookie, "")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student on /analytics/map/buffer, got %d", resp.StatusCode)
	}
}

func TestAnalytics_BufferQuery_Unauthenticated(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	resp := makeRequest(app, http.MethodGet, "/analytics/map/buffer", "", "", "")
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 401/302 for unauthenticated buffer query, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// DELETE /admin/users/:id/fields/:name — admin legacy alias
// ---------------------------------------------------------------------------

func TestAdminUserFields_LegacyDelete_AdminAllowed(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	target := createUser(t, db, "student")
	cookie := loginAs(t, app, db, "admin")

	// Insert a custom field to delete.
	if _, err := db.Exec(
		`INSERT INTO entity_custom_fields (entity_type, entity_id, field_name, field_value)
		 VALUES ('user', ?, 'legacy_test_field', 'some_value')`,
		target.ID,
	); err != nil {
		t.Fatalf("insert custom field: %v", err)
	}

	resp := makeRequest(app, http.MethodDelete,
		fmt.Sprintf("/admin/users/%d/fields/legacy_test_field", target.ID),
		"", cookie, "", htmx())
	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Fatalf("DELETE /admin/users/:id/fields/:name returned 405 — route not registered")
	}
	if resp.StatusCode/100 != 2 && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 2xx/302 for admin deleting legacy field, got %d; body: %s",
			resp.StatusCode, readBody(resp))
	}

	// Confirm field was removed.
	var count int
	_ = db.QueryRow(
		`SELECT COUNT(*) FROM entity_custom_fields WHERE entity_type='user' AND entity_id=? AND field_name='legacy_test_field'`,
		target.ID,
	).Scan(&count)
	if count != 0 {
		t.Error("expected legacy custom field to be deleted from DB")
	}
}

func TestAdminUserFields_LegacyDelete_StudentForbidden(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	target := createUser(t, db, "student")
	cookie := loginAs(t, app, db, "student")

	resp := makeRequest(app, http.MethodDelete,
		fmt.Sprintf("/admin/users/%d/fields/some_field", target.ID),
		"", cookie, "")
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 403/401 for student deleting field, got %d", resp.StatusCode)
	}
}

func TestAdminUserFields_LegacyDelete_NonExistentField(t *testing.T) {
	app, db, cleanup := newTestApp(t)
	defer cleanup()

	target := createUser(t, db, "student")
	cookie := loginAs(t, app, db, "admin")

	resp := makeRequest(app, http.MethodDelete,
		fmt.Sprintf("/admin/users/%d/fields/nonexistent_field", target.ID),
		"", cookie, "", htmx())
	// Handler should succeed (idempotent delete) or return 404/400.
	if resp.StatusCode == http.StatusInternalServerError {
		t.Fatalf("server error deleting non-existent field: %s", readBody(resp))
	}
}
