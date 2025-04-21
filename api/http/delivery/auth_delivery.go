package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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
)

type AuthDelivery struct {
	Env            *bootstrap.Env
	UserUseCase    domain.UserUseCase
	SessionUseCase domain.SessionUseCase
	SinchUseCase   domain.SinchUseCase
	ResendUseCase  domain.ResendUseCase
}

var (
	stateStore = sync.Map{}
)

func (ad *AuthDelivery) GoogleLogin(ctx *gin.Context) {
	googleConfig := bootstrap.GoogleConfig(ad.Env)
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := googleConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ad *AuthDelivery) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	locale := ctx.GetString("locale")
	loginRedirect := fmt.Sprintf("%s/%s/auth/login", ad.Env.FrontEndURL, locale)
	path := fmt.Sprintf("%s/auth/login", locale)

	if _, exists := stateStore.Load(state); !exists {
		utils.SetErrorCookie(ctx, "invalid_state", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	googleConfig := bootstrap.GoogleConfig(ad.Env)

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	client := resty.New()
	resp, respErr := client.R().
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if respErr != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	var userData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &userData)
	if err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	user, err := ad.UserUseCase.FindOneByEmail(userData["email"].(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			jwtClaims := jwt.MapClaims{
				"name":       userData["name"].(string),
				"email":      userData["email"].(string),
				"avatar_url": userData["picture"].(string),
				"auth_type":  types.Google,
			}

			exp1HourUnix := time.Now().Add(1 * time.Hour).Unix() // 1 hour

			tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env, exp1HourUnix)
			if tokenErr != nil {
				utils.SetErrorCookie(ctx, "server_error", path)
				ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
				return
			}

			otpPath := fmt.Sprintf("/%s/auth/otp", ctx.GetString("locale"))
			utils.SetAuthTokenCookie(ctx, tokenString, otpPath, 3600) // 1 hour

			redirectURL := fmt.Sprintf("%s/%s/auth/otp", ad.Env.FrontEndURL, locale)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		} else {
			utils.SetErrorCookie(ctx, "server_error", path)
			ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
			return
		}
	}

	if user.AuthType != types.Google {
		errorType := fmt.Sprintf("exists_%s", user.AuthType)
		utils.SetErrorCookie(ctx, errorType, path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSessionAndUpdateLastLogin(user.ID, user.Role, user.Email)
	if sessionErr != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	utils.SetSIDCookie(ctx, sessionID)

	redirectURL := fmt.Sprintf("%s/%s/auth/verify", ad.Env.FrontEndURL, locale)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (ad *AuthDelivery) GitHubLogin(ctx *gin.Context) {
	githubConfig := bootstrap.GitHubConfig(ad.Env)
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := githubConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ad *AuthDelivery) GitHubCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	locale := ctx.GetString("locale")
	loginRedirect := fmt.Sprintf("%s/%s/auth/login", ad.Env.FrontEndURL, locale)
	path := fmt.Sprintf("%s/auth/login", locale)

	if _, exists := stateStore.Load(state); !exists {
		utils.SetErrorCookie(ctx, "invalid_state", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	githubConfig := bootstrap.GitHubConfig(ad.Env)

	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(token.AccessToken).
		Get("https://api.github.com/user")

	if err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	var userData map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &userData); err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	emailResp, err := client.R().
		SetAuthToken(token.AccessToken).
		Get("https://api.github.com/user/emails")

	if err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	var emails []map[string]interface{}
	if err = json.Unmarshal(emailResp.Body(), &emails); err != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
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
				"auth_type":  types.Github,
			}

			exp1HourUnix := time.Now().Add(1 * time.Hour).Unix() // 1 hour

			tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env, exp1HourUnix)
			if tokenErr != nil {
				utils.SetErrorCookie(ctx, "server_error", path)
				ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
				return
			}

			otpPath := fmt.Sprintf("/%s/auth/otp", ctx.GetString("locale"))
			utils.SetAuthTokenCookie(ctx, tokenString, otpPath, 3600) // 1 hour

			redirectURL := fmt.Sprintf("%s/%s/auth/otp", ad.Env.FrontEndURL, locale)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		} else {
			utils.SetErrorCookie(ctx, "server_error", path)
			ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
			return
		}
	}

	if user.AuthType != types.Github {
		errorType := fmt.Sprintf("exists_%s", user.AuthType)
		utils.SetErrorCookie(ctx, errorType, path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSessionAndUpdateLastLogin(user.ID, user.Role, user.Email)
	if sessionErr != nil {
		utils.SetErrorCookie(ctx, "server_error", path)
		ctx.Redirect(http.StatusTemporaryRedirect, loginRedirect)
		return
	}

	utils.SetSIDCookie(ctx, sessionID)

	redirectURL := fmt.Sprintf("%s/%s/auth/verify", ad.Env.FrontEndURL, locale)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (ad *AuthDelivery) CredentialsLogin(ctx *gin.Context) {
	var body domain.CredentialsLoginBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	user, err := ad.UserUseCase.FindOneByEmail(body.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, utils.NewMessageResponse("User not found. Please register to create an account."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	if user.AuthType != types.Credentials {
		ctx.JSON(http.StatusConflict, utils.NewMessageResponse(
			fmt.Sprintf("An account with this email already exists. Please log in using %s.", utils.ToCamelCase(string(user.AuthType))),
		))
		return
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Incorrect email or password. Please try again."))
		return
	}

	sessionID, sessionErr := ad.SessionUseCase.CreateSessionAndUpdateLastLogin(user.ID, user.Role, user.Email)
	if sessionErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to create a user session. Please try again later or contact support."))
		return
	}

	utils.SetSIDCookie(ctx, sessionID)

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (ad *AuthDelivery) Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("sid")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support2."))
		return
	}

	if err = ad.SessionUseCase.DeleteSession(cookie); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	utils.DeleteCookie(ctx, "sid")

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("User logout successfully"))
}

func (ad *AuthDelivery) SinchSendOTP(ctx *gin.Context) {
	var req domain.PhoneNumberBody

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please provide a valid phone number."))
		return
	}

	if err := ad.SinchUseCase.SendOTP(req.PhoneNumber); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send OTP. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("OTP has been successfully sent to your phone number."))
}

