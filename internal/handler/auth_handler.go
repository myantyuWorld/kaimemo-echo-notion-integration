package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Login(c echo.Context) error
	Callback(c echo.Context) error
	CheckAuth(c echo.Context) error
	Logout(c echo.Context) error
}

type authHandler struct{}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func (a *authHandler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out"})
}

// FE側OnMountedで実行する
func (a *authHandler) CheckAuth(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "authenticated"})
}

// Callback implements AuthHandler.
func (a *authHandler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "code is missing"})
	}

	// アクセストークンを取得
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get token"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp TokenResponse
	json.Unmarshal(body, &tokenResp)

	// ユーザー情報取得
	req, _ := http.NewRequest("GET", profileURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	client := &http.Client{}
	profileResp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get profile"})
	}
	defer profileResp.Body.Close()

	profileBody, _ := io.ReadAll(profileResp.Body)
	var profile Profile
	json.Unmarshal(profileBody, &profile)

	// TODO : ここでDBにユーザー情報を保存（省略）

	// JWT作成
	setAuthCookie(c, profile.UserID)

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful"})
}

// Login implements AuthHandler.
func (a *authHandler) Login(c echo.Context) error {
	baseURL := "https://access.line.me/oauth2/v2.1/authorize"
	u, _ := url.Parse(baseURL)

	query := u.Query()
	query.Set("response_type", "code")
	query.Set("client_id", clientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("state", state)
	query.Set("scope", "profile openid email")

	u.RawQuery = query.Encode()

	return c.Redirect(302, u.String())
}

func setAuthCookie(c echo.Context, userID string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)
	return nil

}

var (
	clientID     = os.Getenv("LINE_CLIENT_ID")     // 環境変数から取得
	clientSecret = os.Getenv("LINE_CLIENT_SECRET") // 環境変数から取得
)

const (
	redirectURI = "https://your-backend.com/api/auth/callback"
	state       = "random_state_string" // CSRF対策用
)

const (
	tokenURL   = "https://api.line.me/oauth2/v2.1/token"
	profileURL = "https://api.line.me/v2/profile"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type Profile struct {
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	UserID      string `json:"userId"`
}

func NewAuthHandler() AuthHandler {
	return &authHandler{}
}
