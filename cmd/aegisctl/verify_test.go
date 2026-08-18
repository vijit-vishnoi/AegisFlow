package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestVerifyHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Strip test args
	args := []string{os.Args[0]}
	for i, arg := range os.Args {
		if arg == "--" {
			args = append(args, os.Args[i+1:]...)
			break
		}
	}
	os.Args = args
	main()
	os.Exit(0)
}

func runVerifyCommand(adminURL string, args ...string) (string, string, int) {
	cmdArgs := append([]string{"-test.run=TestVerifyHelperProcess", "--", "verify"}, args...)
	cmd := exec.Command(os.Args[0], cmdArgs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "AEGISFLOW_ADMIN_URL="+adminURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

func TestVerifyJSON_ValidChain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/v1/audit/verify" {
			t.Errorf("expected path /admin/v1/audit/verify, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"valid": true, "total_records": 10, "message": "all good"}`))
	}))
	defer ts.Close()

	stdout, stderr, code := runVerifyCommand(ts.URL, "--json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d. stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %s", stderr)
	}

	if strings.Contains(stdout, "\x1b") {
		t.Fatalf("output must not contain ANSI escapes: %q", stdout)
	}

	var res VerifyResponse
	if err := json.Unmarshal([]byte(stdout), &res); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if !res.Valid {
		t.Fatal("expected valid to be true")
	}
}

func TestVerifyJSON_InvalidChain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"valid": false, "total_records": 10, "message": "bad hash"}`))
	}))
	defer ts.Close()

	stdout, stderr, code := runVerifyCommand(ts.URL, "--json")
	if code != 1 {
		t.Fatalf("expected exit code 1 for invalid chain, got %d", code)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %s", stderr)
	}

	var res VerifyResponse
	if err := json.Unmarshal([]byte(stdout), &res); err != nil {
		t.Fatalf("stdout is not valid JSON: %v (output: %q)", err, stdout)
	}
	if res.Valid {
		t.Fatal("expected valid to be false")
	}
}

func TestVerifyJSON_SessionEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/v1/evidence/sessions/ses_123/verify" {
			t.Errorf("expected path /admin/v1/evidence/sessions/ses_123/verify, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"valid": true, "total_records": 2, "message": "session valid"}`))
	}))
	defer ts.Close()

	stdout, stderr, code := runVerifyCommand(ts.URL, "--session", "ses_123", "--json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d. stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %s", stderr)
	}

	var res VerifyResponse
	if err := json.Unmarshal([]byte(stdout), &res); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if !res.Valid {
		t.Fatal("expected valid to be true")
	}
}

func TestVerifyJSON_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	stdout, stderr, code := runVerifyCommand(ts.URL, "--json")
	if code != 1 {
		t.Fatalf("expected exit code 1 for http error, got %d", code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %s", stdout)
	}
	if !strings.Contains(stderr, "Error (500)") {
		t.Fatalf("expected Error (500) in stderr, got %s", stderr)
	}
}

func TestFormatVerifyResult_Pass(t *testing.T) {
	r := VerifyResponse{
		Valid:        true,
		TotalRecords: 42,
		Message:      "evidence chain integrity verified",
	}
	out := formatVerifyResult(r)
	if !strings.Contains(out, "PASS") {
		t.Fatal("expected PASS in output")
	}
	if !strings.Contains(out, "Total entries: 42") {
		t.Fatal("expected total entries in output")
	}
	if strings.Contains(out, "Error at index") {
		t.Fatal("should not show error index for valid result")
	}
}

func TestFormatVerifyResult_Fail(t *testing.T) {
	r := VerifyResponse{
		Valid:        false,
		TotalRecords: 10,
		ErrorAtIndex: 5,
		Message:      "hash mismatch at record abc123",
	}
	out := formatVerifyResult(r)
	if !strings.Contains(out, "FAIL") {
		t.Fatal("expected FAIL in output")
	}
	if !strings.Contains(out, "Error at index: 5") {
		t.Fatal("expected error index in output")
	}
	if !strings.Contains(out, "hash mismatch") {
		t.Fatal("expected message in output")
	}
}

func TestFormatVerifyResult_EmptyChain(t *testing.T) {
	r := VerifyResponse{
		Valid:        true,
		TotalRecords: 0,
		Message:      "empty chain is valid",
	}
	out := formatVerifyResult(r)
	if !strings.Contains(out, "PASS") {
		t.Fatal("expected PASS for empty chain")
	}
	if !strings.Contains(out, "Total entries: 0") {
		t.Fatal("expected zero entries")
	}
}
