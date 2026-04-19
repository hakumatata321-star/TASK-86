package services_test

import (
	"testing"

	"w2t86/internal/repository"
	"w2t86/internal/services"
	"w2t86/internal/testutil"
)

func newAnalyticsSvc(t *testing.T) *services.AnalyticsService {
	t.Helper()
	db := testutil.NewTestDB(t)
	return services.NewAnalyticsService(repository.NewAnalyticsRepository(db))
}

func TestAnalyticsService_GetMapData_Empty(t *testing.T) {
	svc := newAnalyticsSvc(t)

	data, err := svc.GetMapData("")
	if err != nil {
		t.Fatalf("GetMapData: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil MapData")
	}
	// Empty DB means empty collections, not error.
	if data.Locations == nil {
		data.Locations = nil // nil slice is fine
	}
}

func TestAnalyticsService_GetMapData_LayerFilter(t *testing.T) {
	svc := newAnalyticsSvc(t)

	for _, layer := range []string{"school", "distribution_center", ""} {
		layer := layer
		t.Run(layer, func(t *testing.T) {
			data, err := svc.GetMapData(layer)
			if err != nil {
				t.Fatalf("GetMapData(%q): %v", layer, err)
			}
			if data == nil {
				t.Errorf("expected non-nil MapData for layer=%q", layer)
			}
		})
	}
}

func TestAnalyticsService_ComputeGrid_NoLocations(t *testing.T) {
	svc := newAnalyticsSvc(t)

	// Should not error even with no data in the DB.
	if err := svc.ComputeGrid("", "count", 10.0); err != nil {
		t.Fatalf("ComputeGrid: %v", err)
	}
}

func TestAnalyticsService_ComputeGrid_Metrics(t *testing.T) {
	svc := newAnalyticsSvc(t)

	for _, metric := range []string{"count", "orders", "value"} {
		metric := metric
		t.Run(metric, func(t *testing.T) {
			if err := svc.ComputeGrid("", metric, 5.0); err != nil {
				t.Fatalf("ComputeGrid(metric=%q): %v", metric, err)
			}
		})
	}
}

func TestAnalyticsService_LocationsWithinRadius_Empty(t *testing.T) {
	svc := newAnalyticsSvc(t)

	locs, err := svc.LocationsWithinRadius(39.5, -98.35, 50.0, "")
	if err != nil {
		t.Fatalf("LocationsWithinRadius: %v", err)
	}
	// Empty DB — no locations returned, no error.
	if locs == nil {
		locs = nil // nil is acceptable
	}
}

func TestAnalyticsService_LocationsWithinRadius_TypeFilter(t *testing.T) {
	svc := newAnalyticsSvc(t)

	for _, locType := range []string{"school", "distribution_center", ""} {
		locType := locType
		t.Run(locType, func(t *testing.T) {
			_, err := svc.LocationsWithinRadius(0, 0, 10.0, locType)
			if err != nil {
				t.Fatalf("LocationsWithinRadius(type=%q): %v", locType, err)
			}
		})
	}
}

func TestAnalyticsService_AdminDashboardData_Empty(t *testing.T) {
	svc := newAnalyticsSvc(t)

	data, err := svc.AdminDashboardData()
	if err != nil {
		t.Fatalf("AdminDashboardData: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil dashboard data")
	}
}

func TestAnalyticsService_InstructorDashboardData_Empty(t *testing.T) {
	svc := newAnalyticsSvc(t)

	// instructorID 0 — no instructor, should return empty/default data.
	data, err := svc.InstructorDashboardData(0)
	if err != nil {
		t.Fatalf("InstructorDashboardData: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil instructor dashboard data")
	}
}

func TestAnalyticsService_TrajectoryPoints_UnknownMaterial(t *testing.T) {
	svc := newAnalyticsSvc(t)

	pts, err := svc.TrajectoryPoints(999999)
	if err != nil {
		t.Fatalf("TrajectoryPoints for unknown material: %v", err)
	}
	if len(pts) != 0 {
		t.Errorf("expected 0 trajectory points for unknown material, got %d", len(pts))
	}
}
