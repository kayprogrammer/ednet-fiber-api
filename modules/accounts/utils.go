package accounts

import (
	"context"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/idtoken"
)

type GooglePayload struct {
	SUB           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func ConvertGoogleToken(ctx context.Context, accessToken string) (*GooglePayload, *config.ErrorResponse) {
	cfg := config.GetConfig()

	payload, err := idtoken.Validate(ctx, accessToken, cfg.GoogleClientID)
	if err != nil {
		errMsg := "Invalid Token"
		if strings.Contains(err.Error(), "audience provided") {
			errMsg = "Invalid Audience"
		}
		errData := config.RequestErr(config.ERR_INVALID_TOKEN, errMsg)
		return nil, &errData
	}

	// Bind JSON into struct
	data := GooglePayload{}
	mapstructure.Decode(payload.Claims, &data)
	return &data, nil
}

func RegisterSocialUser(db *ent.Client, ctx context.Context, email string, name string, avatar *string) (*string, *string, *config.ErrorResponse) {
	cfg := config.GetConfig()

	socialUser, _ := db.User.Query().Where(user.Email(email)).Only(ctx)
	if (socialUser == nil) {
		password := config.HashPassword(cfg.SocialsPassword)
		username := GenerateUsernameFromEmail(db, ctx, email, nil)
		socialUser = db.User.Create().
			SetName(name).
			SetEmail(email).
			SetPassword(password).
			SetUsername(username).	
			SetSocialLogin(true).	
			SetIsVerified(true).
			SaveX(ctx)
	} else {
		if !socialUser.SocialLogin {
			errData := config.RequestErr(config.ERR_INVALID_AUTH, "Requires password to login")
			return nil, nil, &errData
		}
	}
	// Generate tokens
	access := GenerateAccessToken(socialUser.ID, socialUser.Username)
	refresh := GenerateRefreshToken()
	userManager.AddTokens(db, ctx, socialUser, access, refresh)
	return &access, &refresh, nil
}

func GenerateUsernameFromEmail(db *ent.Client, ctx context.Context, email string, username *string) string {
	emailSubstr := strings.Split(email, "@")[0]

	uniqueUsername := slug.Make(emailSubstr)
	if username != nil {
		uniqueUsername = *username
	}

	// Check for uniqueness and adjust if necessary
	for {
		exisitngUser := userManager.GetByUsername(db, ctx, uniqueUsername)
		if exisitngUser == nil {
			// Username is unique
			break
		}
		// Append a random string to make it unique
		randomStr := strconv.FormatUint(uint64(config.GetRandomInt(7)), 10)
		uniqueUsername = slug.Make(emailSubstr) + randomStr
	}
	return uniqueUsername
}