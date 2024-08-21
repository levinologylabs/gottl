package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/web/oauth/providers"
	"github.com/rs/zerolog"
)

type GoogleAuthController struct {
	store OauthSessionStore
	l     zerolog.Logger
	cfg   providers.GoogleConfig
}

func NewGoogleAuthController(l zerolog.Logger, cfg providers.GoogleConfig, store OauthSessionStore) *GoogleAuthController {
	return &GoogleAuthController{
		store: store,
		l:     l,
		cfg:   cfg,
	}
}

func (gac *GoogleAuthController) Authenticate(w http.ResponseWriter, r *http.Request) {
	statestring, err := gac.store.ProviderStateGet(r.Context())
	if err != nil {
		gac.l.Err(err).
			Msg("failed to get state string from provider")

		_ = server.Error().
			Status(http.StatusInternalServerError).
			Msg("internal server error").
			Write(r.Context(), w)
		return
	}

	u := gac.cfg.OathConfig().AuthCodeURL(statestring)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (gac *GoogleAuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	const UserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

	var (
		ctx   = r.Context()
		state = r.FormValue("state")
		code  = r.FormValue("code")
	)

	err := gac.store.ProviderStateUse(r.Context(), state)
	if err != nil {
		return err
	}

	token, err := gac.cfg.OathConfig().Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get(UserInfoURL + "?access_token=" + token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer func() { _ = response.Body.Close() }()

	user := providers.GoogleUser{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to decode user info: %s", err.Error())
	}

	tkn, err := gac.store.ProviderSession(ctx, "google", user.ID, user.Email, user.Name)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, tkn)
}
