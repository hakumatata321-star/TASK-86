package services_test

import (
	"testing"

	"w2t86/internal/repository"
	"w2t86/internal/services"
	"w2t86/internal/testutil"
)

func newCourseService(t *testing.T) (*services.CourseService, int64) {
	t.Helper()
	db := testutil.NewTestDB(t)
	courseRepo := repository.NewCourseRepository(db)
	matRepo := repository.NewMaterialRepository(db)
	svc := services.NewCourseService(courseRepo, matRepo)

	// Insert an instructor user.
	var instructorID int64
	if err := db.QueryRow(
		`INSERT INTO users (username, email, password_hash, role) VALUES ('instr1','i1@x.com','hash','instructor') RETURNING id`,
	).Scan(&instructorID); err != nil {
		t.Fatalf("insert instructor: %v", err)
	}
	return svc, instructorID
}

func TestCourseService_CreateCourse_Valid(t *testing.T) {
	svc, instructorID := newCourseService(t)

	course, err := svc.CreateCourse(instructorID, "Math 101", "math", "grade-6", "2025-2026")
	if err != nil {
		t.Fatalf("CreateCourse: %v", err)
	}
	if course.ID == 0 {
		t.Error("expected non-zero course ID")
	}
	if course.Name != "Math 101" {
		t.Errorf("expected name=Math 101, got %q", course.Name)
	}
	if course.InstructorID != instructorID {
		t.Errorf("expected instructorID=%d, got %d", instructorID, course.InstructorID)
	}
}

func TestCourseService_CreateCourse_EmptyName(t *testing.T) {
	svc, instructorID := newCourseService(t)

	_, err := svc.CreateCourse(instructorID, "", "math", "", "")
	if err == nil {
		t.Error("expected error for empty course name, got nil")
	}
}

func TestCourseService_ListCourses_Empty(t *testing.T) {
	svc, instructorID := newCourseService(t)

	courses, err := svc.ListCourses(instructorID)
	if err != nil {
		t.Fatalf("ListCourses: %v", err)
	}
	if len(courses) != 0 {
		t.Errorf("expected empty course list, got %d", len(courses))
	}
}

func TestCourseService_ListCourses_NonEmpty(t *testing.T) {
	svc, instructorID := newCourseService(t)

	_, _ = svc.CreateCourse(instructorID, "Science", "science", "grade-7", "2025")
	_, _ = svc.CreateCourse(instructorID, "English", "english", "grade-7", "2025")

	courses, err := svc.ListCourses(instructorID)
	if err != nil {
		t.Fatalf("ListCourses: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("expected 2 courses, got %d", len(courses))
	}
}

func TestCourseService_GetCourse_OwnerAccess(t *testing.T) {
	svc, instructorID := newCourseService(t)

	created, err := svc.CreateCourse(instructorID, "History", "history", "grade-8", "2025")
	if err != nil {
		t.Fatalf("CreateCourse: %v", err)
	}

	fetched, err := svc.GetCourse(created.ID, instructorID, false)
	if err != nil {
		t.Fatalf("GetCourse owner: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected course ID=%d, got %d", created.ID, fetched.ID)
	}
}

func TestCourseService_GetCourse_NonOwnerDenied(t *testing.T) {
	svc, instructorID := newCourseService(t)

	created, _ := svc.CreateCourse(instructorID, "Art", "art", "grade-5", "2025")

	_, err := svc.GetCourse(created.ID, instructorID+99, false)
	if err == nil {
		t.Error("expected error for non-owner accessing course, got nil")
	}
}

func TestCourseService_GetCourse_AdminBypasses(t *testing.T) {
	svc, instructorID := newCourseService(t)

	created, _ := svc.CreateCourse(instructorID, "PE", "pe", "grade-4", "2025")

	fetched, err := svc.GetCourse(created.ID, instructorID+99, true)
	if err != nil {
		t.Fatalf("admin GetCourse: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected course ID=%d, got %d", created.ID, fetched.ID)
	}
}

func TestCourseService_AddSection_Valid(t *testing.T) {
	svc, instructorID := newCourseService(t)

	course, _ := svc.CreateCourse(instructorID, "Biology", "biology", "grade-9", "2025")
	section, err := svc.AddSection(course.ID, "Period 1", "1", "Room 101")
	if err != nil {
		t.Fatalf("AddSection: %v", err)
	}
	if section.ID == 0 {
		t.Error("expected non-zero section ID")
	}
	if section.Name != "Period 1" {
		t.Errorf("expected name=Period 1, got %q", section.Name)
	}
}

func TestCourseService_AddSection_EmptyName(t *testing.T) {
	svc, instructorID := newCourseService(t)

	course, _ := svc.CreateCourse(instructorID, "Chemistry", "chem", "grade-10", "2025")
	_, err := svc.AddSection(course.ID, "", "", "")
	if err == nil {
		t.Error("expected error for empty section name, got nil")
	}
}
