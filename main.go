package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/getNeonConnection"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/testNeonDatabase"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/mainPage"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goConstants"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goEnv"

	"github.com/gin-gonic/gin"
)

/*func SignUpUser(email, password string) error {
	// Create a context
	ctx := context.Background()

	// Create a new user with the provided email and password
	newUser, err := user.Create(ctx, &user.CreateParams{
		EmailAddresses: &[]string{email},
		Password:       clerk.String(password),
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	fmt.Printf("User created with ID: %s\n", newUser.ID)
	return nil
}*/

func main() {

	err := goEnv.GetEnvVar()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

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
		cookie, err := c.Request.Cookie("clerkSession")
		if err != nil {
			fmt.Println("No cookie found")
		} else {
			fmt.Printf("Cookie value: %s\n", cookie.Value)
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

func ClerkSessionMiddleware(c *gin.Context) {
	// Check if the session cookie exists
	sessionToken, err := c.Cookie("clerk_session")
	if err != nil || sessionToken == "" {
		// Proceed without session if not found
		fmt.Println("No session token found in cookies")
		c.Next()
		return
	}

	// Store the session token in the context for later retrieval
	c.Set("clerk_session", sessionToken)
	c.Next()
}

func GetSessionToken(c *gin.Context) (string, error) {
	// Attempt to retrieve the session token from the context
	token, exists := c.Get("clerk_session")
	if !exists {
		return "", fmt.Errorf("session token not found in context")
	}

	// Assert that the token is a string
	sessionToken, ok := token.(string)
	if !ok {
		return "", fmt.Errorf("session token is not a valid string")
	}

	return sessionToken, nil
}

const clerkAPIBaseURL = "https://api.clerk.com/v1"

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
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Reattach body for further use
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
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Sign up failed: %s", string(bodyBytes)))
		return
	}

	// Parse the user_id from the response
	var signupResponse struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&signupResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse sign-up response")
		return
	}

	// Automatically sign in the user
	signIn(c, signupResponse.UserID)
}

func signIn(c *gin.Context, userID string) {
	payload := map[string]string{
		"user_id": userID,
	}

	resp, err := clerkAPIRequest("POST", "/sessions", payload)
	if err != nil || resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to sign in: %s", string(bodyBytes)))
		return
	}
	defer resp.Body.Close()

	var sessionResponse struct {
		SessionToken string `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse sign-in response")
		return
	}

	// Set session cookie
	c.SetCookie("clerk_session", sessionResponse.SessionToken, 3600, "/", "", false, true)
	c.String(http.StatusOK, "Sign-in successful")
}

func SignInHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Fetch user_id using email and password
	payload := map[string]string{
		"email_address": email,
		"password":      password,
	}

	resp, err := clerkAPIRequest("POST", "/users/password", payload)
	if err != nil || resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to fetch user: %s", string(bodyBytes)))
		return
	}
	defer resp.Body.Close()

	var userResponse struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse user response")
		return
	}

	// Use user_id to create a session
	signIn(c, userResponse.UserID)
}

func SignOutHandler(c *gin.Context) {
	// Retrieve the session token from the context
	sessionToken, err := GetSessionToken(c)
	if err != nil || sessionToken == "" {
		c.String(http.StatusUnauthorized, "No active session found")
		return
	}

	// Call Clerk API to revoke the session
	resp, err := clerkAPIRequest("POST", "/sessions/"+sessionToken+"/revoke", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to sign out: %s", string(bodyBytes)))
		return
	}
	defer resp.Body.Close()

	// Clear the cookie
	c.SetCookie("clerk_session", "", -1, "/", "", false, true)
	c.String(http.StatusOK, "Sign-out successful")
}

func ProtectedRoute(c *gin.Context) {
	sessionToken, err := GetSessionToken(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Your session token is: %s", sessionToken))
}