func (ad *AuthDelivery) SendSetupNewPasswordEmail(ctx *gin.Context) {
	var body domain.EmailBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request. Please provide a valid email address."))
		return
	}

	user, err := ad.UserUseCase.FindOneByEmailAndAuthType(body.Email, types.Credentials)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, utils.NewMessageResponse("The email address is not associated with any credentials account. Please check and try again."))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
			return
		}
	}

	jwtClaims := jwt.MapClaims{
		"process": "update_password",
		"id":      user.ID,
	}

	exp3MinUnix := time.Now().Add(5 * time.Minute).Unix() // 5 min

	tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env, exp3MinUnix)
	if tokenErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	path := fmt.Sprintf("/%s/auth/reset-password", ctx.GetString("locale"))

	utils.SetAuthTokenCookie(ctx, tokenString, path, 300) // 5 min

	setupPasswordLink := fmt.Sprintf("%s/%s/auth/reset-password", ad.Env.FrontEndURL, ctx.GetString("locale"))

	_, sentErr := ad.ResendUseCase.SendSetupPasswordEmail(body.Email, setupPasswordLink)
	if sentErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send new password setup email. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("New password setup email sent successfully. Please check your inbox."))
}

func (ad *AuthDelivery) UpdatePassword(ctx *gin.Context) {
	var body domain.PasswordBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request. Please provide a valid password."))
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	processStr, ok := jwtClaims["process"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	if processStr != "update_password" {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userIDStr, ok := jwtClaims["id"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	if updateErr := ad.UserUseCase.UpdateCredentialsPasswordByID(userID, hashedPassword); updateErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to update password. Please try again later or contact support."))
		return
	}

	utils.DeleteCookie(ctx, "token")

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("Password updated successfully."))
}

func (ad *AuthDelivery) VerifyOTPAndCreate(ctx *gin.Context) {
	var body domain.VerifyOTPAndCreateBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	valid, err := ad.SinchUseCase.VerifyOTP(body.PhoneNumber, body.OTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to verify OTP. Please try again later or contact support."))
		return
	} else if !valid {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Incorrect OTP. Please check and try again."))
		return
	}

	newUser := &domain.User{
		Name:        body.Name,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		AvatarURL:   body.AvatarURL,
		AuthType:    body.AuthType,
	}

	if body.Password != "" && newUser.AuthType == types.Credentials {
		hashedPassword, hashErr := utils.HashPassword(body.Password)
		if hashErr != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
			return
		}

		newUser.Password = hashedPassword
	}

	if err = ad.UserUseCase.Create(newUser); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, utils.NewMessageResponse("A user with this information already exists. Please try a different email or phone number."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		ctx.JSON(http.StatusInternalServerError, err.Error())

		return
	}

	utils.DeleteCookie(ctx, "token")

	ctx.JSON(http.StatusCreated, utils.NewMessageResponse("User created successfully"))
}

func (ad *AuthDelivery) SendDeleteAccountMail(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	userData := user.(*domain.User)

	jwtClaims := jwt.MapClaims{
		"process": "delete_account",
		"id":      userData.ID,
		"email":   userData.Email,
		"image":   userData.AvatarURL,
	}

	exp3MinUnix := time.Now().Add(5 * time.Minute).Unix() // 5 min

	tokenString, tokenErr := utils.GenerateJWT(jwtClaims, ad.Env, exp3MinUnix)
	if tokenErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	path := fmt.Sprintf("/%s/auth/account/delete", ctx.GetString("locale"))
	utils.SetAuthTokenCookie(ctx, tokenString, path, 300)

	deleteAccountLink := fmt.Sprintf("%s/%s/auth/account/delete", ad.Env.FrontEndURL, ctx.GetString("locale"))

	_, sentErr := ad.ResendUseCase.SendDeleteAccountEmail(userData.Email, deleteAccountLink)
	if sentErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send delete account email. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("Delete account email sent successfully. Please check your inbox."))
}

func (ad *AuthDelivery) DeleteAccount(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	processStr, ok := jwtClaims["process"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	if processStr != "delete_account" {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userIDStr, ok := jwtClaims["id"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	sessionID, err := ctx.Cookie("sid")
	if err == nil {
		if err = ad.SessionUseCase.DeleteSession(sessionID); err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
			return
		}
		utils.DeleteCookie(ctx, "sid")
	}

	if err = ad.UserUseCase.DeleteUser(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	utils.DeleteCookie(ctx, "token")

	ctx.JSON(http.StatusNoContent, utils.NewMessageResponse("Account deleted successfully!"))
}
