package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payments-core/internal/config"
	"github.com/payments-core/internal/db"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()

	cfg := config.Config{
		Addr:     ":0",
		DBDriver: "sqlite",
		DBDSN:    "file:memdb1?mode=memory&cache=shared",
	}
	conn, err := db.Open(cfg)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, _ := conn.DB()
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := db.Migrate(cfg, conn); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return NewServer(cfg, conn)
}

func TestCreateAndGetAccount(t *testing.T) {
	s := newTestServer(t)
	h := s.Router()

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"account_id":1,"initial_balance":"100.00"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w2.Code, w2.Body.String())
	}
}

func TestTransferUpdatesBalances(t *testing.T) {
	s := newTestServer(t)
	h := s.Router()

	post := func(path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		return w
	}

	if w := post("/accounts", `{"account_id":1,"initial_balance":"50.00"}`); w.Code != http.StatusCreated {
		t.Fatalf("create src: %d %s", w.Code, w.Body.String())
	}
	if w := post("/accounts", `{"account_id":2,"initial_balance":"10.00"}`); w.Code != http.StatusCreated {
		t.Fatalf("create dst: %d %s", w.Code, w.Body.String())
	}

	if w := post("/transactions", `{"source_account_id":1,"destination_account_id":2,"amount":"15.00"}`); w.Code != http.StatusCreated {
		t.Fatalf("tx: %d %s", w.Code, w.Body.String())
	}

	getBal := func(id string) string {
		req := httptest.NewRequest(http.MethodGet, "/accounts/"+id, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("get %s: %d %s", id, w.Code, w.Body.String())
		}
		return w.Body.String()
	}

	srcBody := getBal("1")
	dstBody := getBal("2")
	if !contains(srcBody, `"balance":"35`) {
		t.Fatalf("unexpected src balance: %s", srcBody)
	}
	if !contains(dstBody, `"balance":"25`) {
		t.Fatalf("unexpected dst balance: %s", dstBody)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
