package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func main() {
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		fmt.Fprintln(os.Stderr, "ADMIN_PASSWORD env var required")
		os.Exit(1)
	}
	baseURL := "http://localhost:3000"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// Obtain CSRF token.
	resp, err := client.Get(baseURL + "/login")
	if err != nil {
		fmt.Fprintf(os.Stderr, "GET /login: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()

	csrf := extractCookie(jar, baseURL, "csrf_token")

	// Admin login.
	form := url.Values{"username": {"admin"}, "password": {adminPass}}
	req, _ := http.NewRequest("POST", baseURL+"/login", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Csrf-Token", csrf)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "POST /login: %v\n", err)
		os.Exit(1)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "admin login failed: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	csrf = extractCookie(jar, baseURL, "csrf_token")

	accounts := []struct{ username, email, role string }{
		{"demo_student", "demo_student@portal.local", "student"},
		{"demo_instructor", "demo_instructor@portal.local", "instructor"},
		{"demo_clerk", "demo_clerk@portal.local", "clerk"},
		{"demo_moderator", "demo_moderator@portal.local", "moderator"},
		{"demo_manager", "demo_manager@portal.local", "manager"},
	}

	for _, a := range accounts {
		csrf = createUser(client, baseURL, csrf, a.username, a.email, "Demo1234!", a.role)
	}
}

func createUser(client *http.Client, baseURL, csrf, username, email, password, role string) string {
	form := url.Values{
		"username": {username}, "email": {email},
		"password": {password}, "role": {role},
	}
	req, _ := http.NewRequest("POST", baseURL+"/admin/users", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Csrf-Token", csrf)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create user %s: %v\n", username, err)
		return csrf
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 302 || resp.StatusCode == 303 {
		fmt.Printf("  OK   %s (%s)\n", username, role)
	} else {
		fmt.Printf("  SKIP %s — HTTP %d\n", username, resp.StatusCode)
	}
	jar := client.Jar
	return extractCookie(jar, baseURL, "csrf_token")
}

func extractCookie(jar http.CookieJar, rawURL, name string) string {
	u, _ := url.Parse(rawURL)
	for _, c := range jar.Cookies(u) {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}
