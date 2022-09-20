package db

import (
	"encoding/hex"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"habiko-go/pkg/utils"
	"strings"
	"time"

	"habiko-go/pkg/models"
)

type UserRepository struct {
	*TransactionProvider
}

func (r *UserRepository) IsEmailConfirmed(user models.User) bool {
	return user.EmailConfirmationToken == nil || *user.EmailConfirmationToken == ""
}

func (r *UserRepository) FindByID(transaction models.Transaction, id uuid.UUID) (*models.User, error) {
	var user models.User
	return &user, dbConn(transaction, r.DBConn).First(&user, "id = ?", id).Error
}

func (r *UserRepository) FindByEmail(transaction models.Transaction, email string) (*models.User, error) {
	var user models.User
	return &user, dbConn(transaction, r.DBConn).First(&user, "email = ?", email).Error
}

func (r *UserRepository) FindByCredentials(
	transaction models.Transaction, email string, password string,
) (*models.User, error) {
	var user models.User

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		return nil, err
	}

	db := dbConn(transaction, r.DBConn).First(&user, "email = ?", email)
	if db.Error != nil {
		return nil, db.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindPublicDataByName(transaction models.Transaction, name string) (*models.PublicUserData, error) {
	var user models.PublicUserData
	return &user, dbConn(transaction, r.DBConn).Raw(`
		SELECT id, login, users.public 
		FROM users 
		WHERE users.public -> 'is_api_enabled' = 'true' AND login = ?`,
		name,
	).Find(&user).Error
}

func generatePassword(password string) (*string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, models.ErrDbPassword
	}
	generatedPassword := string(passwordBytes[:])
	return &generatedPassword, nil
}

func (r *UserRepository) Insert(transaction models.Transaction, newUser models.NewUser) (*models.User, error) {
	user := models.User{
		ID: uuid.New(),
	}

	// set login
	trimmedLogin := strings.TrimSpace(newUser.Login)
	user.Login = trimmedLogin

	// set e-mail
	email, err := govalidator.NormalizeEmail(newUser.Email)
	if err != nil {
		return nil, models.ErrDbEmail
	}
	newToken := hex.EncodeToString(utils.NewEntropy(32))
	user.Email = email
	user.EmailConfirmationToken = &newToken
	user.EmailConfirmationTokenCreatedAt = utils.NewNullTime()

	// generate password
	trimmedPassword := strings.TrimSpace(newUser.Password)
	if len(trimmedPassword) < 6 {
		return nil, models.ErrDbPassword
	}
	generatedPassword, err := generatePassword(trimmedPassword)
	if err != nil {
		return nil, models.ErrDbPassword
	}
	user.Password = *generatedPassword

	// create in db
	if err := dbConn(transaction, r.DBConn).Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Patch(
	transaction models.Transaction, patchedUser models.PatchedUser, previousUser models.User,
) (*models.User, error) {
	user := models.User{
		ID: patchedUser.ID,
	}

	if patchedUser.Login != nil {
		login := strings.TrimSpace(*patchedUser.Login)
		if login == "" || len(login) < 3 {
			return nil, models.ErrDbLogin
		}
		user.Login = *patchedUser.Login
	}

	// reset email confirmation token (means the user is confirmed)
	if patchedUser.EmailConfirmationToken != nil {
		after := time.Now().Add(-12 * time.Hour)
		if *patchedUser.EmailConfirmationToken != *previousUser.EmailConfirmationToken ||
			previousUser.EmailConfirmationTokenCreatedAt == nil ||
			previousUser.EmailConfirmationTokenCreatedAt.Time.Before(after) {
			return nil, models.ErrIncorrectConfirmationToken
		}
		newToken := ""
		user.EmailConfirmationToken = &newToken
	}

	// update email
	if patchedUser.Email != nil {
		email, err := govalidator.NormalizeEmail(*patchedUser.Email)
		if err != nil {
			return nil, models.ErrDbEmail
		}
		newToken := hex.EncodeToString(utils.NewEntropy(32))
		user.Email = email
		user.EmailConfirmationToken = &newToken
		user.EmailConfirmationTokenCreatedAt = utils.NewNullTime()
	}

	// update password
	if patchedUser.Password != nil {
		trimmedPassword := strings.TrimSpace(*patchedUser.Password)
		if len(trimmedPassword) < 6 {
			return nil, models.ErrDbPassword
		}

		if patchedUser.RemindPasswordToken != nil &&
			previousUser.RemindPasswordToken != nil &&
			*previousUser.RemindPasswordToken != "" &&
			previousUser.RemindPasswordTokenCreatedAt != nil &&
			previousUser.RemindPasswordTokenCreatedAt.Valid {
			after := time.Now().Add(-2 * time.Hour)

			if *patchedUser.RemindPasswordToken != *previousUser.RemindPasswordToken ||
				previousUser.RemindPasswordTokenCreatedAt.Time.Before(after) {
				return nil, models.ErrIncorrectRemindPasswordToken
			}

			remindPasswordToken := ""
			user.RemindPasswordToken = &remindPasswordToken
		} else {
			if patchedUser.OldPassword == nil || *patchedUser.OldPassword == "" {
				return nil, models.ErrIncorrectOldPassword
			}
			if err := bcrypt.CompareHashAndPassword(
				[]byte(previousUser.Password),
				[]byte(*patchedUser.OldPassword),
			); err != nil {
				return nil, models.ErrIncorrectOldPassword
			}
		}

		// generate new password
		generatedPassword, err := generatePassword(trimmedPassword)
		if err != nil {
			return nil, models.ErrDbPassword
		}
		user.Password = *generatedPassword
	}

	// remind password
	if patchedUser.ShouldGenerateRemindPasswordToken {
		newToken := hex.EncodeToString(utils.NewEntropy(64))
		user.RemindPasswordToken = &newToken
		user.RemindPasswordTokenCreatedAt = utils.NewNullTime()
	}

	if patchedUser.Public != nil {
		user.Public = jsonMarshal(patchedUser.Public)
	}

	// update in db
	db := dbConn(transaction, r.DBConn).Model(&user).Updates(user)
	return &user, db.Error
}

func (r *UserRepository) PatchStripeData(
	transaction models.Transaction, patchedUser models.PatchedUser,
) (*models.User, error) {
	user := models.User{
		ID: patchedUser.ID,
	}

	if patchedUser.StripeSessionId != nil {
		user.StripeSessionId = patchedUser.StripeSessionId
	}
	if patchedUser.IsSubscribed != nil {
		if *patchedUser.IsSubscribed {
			user.SubscribedAt = utils.NewNullTime()
		} else {
			user.SubscribedAt = &utils.NullTime{}
		}
	}

	db := dbConn(transaction, r.DBConn).Model(&user).Updates(user)
	return &user, db.Error
}
