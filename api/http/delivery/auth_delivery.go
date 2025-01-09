package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"sync"
)

type AuthDelivery struct {
	Env            *bootstrap.Env
	UserUseCase    domain.UserUseCase
	SessionUseCase domain.SessionUseCase
	SinchUseCase   domain.SinchUseCase
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

	user, err := ad.UserUseCase.FindOneByEmail(userData["email"].(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			jwtClaims := jwt.MapClaims{
				"name":       userData["name"].(string),
				"email":      userData["email"].(string),
				"avatar_url": userData["picture"].(string),
				"auth_with":  types.Google,
			}

			tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env)
			if tokenErr != nil {
				ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("JWT Token Failed"))
				return
			}

			redirectURL := fmt.Sprintf("%s/en/auth/otp?auth=%s", ad.Env.FrontEndURL, tokenString)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("User data fetch failed"))
			return
		}
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession(user.ID)
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

	redirectURL := fmt.Sprintf("%s/en/auth/verify", ad.Env.FrontEndURL)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
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

	user, err := ad.UserUseCase.FindOneByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			jwtClaims := jwt.MapClaims{
				"name":       userData["name"].(string),
				"email":      email,
				"avatar_url": userData["avatar_url"].(string),
				"auth_with":  types.Github,
			}

			tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env)
			if tokenErr != nil {
				ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("JWT Token Failed"))
				return
			}

			redirectURL := fmt.Sprintf("%s/en/auth/otp?auth=%s", ad.Env.FrontEndURL, tokenString)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("User data fetch failed"))
			return
		}
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession(user.ID)
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

	redirectURL := fmt.Sprintf("%s/en/auth/verify", ad.Env.FrontEndURL)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (ad *AuthDelivery) VerifyOTPAndCreate(ctx *gin.Context) {
	var body domain.VerifyOTPAndCreateBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	newUser := &domain.User{
		Name:        body.Name,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		AvatarURL:   body.AvatarURL,
		AuthWith:    body.AuthWith,
	}

	if body.Password != "" {
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to hash password"))
			return
		}

		newUser.Password = hashedPassword
	}

	valid, err := ad.SinchUseCase.VerifyOTP(newUser.PhoneNumber, body.OTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to verify OTP. Please try again later or contact support."))
		return
	} else if !valid {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Incorrect OTP. Please check and try again."))
		return
	}

	if err = ad.UserUseCase.Create(newUser); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, utils.NewMessageResponse("A user with this information already exists. Please try a different email or phone number."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An unexpected error occurred. Please try again later."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("User created successfully"))
}

func (ad *AuthDelivery) CredentialsSignIn(ctx *gin.Context) {
	var body domain.CredentialsSignInBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	user, err := ad.UserUseCase.FindOneByEmailAndAuthWith(body.Email, types.Credentials)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Incorrect email or password. Please try again."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch user. Please try again later or contact support."))
		return
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Incorrect email or password. Please try again."))
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSession(user.ID)
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

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
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

func (ad *AuthDelivery) SinchSendOTP(ctx *gin.Context) {
	var req domain.PhoneNumberBody

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	if err := ad.SinchUseCase.SendOTP(req.PhoneNumber); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send OTP. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("OTP has been successfully sent to your phone number."))
}

func (ad *AuthDelivery) Check(ctx *gin.Context) {
	sessionUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("User ID not found in session"))
		return
	}

	userIDStr, ok := sessionUserID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid user ID format"))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid user ID"))
		return
	}

	user, err := ad.UserUseCase.FindOneByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch user"))
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (ad *AuthDelivery) IsEmailExists(ctx *gin.Context) {
	var body domain.IsEmailExistsBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	exists, err := ad.UserUseCase.IsEmailExists(body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch user. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (ad *AuthDelivery) IsPhoneExists(ctx *gin.Context) {
	var body domain.PhoneNumberBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	exists, err := ad.UserUseCase.IsPhoneExists(body.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to fetch user. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"exists": exists})
}
