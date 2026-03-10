package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Result struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Add(a, b float64) float64 {
	return a + b
}

func Subtract(a, b float64) float64 {
	return a - b
}

func Multiply(a, b float64) float64 {
	return a * b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func parseParams(r *http.Request) (float64, float64, error) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid parameter 'a': %s", aStr)
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid parameter 'b': %s", bStr)
	}

	return a, b, nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	a, b, err := parseParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, Result{Operation: "add", A: a, B: b, Result: Add(a, b)})
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	a, b, err := parseParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, Result{Operation: "subtract", A: a, B: b, Result: Subtract(a, b)})
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	a, b, err := parseParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, Result{Operation: "multiply", A: a, B: b, Result: Multiply(a, b)})
}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	a, b, err := parseParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	result, err := Divide(a, b)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, Result{Operation: "divide", A: a, B: b, Result: result})
}

func main() {
	http.HandleFunc("/api/add", addHandler)
	http.HandleFunc("/api/subtract", subtractHandler)
	http.HandleFunc("/api/multiply", multiplyHandler)
	http.HandleFunc("/api/divide", divideHandler)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
