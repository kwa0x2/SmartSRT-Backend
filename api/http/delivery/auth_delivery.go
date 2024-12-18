package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"sync"
)

type AuthDelivery struct {
	Env            *bootstrap.Env
	UserUseCase    domain.UserUseCase
	SessionUseCase domain.SessionUseCase
}

var (
	stateStore = sync.Map{}
)

func (ad *AuthDelivery) GoogleSignIn(ctx *gin.Context) {
	googleConfig := bootstrap.GoogleConfig(ad.Env)
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := googleConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ad *AuthDelivery) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if _, exists := stateStore.Load(state); !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid state parameter. Please try again"))
		return
	}

	googleConfig := bootstrap.GoogleConfig(ad.Env)

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Code-Token Exchange Failed"))
		return
	}

	resp, respErr := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if respErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("User data fetch failed"))
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("JSON Parsing Failed"))
		return
	}

	newUser := &domain.User{
		Name:      userData["name"].(string),
		Email:     userData["email"].(string),
		AvatarURL: userData["picture"].(string),
	}

	if err = ad.UserUseCase.Create(newUser); err != nil && !mongo.IsDuplicateKeyError(err) {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession()
	if sessionErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(sessionErr.Error()))
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "",
	})

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (ad *AuthDelivery) GitHubSignIn(ctx *gin.Context) {
	githubConfig := bootstrap.GitHubConfig(ad.Env)
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := githubConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ad *AuthDelivery) GitHubCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if _, exists := stateStore.Load(state); !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid state parameter. Please try again"))
		return
	}

	githubConfig := bootstrap.GitHubConfig(ad.Env)

	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Code-Token Exchange Failed"))
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(token.AccessToken).
		Get("https://api.github.com/user")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send request to GitHub"))
		return
	}

	var userData map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &userData); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to decode user data"))
		return
	}

	emailResp, err := client.R().
		SetAuthToken(token.AccessToken).
		Get("https://api.github.com/user/emails")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch emails"))
		return
	}

	var emails []map[string]interface{}
	if err = json.Unmarshal(emailResp.Body(), &emails); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to decode email data"))
		return
	}

	var email string
	if len(emails) > 0 {
		email = emails[0]["email"].(string)
	} else {
		email = "Email not available"
	}

	newUser := &domain.User{
		Name:      userData["name"].(string),
		Email:     email,
		AvatarURL: userData["avatar_url"].(string),
	}

	if err = ad.UserUseCase.Create(newUser); err != nil && !mongo.IsDuplicateKeyError(err) {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("user create error"))
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession()
	if sessionErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(sessionErr.Error()))
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "",
	})

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (ad *AuthDelivery) CredentialsSignUp(ctx *gin.Context) {
	var body domain.CredentialsSignUpBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Body Parsing Failed"))
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to hash password"))
		return
	}

	newUser := &domain.User{
		Name:      body.Name,
		Email:     body.Email,
		Password:  hashedPassword,
		AvatarURL: body.AvatarURL,
	}

	if err = ad.UserUseCase.Create(newUser); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusFound, utils.NewMessageResponse("already exists"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("User created successfully"))
}

func (ad *AuthDelivery) CredentialsSignIn(ctx *gin.Context) {
	var body domain.CredentialsSignInBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Body Parsing Failed"))
		return
	}

	user, err := ad.UserUseCase.FindOneByEmail(body.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, utils.NewMessageResponse("Email not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch user"))
		return
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Invalid credentials"))
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession()
	if sessionErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(sessionErr.Error()))
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "",
	})

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("User sign in successfully"))
}

func (ad *AuthDelivery) SignOut(ctx *gin.Context) {
	cookie, err := ctx.Cookie("sid")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Session ID not found in cookie"))
		return
	}

	if err = ad.SessionUseCase.DeleteSession(cookie); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to delete session"))
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "",
	})

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("User signed out successfully"))
}
