//go:build unit

package models

import (
	"github.com/RafalSalwa/interview-app-srv/pkg/jwt"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type usersModels struct {
	HTTPResponse UserResponse
	HTTPRequest  UserRequest
	DBResponse   UserDBResponse
	DBModel      UserDBModel
}

func initStructs() usersModels {
	uid := int64(1)
	username := "username"
	email := username + "@interview.com"
	firstname := "firstname"
	lastname := "lastname"
	password := "password"
	phoneNo := "123123123"
	VerificationCode := "abcdefgh"
	timeNow := time.Now()
	uresp := UserResponse{
		Id:               uid,
		Username:         username,
		Firstname:        firstname,
		Verified:         true,
		VerificationCode: "",
		Active:           true,
		AccessToken:      "",
		RefreshToken:     "",
		CreatedAt:        timeNow,
		UpdatedAt:        &timeNow,
		LastLogin:        &timeNow,
		DeletedAt:        &timeNow,
	}
	ureq := UserRequest{
		Id:       uid,
		Username: username,
		Email:    email,
	}
	udbresp := UserDBResponse{
		Id:        uid,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
		Verified:  true,
		Active:    true,
		CreatedAt: timeNow,
		LastLogin: timeNow,
	}
	udb := UserDBModel{
		Id:               uid,
		Username:         username,
		Password:         password,
		Firstname:        firstname,
		Lastname:         lastname,
		Email:            email,
		Phone:            phoneNo,
		VerificationCode: VerificationCode,
		Verified:         true,
		Active:           true,
		CreatedAt:        timeNow,
		UpdatedAt:        &timeNow,
		LastLogin:        &timeNow,
		DeletedAt:        &timeNow,
	}
	return usersModels{
		HTTPResponse: uresp,
		HTTPRequest:  ureq,
		DBResponse:   udbresp,
		DBModel:      udb,
	}
}

func TestUserResponseMappers(t *testing.T) {
	testModel := initStructs()
	ur := UserResponse{}
	err := ur.FromDBModel(&testModel.DBModel)
	assert.NoError(t, err)
	ur = UserResponse{}
	err = ur.FromDBResponse(&testModel.DBResponse)
	assert.NoError(t, err)
	tp := &jwt.TokenPair{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	ur.AssignTokenPair(tp)
	assert.NoError(t, err)
	assert.Equal(t, ur.AccessToken, "access")
	assert.Equal(t, ur.RefreshToken, "refresh")

	ur = UserResponse{}
	ur.FromProtoSignIn(&intrvproto.SignInUserResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
	})
	assert.NoError(t, err)
	assert.Equal(t, ur.AccessToken, "access")
	assert.Equal(t, ur.RefreshToken, "refresh")
	err = ur.FromProtoSignUp(&intrvproto.SignUpUserResponse{
		Id:                1,
		Username:          "username",
		VerificationToken: "abcd",
		CreatedAt:         nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, ur.Username, "username")
	assert.Equal(t, ur.VerificationCode, "abcd")

	dbu := UserDBModel{
		Id: 1,
	}
	sup := SignUpUserRequest{
		Email:           "test@test.com",
		Password:        "password",
		PasswordConfirm: "password",
	}
	err = dbu.FromCreateUserReq(sup, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, dbu)
	intrvU := &intrvproto.User{
		Id:        1,
		Username:  "username",
		Email:     "email",
		CreatedAt: nil,
	}
	err = ur.FromProtoUserResponse(&intrvproto.UserResponse{User: intrvU})
	assert.NoError(t, err)
	assert.Equal(t, ur.Username, "username")

	ur = UserResponse{Id: 1}
	ur.FromProtoUserDetails(&intrvproto.UserDetails{Id: 1})
	assert.NotEmpty(t, ur)
}
