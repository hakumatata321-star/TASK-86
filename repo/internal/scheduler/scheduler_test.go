package scheduler

import (
	"testing"
	"time"

	"w2t86/internal/testutil"
)

// TestAutoCloseOrders_CancelsExpiredPendingPayment verifies that an order stuck
// in pending_payment whose auto_close_at has elapsed is canceled and its
// inventory is rolled back.
func TestAutoCloseOrders_CancelsExpiredPendingPayment(t *testing.T) {
	db := testutil.NewTestDBNoFK(t)
	s := NewOrderScheduler(db)

	// Insert a material with 5 total / 0 available (all reserved).
	var materialID int64
	if err := db.QueryRow(`
		INSERT INTO materials (title, total_qty, available_qty, reserved_qty, price, status)
		VALUES ('Test Book', 5, 0, 5, 9.99, 'active')
		RETURNING id`).Scan(&materialID); err != nil {
		t.Fatalf("insert material: %v", err)
	}

	// Insert a user (FK off, so we can use a placeholder).
	const userID int64 = 999

	// Insert an order in pending_payment whose auto_close_at is in the past.
	var orderID int64
	if err := db.QueryRow(`
		INSERT INTO orders (user_id, status, total_amount, auto_close_at)
		VALUES (?, 'pending_payment', 49.95, datetime('now', '-1 hour'))
		RETURNING id`, userID).Scan(&orderID); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	// Insert an order item.
	if _, err := db.Exec(`
		INSERT INTO order_items (order_id, material_id, qty, unit_price)
		VALUES (?, ?, 5, 9.99)`, orderID, materialID); err != nil {
		t.Fatalf("insert order_item: %v", err)
	}

	// Run the auto-close job.
	if err := s.autoCloseOrders("pending_payment"); err != nil {
		t.Fatalf("autoCloseOrders: %v", err)
	}

	// Verify the order is now canceled.
	var status string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id=?`, orderID).Scan(&status); err != nil {
		t.Fatalf("query order status: %v", err)
	}
	if status != "canceled" {
		t.Errorf("expected status=canceled, got %q", status)
	}

	// Verify inventory was rolled back: available_qty should be 5, reserved 0.
	var availQty, reservedQty int
	if err := db.QueryRow(`SELECT available_qty, reserved_qty FROM materials WHERE id=?`,
		materialID).Scan(&availQty, &reservedQty); err != nil {
		t.Fatalf("query material qty: %v", err)
	}
	if availQty != 5 {
		t.Errorf("expected available_qty=5 after rollback, got %d", availQty)
	}
	if reservedQty != 0 {
		t.Errorf("expected reserved_qty=0 after rollback, got %d", reservedQty)
	}

	// Verify an order_event was recorded.
	var eventCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM order_events WHERE order_id=? AND to_status='canceled'`,
		orderID).Scan(&eventCount); err != nil {
		t.Fatalf("query order_events: %v", err)
	}
	if eventCount != 1 {
		t.Errorf("expected 1 order_event for canceled order, got %d", eventCount)
	}
}

// TestAutoCloseOrders_SkipsFutureAutoClose verifies that an order whose
// auto_close_at is in the future is NOT canceled.
func TestAutoCloseOrders_SkipsFutureAutoClose(t *testing.T) {
	db := testutil.NewTestDBNoFK(t)
	s := NewOrderScheduler(db)

	const userID int64 = 998
	var orderID int64
	if err := db.QueryRow(`
		INSERT INTO orders (user_id, status, total_amount, auto_close_at)
		VALUES (?, 'pending_payment', 9.99, datetime('now', '+1 hour'))
		RETURNING id`, userID).Scan(&orderID); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	if err := s.autoCloseOrders("pending_payment"); err != nil {
		t.Fatalf("autoCloseOrders: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id=?`, orderID).Scan(&status); err != nil {
		t.Fatalf("query order status: %v", err)
	}
	if status != "pending_payment" {
		t.Errorf("expected order to remain pending_payment, got %q", status)
	}
}

// TestAutoCloseOrders_NullAutoClose verifies orders without auto_close_at set
// are never auto-closed.
func TestAutoCloseOrders_NullAutoClose(t *testing.T) {
	db := testutil.NewTestDBNoFK(t)
	s := NewOrderScheduler(db)

	const userID int64 = 997
	var orderID int64
	if err := db.QueryRow(`
		INSERT INTO orders (user_id, status, total_amount)
		VALUES (?, 'pending_payment', 9.99)
		RETURNING id`, userID).Scan(&orderID); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	if err := s.autoCloseOrders("pending_payment"); err != nil {
		t.Fatalf("autoCloseOrders: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM orders WHERE id=?`, orderID).Scan(&status); err != nil {
		t.Fatalf("query order status: %v", err)
	}
	if status == "canceled" {
		t.Error("expected order with NULL auto_close_at not to be canceled")
	}
}

// TestAutoCloseOrders_CancelsExpiredPendingShipment verifies the shipment-
// timeout note is correct.
func TestAutoCloseOrders_CancelsExpiredPendingShipment(t *testing.T) {
	db := testutil.NewTestDBNoFK(t)
	s := NewOrderScheduler(db)

	const userID int64 = 996
	var orderID int64
	if err := db.QueryRow(`
		INSERT INTO orders (user_id, status, total_amount, auto_close_at)
		VALUES (?, 'pending_shipment', 9.99, datetime('now', '-2 hours'))
		RETURNING id`, userID).Scan(&orderID); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	if err := s.autoCloseOrders("pending_shipment"); err != nil {
		t.Fatalf("autoCloseOrders: %v", err)
	}

	var note string
	if err := db.QueryRow(
		`SELECT note FROM order_events WHERE order_id=? AND to_status='canceled'`,
		orderID,
	).Scan(&note); err != nil {
		t.Fatalf("query event note: %v", err)
	}
	if note != "auto-closed: shipment timeout" {
		t.Errorf("expected shipment-timeout note, got %q", note)
	}
}

// TestNewOrderScheduler_StartStop verifies the scheduler can be started and
// stopped without panicking.
func TestNewOrderScheduler_StartStop(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewOrderScheduler(db)

	s.Start()

	// Give the cron a moment to initialise.
	time.Sleep(50 * time.Millisecond)

	// Stop should complete without hanging.
	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop() timed out after 5 seconds")
	}
}
