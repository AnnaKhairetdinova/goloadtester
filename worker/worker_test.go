package worker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	defer server.Close()

	result := Do(server.URL, 1*time.Second)
	if result.Err != nil {
		t.Fatalf("Ошибка: %s", result.Err)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен: %d", result.StatusCode)
	}

	if result.Duration <= 0 {
		t.Errorf("Длительность должна быть положительной, но получена: %d", result.Duration)
	}
}

func TestDo_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	result := Do(server.URL, 100*time.Millisecond)
	if result.Err == nil {
		t.Fatal("Нет ошишбки таймаута")
	}

	if result.Duration > 200*time.Millisecond {
		t.Errorf("Слишком долгая длительность таймаута: %v", result.Duration)
	}
}
