package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"w2t86/internal/config"
	"w2t86/internal/middleware"
	"w2t86/internal/models"
	"w2t86/internal/repository"
	"w2t86/internal/services"
	"w2t86/internal/testutil"
)

// ---------------------------------------------------------------------------
// RequireRole tests
// ---------------------------------------------------------------------------

func TestRequireRole_AllowsMatchingRole(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &models.User{ID: 1, Role: "admin"})
		return c.Next()
	})
	app.Get("/admin", middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for admin role, got %d", resp.StatusCode)
	}
}

func TestRequireRole_AllowsOneOfMultipleRoles(t *testing.T) {
	for _, role := range []string{"clerk", "admin"} {
		role := role
		t.Run(role, func(t *testing.T) {
			a := fiber.New()
			a.Use(func(c *fiber.Ctx) error {
				c.Locals("user", &models.User{ID: 1, Role: role})
				return c.Next()
			})
			a.Get("/orders", middleware.RequireRole("clerk", "admin"), func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			resp, err := a.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200 for role %q, got %d", role, resp.StatusCode)
			}
		})
	}
}

func TestRequireRole_DeniesUnauthorizedRole(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &models.User{ID: 1, Role: "student"})
		return c.Next()
	})
	app.Get("/admin", middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for student on admin route, got %d", resp.StatusCode)
	}
}

func TestRequireRole_DeniesWhenNoUser(t *testing.T) {
	app := fiber.New()
	// No middleware sets the user — simulates missing auth.
	app.Get("/admin", middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 when user is nil, got %d", resp.StatusCode)
	}
}

func TestRequireRole_MultiRoleDeniesOthers(t *testing.T) {
	for _, role := range []string{"student", "instructor", "moderator"} {
		role := role
		t.Run(role, func(t *testing.T) {
			a := fiber.New()
			a.Use(func(c *fiber.Ctx) error {
				c.Locals("user", &models.User{ID: 1, Role: role})
				return c.Next()
			})
			a.Get("/clerk-only", middleware.RequireRole("clerk", "admin"), func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/clerk-only", nil)
			resp, err := a.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusForbidden {
				t.Fatalf("expected 403 for role %q on clerk-only route, got %d", role, resp.StatusCode)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// RequireAuth tests
// ---------------------------------------------------------------------------

func TestRequireAuth_RejectsEmptyCookie(t *testing.T) {
	db := testutil.NewTestDB(t)
	authMW := middleware.NewAuthMiddleware(
		repository.NewSessionRepository(db),
		repository.NewUserRepository(db),
		"test-secret",
	)

	app := fiber.New()
	app.Get("/protected", authMW.RequireAuth(), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without cookie, got %d", resp.StatusCode)
	}
}

func TestRequireAuth_RejectsInvalidToken(t *testing.T) {
	db := testutil.NewTestDB(t)
	authMW := middleware.NewAuthMiddleware(
		repository.NewSessionRepository(db),
		repository.NewUserRepository(db),
		"test-secret",
	)

	app := fiber.New()
	app.Get("/protected", authMW.RequireAuth(), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", "session_token=bogus-invalid-token")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for bogus token, got %d", resp.StatusCode)
	}
}

func TestRequireAuth_AllowsValidSession(t *testing.T) {
	db := testutil.NewTestDB(t)
	cfg := &config.Config{SessionSecret: "test-secret"}
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authSvc := services.NewAuthService(userRepo, sessionRepo, cfg)

	// Create and log in a user to get a real session token.
	if _, err := authSvc.Register("mwuser", "mw@example.com", "TestPassword123!", "student"); err != nil {
		t.Fatalf("register: %v", err)
	}
	token, user, err := authSvc.Login("mwuser", "TestPassword123!")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	authMW := middleware.NewAuthMiddleware(sessionRepo, userRepo, cfg.SessionSecret)

	app := fiber.New()
	app.Get("/protected", authMW.RequireAuth(), func(c *fiber.Ctx) error {
		u := middleware.GetUser(c)
		if u == nil || u.ID != user.ID {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Cookie", "session_token="+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for valid session, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// GetUser helper test
// ---------------------------------------------------------------------------

func TestGetUser_ReturnsNilWhenNotSet(t *testing.T) {
	app := fiber.New()
	var gotUser *models.User
	app.Get("/check", func(c *fiber.Ctx) error {
		gotUser = middleware.GetUser(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	if _, err := app.Test(req, -1); err != nil {
		t.Fatal(err)
	}
	if gotUser != nil {
		t.Errorf("expected nil user when not set, got %+v", gotUser)
	}
}

func TestGetUser_ReturnsUserWhenSet(t *testing.T) {
	expected := &models.User{ID: 42, Role: "admin"}
	app := fiber.New()
	var gotUser *models.User

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", expected)
		return c.Next()
	})
	app.Get("/check", func(c *fiber.Ctx) error {
		gotUser = middleware.GetUser(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	if _, err := app.Test(req, -1); err != nil {
		t.Fatal(err)
	}
	if gotUser == nil || gotUser.ID != expected.ID {
		t.Errorf("expected user ID=%d, got %+v", expected.ID, gotUser)
	}
}
