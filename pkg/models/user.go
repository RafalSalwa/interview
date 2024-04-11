package models

import (
	"time"

	"gorm.io/gorm"
)

type StringInterfaceMap map[string]interface{}

type UserDBModel struct {
	Id               int64              `gorm:"id;primaryKey;autoIncrement"`
	Username         string             `gorm:"type:varchar(180);not null;uniqueIndex;not null"`
	Password         string             `gorm:"type:varchar(255);not null"`
	Firstname        string             `gorm:"column:first_name;type:varchar(255)"`
	Lastname         string             `gorm:"column:last_name;type:varchar(255)"`
	Email            string             `gorm:"type:varchar(255);not null;uniqueIndex;not null"`
	Phone            string             `gorm:"column:phone_no;type:varchar(11)"`
	VerificationCode string             `gorm:"column:verification_code;type:varchar(12)"`
	Roles            StringInterfaceMap `gorm:"column:roles;type:json"`
	CreatedAt        time.Time          `gorm:"column:created_at"`
	UpdatedAt        *time.Time         `gorm:"column:updated_at"`
	LastLogin        *time.Time         `gorm:"column:last_login"`
	DeletedAt        *time.Time         `gorm:"column:deleted_at"`
	Verified         bool               `gorm:"column:is_verified;default:false"`
	Active           bool               `gorm:"column:is_active;default:false"`
}

type UserMongoModel struct {
	Username         string     `bson:"username,omitempty"`
	Password         string     `bson:"password,omitempty"`
	Firstname        string     `bson:"firstname,omitempty"`
	Lastname         string     `bson:"lastname,omitempty"`
	Email            string     `bson:"email,omitempty"`
	Phone            string     `bson:"phone,omitempty"`
	VerificationCode string     `bson:"verification_code,omitempty"`
	Verified         bool       `bson:"is_verified,omitempty"`
	Active           bool       `bson:"is_active,omitempty"`
	Roles            Roles      `bson:"roles,omitempty"`
	CreatedAt        time.Time  `bson:"createdAt,omitempty"`
	UpdatedAt        *time.Time `bson:"updatedAt,omitempty"`
	LastLogin        *time.Time `bson:"lastLogin,omitempty"`
}

func (um *UserDBModel) BeforeCreate(tx *gorm.DB) (err error) {
	um.Active = false
	um.Verified = false
	return nil
}

func (UserDBModel) TableName() string {
	return "user"
}

type Users struct {
	Users []UserDBModel `json:"users"`
}

type UserRequest struct {
	Id               int64  `json:"id" govalidator:"int"`
	Email            string `json:"email,omitempty"`
	VerificationCode string `json:"verification_code,omitempty"`
	AccessToken      string `json:"token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
}

type SignUpUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8,max=32"`
}

type SignInUserRequest struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Password string `json:"password" validate:"required"`
}

type UserDBResponse struct {
	Id        int64
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Verified  bool
	Active    bool
	CreatedAt time.Time
	LastLogin time.Time
}

type UserResponse struct {
	Id               int64      `json:"id,omitempty"`
	Username         string     `json:"username,omitempty"`
	Firstname        string     `json:"firstname,omitempty"`
	Lastname         string     `json:"lastname,omitempty"`
	Email            string     `json:"email,omitempty"`
	Verified         bool       `json:"is_verified,omitempty"`
	VerificationCode string     `json:"verification_token,omitempty"`
	Active           bool       `json:"is_active,omitempty"`
	AccessToken      string     `json:"token,omitempty"`
	RefreshToken     string     `json:"refresh_token,omitempty"`
	CreatedAt        time.Time  `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type UpdateUserRequest struct {
	Id        int64   `json:"id,omitempty"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Password  *string `json:"password"`
}

type ChangePasswordRequest struct {
	Id              int64  `json:"id"`
	Email           string `json:"email"`
	OldPassword     string `json:"oldPassword"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type User struct {
	ID int64
	// Creates a primary key for UserID.
	UserID    int64  `gorm:"primary_key"`
	Username  string `sql:"type:VARCHAR(15);not null"`
	FName     string `sql:"size:100;not null" gorm:"column:FirstName"`
	LName     string `sql:"unique;unique_index;not null;DEFAULT:'Unknown'" gorm:"column:LastName"`
	Count     int    `gorm:"AUTO_INCREMENT"`
	TempField bool   `sql:"-"` // Ignore a Field
}

type Roles struct {
	Roles []string `json:"roles"`
}

type VerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}
