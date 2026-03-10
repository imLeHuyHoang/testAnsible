package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ===================== Unit Tests cho 4 hàm tính toán =====================

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{2, 3, 5},
		{-1, 1, 0},
		{0, 0, 0},
		{100.5, 200.3, 300.8},
	}
	for _, tc := range tests {
		result := Add(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Add(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{5, 3, 2},
		{0, 5, -5},
		{-3, -3, 0},
		{100.5, 0.5, 100},
	}
	for _, tc := range tests {
		result := Subtract(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Subtract(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{2, 3, 6},
		{-2, 3, -6},
		{0, 100, 0},
		{1.5, 2, 3},
	}
	for _, tc := range tests {
		result := Multiply(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Multiply(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
		hasErr   bool
	}{
		{10, 2, 5, false},
		{-6, 3, -2, false},
		{7, 2, 3.5, false},
		{5, 0, 0, true},
	}
	for _, tc := range tests {
		result, err := Divide(tc.a, tc.b)
		if tc.hasErr && err == nil {
			t.Errorf("Divide(%v, %v) expected error, got nil", tc.a, tc.b)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("Divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
		}
		if !tc.hasErr && result != tc.expected {
			t.Errorf("Divide(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

// ===================== Integration Tests cho 4 API endpoints =====================

func TestAPIAdd(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/add?a=10&b=5", nil)
	w := httptest.NewRecorder()
	addHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var result Result
	json.NewDecoder(w.Body).Decode(&result)
	if result.Result != 15 {
		t.Errorf("API add: expected 15, got %v", result.Result)
	}
}

func TestAPISubtract(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/subtract?a=20&b=8", nil)
	w := httptest.NewRecorder()
	subtractHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var result Result
	json.NewDecoder(w.Body).Decode(&result)
	if result.Result != 12 {
		t.Errorf("API subtract: expected 12, got %v", result.Result)
	}
}

func TestAPIMultiply(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/multiply?a=6&b=7", nil)
	w := httptest.NewRecorder()
	multiplyHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var result Result
	json.NewDecoder(w.Body).Decode(&result)
	if result.Result != 42 {
		t.Errorf("API multiply: expected 42, got %v", result.Result)
	}
}

func TestAPIDivide(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/divide?a=100&b=4", nil)
	w := httptest.NewRecorder()
	divideHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var result Result
	json.NewDecoder(w.Body).Decode(&result)
	if result.Result != 25 {
		t.Errorf("API divide: expected 25, got %v", result.Result)
	}
}

func TestAPIDivideByZero(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/divide?a=10&b=0", nil)
	w := httptest.NewRecorder()
	divideHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}

	var errResp ErrorResponse
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp.Error == "" {
		t.Errorf("Expected error message for divide by zero")
	}
}

// ===================== Full Server Integration Test =====================

func TestFullServerIntegration(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/add", addHandler)
	mux.HandleFunc("/api/subtract", subtractHandler)
	mux.HandleFunc("/api/multiply", multiplyHandler)
	mux.HandleFunc("/api/divide", divideHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name     string
		endpoint string
		expected float64
	}{
		{"Add", "/api/add?a=15&b=25", 40},
		{"Subtract", "/api/subtract?a=50&b=30", 20},
		{"Multiply", "/api/multiply?a=8&b=9", 72},
		{"Divide", "/api/divide?a=144&b=12", 12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s%s", server.URL, tc.endpoint))
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected 200, got %d", resp.StatusCode)
			}

			var result Result
			json.NewDecoder(resp.Body).Decode(&result)
			if result.Result != tc.expected {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result.Result)
			}
		})
	}
}
