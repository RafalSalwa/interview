package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RafalSalwa/auth-api/pkg/encdec"
	"github.com/RafalSalwa/auth-api/pkg/hashing"
	"github.com/RafalSalwa/auth-api/pkg/jwt"
	intrvproto "github.com/RafalSalwa/auth-api/proto/grpc"
	"github.com/jinzhu/copier"
)

func (r *UserResponse) FromDBResponse(user *UserDBResponse) error {
	err := copier.Copy(r, &user)
	if err != nil {
		return fmt.Errorf("from db to response error: %w", err)
	}
	return nil
}

func (m *UserDBModel) FromCreateUserReq(cur SignUpUserRequest, enc ...bool) error {
	err := copier.Copy(m, &cur)
	if err != nil {
		return fmt.Errorf("from create to db model error: %w", err)
	}
	if len(enc) > 0 && enc[0] == true {
		m.Email = encdec.Encrypt(cur.Email)
		m.Password = hashing.Argon2ID(cur.Password)
	}
	return nil
}

func (r *UserResponse) FromProtoUserDetails(pbu *intrvproto.UserDetails) {
	r.Id = pbu.GetId()
	r.Username = pbu.GetUsername()
	r.Firstname = pbu.GetFirstname()
	r.Lastname = pbu.GetLastname()
	r.Verified = pbu.GetVerified()
	r.VerificationCode = pbu.GetVerificationCode()
	r.Active = pbu.GetActive()
	r.CreatedAt = pbu.GetCreatedAt().AsTime()

	ll := pbu.GetLastLogin().AsTime()
	r.LastLogin = &ll
	if pbu.GetEmail() != "" {
		dec, _ := encdec.Decrypt(pbu.GetEmail())
		r.Email = dec
	}
}

func (r *UserResponse) FromDBModel(um *UserDBModel) error {
	err := copier.Copy(r, &um)
	if err != nil {
		return fmt.Errorf("from response to db error: %w", err)
	}
	r.Username = um.Username
	r.CreatedAt = um.CreatedAt
	return nil
}

func (r *UserResponse) AssignTokenPair(tp *jwt.TokenPair) {
	r.AccessToken = tp.AccessToken
	r.RefreshToken = tp.RefreshToken
}

func (r *UserResponse) FromProtoSignIn(pbu *intrvproto.SignInUserResponse) {
	r.AccessToken = pbu.GetAccessToken()
	r.RefreshToken = pbu.GetRefreshToken()
}
func (r *UserResponse) FromProtoSignUp(pbu *intrvproto.SignUpUserResponse) error {
	r.Username = pbu.GetUsername()
	r.VerificationCode = pbu.GetVerificationToken()
	c := pbu.GetCreatedAt().AsTime()
	r.CreatedAt = c
	return nil
}

func (r *UserResponse) FromProtoUserResponse(pu *intrvproto.UserResponse) error {
	err := copier.Copy(r, &pu)
	r.Id = pu.GetUser().GetId()
	r.Username = pu.GetUser().GetUsername()

	if err != nil {
		return fmt.Errorf("from response to db error: %w", err)
	}
	return nil
}

func (r *UserDBResponse) FromProtoUserDetails(pw *intrvproto.UserDetails) error {
	err := copier.Copy(r, pw)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserMongoModel) FromDBModel(user *UserDBModel) error {
	err := copier.Copy(m, user)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserDBModel) FromMongoUser(um2 UserMongoModel) error {
	err := copier.Copy(um, um2)
	if err != nil {
		return err
	}
	return nil
}

func (m *StringInterfaceMap) Scan(src interface{}) error {
	var source []byte
	_m := make(map[string]interface{})

	switch src.(type) {
	case []uint8:
		source = []byte(src.([]uint8))
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = StringInterfaceMap(_m)
	return nil
}

func (m StringInterfaceMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value([]byte(j)), nil
}
