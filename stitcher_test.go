package stitcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_it_calls_the_services_and_concatenates_the_results(t *testing.T) {

	helloAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}))

	worldAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "world")
	}))

	result := Stitcher(helloAPI.URL, worldAPI.URL)
	expected := "Hello, world"

	if result != expected {
		t.Errorf("Stitcher failed, expected [%s], got [%s]", expected, result)
	}
}

func Test_it_handles_non_ok_from_hello_api(t *testing.T) {

	helloAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "oops", http.StatusInternalServerError)
	}))

	worldAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "world")
	}))

	result := Stitcher(helloAPI.URL, worldAPI.URL)

	if result != errorMsg {
		t.Errorf("Stitcher didnt fail properly, expected [%s], got [%s]", errorMsg, result)
	}
}

func Benchmark_the_stitcher(b *testing.B) {
	slowHelloAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(w, "Hello")
	}))

	slowWorldAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(w, "world")
	}))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Stitcher(slowHelloAPI.URL, slowWorldAPI.URL)
	}
}
