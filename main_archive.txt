package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/getNeonConnection"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/testNeonDatabase"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/mainPage"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goConstants"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goEnv"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
)

const clerkSessionCookieName = "clerkSession"
const clerkAPIBaseURL = "https://api.clerk.com/v1"

func main() {
	err := goEnv.GetEnvVar()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	clerk.SetKey(goEnv.GlobalEnvVar.ClerkSecretKey)

	ctx := context.Background()

	// Initialize Gin
	r := gin.Default()

	neonConnection := getNeonConnection.GetNeonConnection(ctx)
	defer neonConnection.Close(ctx)

	// Serve static files from the "static" directory
	r.Static("/"+goConstants.StaticFolder, "./"+goConstants.StaticFolder)

	// Serve the initial page
	r.GET("/test-neon-database", func(c *gin.Context) {
		testNeonDatabase.TestNeonDatabase(c.Request.Context(), neonConnection)
	})

	// Serve the initial page
	r.GET("/", func(c *gin.Context) {
		cookie, err := c.Cookie(clerkSessionCookieName)
		if err != nil {
			fmt.Println("No cookie found")
		} else {
			fmt.Printf("Cookie value: %s\n", cookie)
		}

		c.Status(200)
		mainPage.Page().Render(c.Request.Context(), c.Writer)
	})

	// Handle HTMX request to update the greeting
	r.GET("/update", func(c *gin.Context) {
		c.Status(200)
		mainPage.Hello("HTMX").Render(c.Request.Context(), c.Writer)
	})

	r.POST("/signup", SignUpHandler)
	r.POST("/signin", SignInHandler)
	r.POST("/signout", SignOutHandler)

	// Protected route
	r.GET("/protected", ProtectedRoute)

	// must run on 0.0.0.0, not localhost or anything else so that it works on fly.io
	r.Run("0.0.0.0:8080")
}

func GetSessionToken(c *gin.Context) (string, error) {
	// Attempt to retrieve the session token from the cookie
	sessionToken, err := c.Cookie(clerkSessionCookieName)
	if err != nil {
		return "", fmt.Errorf("session token not found in cookies")
	}
	return sessionToken, nil
}

func clerkAPIRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	apiKey := goEnv.GlobalEnvVar.ClerkSecretKey
	url := clerkAPIBaseURL + endpoint

	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	// Log the request and response for debugging
	fmt.Printf("Request: %s %s\n", method, url)
	if payload != nil {
		fmt.Printf("Payload: %s\n", string(body))
	}
	if resp != nil {
		fmt.Printf("Response Status: %d\n", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response Body: %s\n", string(bodyBytes))
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reattach body for further use
	}

	return resp, err
}

func SignUpHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Create the user
	payload := map[string]interface{}{
		"email_address": []string{email},
		"password":      password,
	}
	resp, err := clerkAPIRequest("POST", "/users", payload)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to sign up: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Sign up failed: %s", string(bodyBytes)))
		return
	}

	// Parse the user_id from the response
	var signupResponse struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&signupResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse sign-up response")
		return
	}

	// Automatically sign in the user
	signIn(c, signupResponse.ID)
}

func signIn(c *gin.Context, userID string) {
	// Create a session
	payload := map[string]string{
		"user_id": userID,
	}

	resp, err := clerkAPIRequest("POST", "/sessions", payload)
	if err != nil || resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to sign in: %s", string(bodyBytes)))
		return
	}
	defer resp.Body.Close()

	var sessionResponse struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse sign-in response")
		return
	}

	sessionID := sessionResponse.ID

	// Now, create a session token (JWT)
	tokenResp, err := clerkAPIRequest("POST", fmt.Sprintf("/sessions/%s/tokens", sessionID), nil)
	if err != nil || tokenResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to create session token: %s", string(bodyBytes)))
		return
	}
	defer tokenResp.Body.Close()

	var tokenResponse struct {
		Object string `json:"object"`
		JWT    string `json:"jwt"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse session token response")
		return
	}

	sessionToken := tokenResponse.JWT

	// Set session cookie
	c.SetCookie(clerkSessionCookieName, sessionToken, 3600, "/", "", false, true)
	c.String(http.StatusOK, "Sign-in successful")
}

func SignInHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Fetch user details using the email address
	userListResp, err := clerkAPIRequest("GET", fmt.Sprintf("/users?email_address=%s", email), nil)
	if err != nil || userListResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(userListResp.Body)
		println("Failed to fetch user by email:", err)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to fetch user by email: %s", string(bodyBytes)))
		return
	}
	defer userListResp.Body.Close()

	var userListResponse []struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(userListResp.Body).Decode(&userListResponse); err != nil || len(userListResponse) == 0 {
		fmt.Printf("Failed to parse user list response: %v", err)
		c.String(http.StatusUnauthorized, "User not found")
		return
	}

	userID := userListResponse[0].ID
	println("User ID:", userID)

	// Verify the user's password
	verifyPasswordPayload := map[string]string{
		"password": password,
	}
	verifyResp, err := clerkAPIRequest("POST", fmt.Sprintf("/users/%s/verify_password", userID), verifyPasswordPayload)
	if err != nil || verifyResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(verifyResp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Password verification failed: %s", string(bodyBytes)))
		return
	}
	defer verifyResp.Body.Close()

	// Use user_id to create a session and get a session token
	signIn(c, userID)
}

func SignOutHandler(c *gin.Context) {
	// Retrieve the session token from the cookie
	sessionToken, err := GetSessionToken(c)
	if err != nil || sessionToken == "" {
		c.String(http.StatusUnauthorized, "No active session found")
		return
	}

	// Verify the session token to get the claims
	claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
		Token: sessionToken,
	})
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid session token")
		return
	}

	sessionID := claims.Claims.SessionID

	// Call Clerk API to revoke the session
	resp, err := clerkAPIRequest("POST", fmt.Sprintf("/sessions/%s/revoke", sessionID), nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to sign out: %s", string(bodyBytes)))
		return
	}
	defer resp.Body.Close()

	// Clear the cookie
	c.SetCookie(clerkSessionCookieName, "", -1, "/", "", false, true)
	c.String(http.StatusOK, "Sign-out successful")
}

func ProtectedRoute(c *gin.Context) {
	sessionToken, err := GetSessionToken(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify the session token
	claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
		Token: sessionToken,
	})
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid session token")
		return
	}

	userID := claims.Subject
	c.String(http.StatusOK, fmt.Sprintf("Your user ID is: %s", userID))
}
