package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RafalSalwa/interview-app-srv/pkg/generator"
	"github.com/RafalSalwa/interview-app-srv/pkg/hashing"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	mySql "github.com/RafalSalwa/interview-app-srv/pkg/sql"
)

type SqlServiceImpl struct {
	db     *mySql.DB
	logger *logger.Logger
}

type SQLService interface {
	GetByID(id int) (user *models.UserDBResponse, err error)
	GetByCode(code string) (user *models.UserDBModel, err error)
	UsernameInUse(user *models.SignUpUserRequest) bool
	StoreVerificationData(user *models.UserDBModel) bool
	UpdateUser(user *models.UpdateUserRequest) (err error)
	LoginUser(user *models.SignInUserRequest) (*models.UserResponse, error)
	UpdateUserPassword(user *models.UpdateUserRequest) (err error)
	CreateUser(user *models.SignUpUserRequest) (*models.UserResponse, error)
}

func NewMySQLService(db *mySql.DB, l *logger.Logger) *SqlServiceImpl {
	return &SqlServiceImpl{db, l}
}

func (s SqlServiceImpl) GetByCode(code string) (*models.UserDBModel, error) {
	user := &models.UserDBModel{}

	row := s.db.QueryRow("SELECT id,verification_code FROM `user` WHERE verification_code = ?", code)
	err := row.Scan(&user.Id,
		&user.VerificationCode)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SqlServiceImpl) GetByID(id int) (user *models.UserDBResponse, err error) {
	user = &models.UserDBResponse{}

	row := s.db.QueryRow("SELECT "+
		"id,"+
		"username,"+
		"first_name,"+
		"last_name, "+
		"password, "+
		"is_verified, "+
		"is_active, "+
		"created_at "+
		"FROM `user` WHERE is_active = 1 AND id=?", id)
	err = row.Scan(&user.Id,
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.Password,
		&user.Verified,
		&user.Active,
		&user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SqlServiceImpl) UsernameInUse(user *models.SignUpUserRequest) bool {
	dbUser := &models.UserDBResponse{}
	row := s.db.QueryRow("SELECT id FROM `user` WHERE email = ?", user.Email)
	err := row.Scan(&dbUser.Id)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

func (s *SqlServiceImpl) StoreVerificationData(user *models.UserDBModel) bool {
	_, err := s.db.Exec("UPDATE `user` SET is_verified = 1, is_active=1 WHERE id = ?", user.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

func (s *SqlServiceImpl) UpdateUser(user *models.UpdateUserRequest) (err error) {
	sqlStatement := "UPDATE user SET first_name=?, last_name=? WHERE id=?;"
	_, err = s.db.ExecContext(getContext(), sqlStatement, &user.Firstname, &user.Lastname, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqlServiceImpl) LoginUser(u *models.SignInUserRequest) (*models.UserResponse, error) {
	user := models.UserResponse{}

	row := s.db.QueryRow("SELECT id,username,first_name,last_name FROM `user` "+
		"WHERE (username=? OR email=?) AND (is_active = 1 AND is_verified = 1)", u.Username, u.Email)
	err := row.Scan(&user.Id,
		&user.Username,
		&user.Firstname)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// roles, err := phpserialize.Decode(user.RolesJson)

	if err != nil {
		return nil, err
	}

	// v, ok := roles.(map[interface{}]interface{})
	// if ok {
	//	for _, s := range v {
	//		user.Roles = append(user.Roles, fmt.Sprintf("%v", s))
	//	}
	//}

	return &user, nil
}

func (s *SqlServiceImpl) UpdateUserPassword(user *models.UpdateUserRequest) (err error) {
	sqlStatement := "UPDATE `user` SET password=? WHERE id=?;"
	_, err = s.db.ExecContext(getContext(), sqlStatement, user.Password, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqlServiceImpl) CreateUser(newUserRequest *models.SignUpUserRequest) (*models.UserResponse, error) {
	if err := hashing.Validate(newUserRequest.Password, newUserRequest.PasswordConfirm); err != nil {
		return nil, err
	}

	if s.UsernameInUse(newUserRequest) {
		return nil, errors.New("create user: username already in use")
	}

	vcode, err := generator.RandomString(6)
	if err != nil {
		return nil, err
	}

	newUserRequest.Password, err = hashing.HashPassword(newUserRequest.Password)
	if err != nil {
		return nil, err
	}

	dbUser := &models.UserDBModel{
		Password:         newUserRequest.Password,
		Email:            newUserRequest.Email,
		VerificationCode: vcode,
	}
	ctx := getContext()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	sqlStatement := "INSERT INTO `user` ( `password`, `email`, `verification_code`, `is_verified`,`is_active`) VALUES (?,?,?,0,1);"
	rows, err := tx.ExecContext(ctx,
		sqlStatement,
		dbUser.Password,
		dbUser.Email,
		dbUser.VerificationCode)

	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return nil, errTx
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	ur := &models.UserResponse{
		Id:       id,
		Username: dbUser.Username,
	}

	return ur, nil
}

func getContext() context.Context {
	ctx := context.Background()
	return ctx
}
