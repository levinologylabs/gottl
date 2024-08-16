package handlers

import (
	"net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web/extractors"
)

type UserController struct {
	userservice *services.UserService
	pwservice   *services.PasswordService
}

func NewAuthController(userservice *services.UserService, pwservice *services.PasswordService) *UserController {
	return &UserController{
		userservice: userservice,
		pwservice:   pwservice,
	}
}

// Register godoc
//
//	@Tags			User
//	@Summary		Register a new user
//	@Description	Register a new user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.UserRegister	true	"User details"
//	@Success		201		{object}	dtos.User
//	@Router			/api/v1/users [POST]
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.UserRegister](r)
	if err != nil {
		return err
	}

	user, err := uc.userservice.Register(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusCreated, user)
}

// Authenticate godoc
//
//	@Tags			Auth
//	@Summary		Authenticate a user
//	@Description	Authenticate a user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.UserAuthenticate	true	"User credentials"
//	@Success		200		{object}	dtos.UserSession
//	@Router			/api/v1/users/authenticate [POST]
func (uc *UserController) Authenticate(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.UserAuthenticate](r)
	if err != nil {
		return err
	}

	session, err := uc.userservice.Authenticate(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, session)
}

// Self godoc
//
//	@Tags			User
//	@Summary		Get the current user
//	@Description	Get the current user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.User
//	@Router			/api/v1/users/self [GET]
//	@Security		Bearer
func (uc *UserController) Self(w http.ResponseWriter, r *http.Request) error {
	// Note: Using the context user is _fine_ for now, but if we switch to a JWT
	// and embed the user in the token, we'll need to update this to fetch the
	// user from the session token.
	user := services.UserFrom(r.Context())
	return server.JSON(w, http.StatusOK, user)
}

// Update godoc
//
//	@Tags			User
//	@Summary		Update the current user's details
//	@Description	Update the current user's details
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.UserUpdate	true	"User details"
//	@Success		200		{object}	dtos.User
//	@Router			/api/v1/users/self [PATCH]
//	@Security		Bearer
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.UserUpdate](r)
	if err != nil {
		return err
	}

	user := services.UserFrom(r.Context())

	updated, err := uc.userservice.UpdateDetails(r.Context(), user.ID, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, updated)
}

// ResetPasswordRequest godoc
//
//	@Tags			User
//	@Summary		Request a password reset
//	@Description	Request a password reset
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.PasswordResetRequest	true	"User details"
//	@Success		204		{object} nil
//	@Router			/api/v1/users/request-password-reset [POST]
func (uc *UserController) ResetPasswordRequest(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.PasswordResetRequest](r)
	if err != nil {
		return err
	}

	_, err = uc.pwservice.RequestReset(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusNoContent, nil)
}

// ResetPassword godoc
//
//	@Tags			User
//	@Summary		Reset a user's password
//	@Description	Reset a user's password
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.PasswordReset	true	"User details"
//	@Success		204		{object} nil
//	@Router			/api/v1/users/reset-password [POST]
func (uc *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.PasswordReset](r)
	if err != nil {
		return err
	}

	err = uc.pwservice.Reset(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusNoContent, nil)
}
