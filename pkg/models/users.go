package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"

	"habiko-go/pkg/utils"
)

var ErrDbEmail = errors.New("e-mail is incorrect")
var ErrDbLogin = errors.New("login is incorrect")
var ErrDbPassword = errors.New("password is incorrect")
var ErrIncorrectOldPassword = errors.New("old password is incorrect")
var ErrIncorrectConfirmationToken = errors.New("confirmation token is incorrect")
var ErrIncorrectRemindPasswordToken = errors.New("remind password code is incorrect")

type NewUser struct {
	Email    string `json:"email" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PatchedUser struct {
	ID                     uuid.UUID       `json:"-"`
	Email                  *string         `json:"email"`
	Login                  *string         `json:"login"`
	OldPassword            *string         `json:"old_password"`
	Password               *string         `json:"password"`
	EmailConfirmationToken *string         `json:"email_confirmation_token"`
	RemindPasswordToken    *string         `json:"remind_password_token"`
	IsSubscribed           *bool           `json:"-"`
	Public                 *datatypes.JSON `json:"public"`

	// stripe
	StripeSessionId *string `json:"stripe_session_id"`

	// extra flags
	ShouldGenerateRemindPasswordToken bool `json:"-"`
}

type User struct {
	ID                              uuid.UUID       `json:"id" gorm:"primary_key"`
	Email                           string          `json:"email" gorm:"unique"`
	Login                           string          `json:"login" gorm:"unique"`
	Password                        string          `json:"-"`
	EmailConfirmationToken          *string         `json:"-"`
	EmailConfirmationTokenCreatedAt *utils.NullTime `json:"-"`
	RemindPasswordToken             *string         `json:"-"`
	RemindPasswordTokenCreatedAt    *utils.NullTime `json:"-"`
	CreatedAt                       *time.Time      `json:"created_at"`
	UpdatedAt                       *time.Time      `json:"updated_at"`
	DeletedAt                       *time.Time      `json:"deleted_at,omitempty" sql:"index"`
	Public                          datatypes.JSON  `json:"public" gorm:"type:JSONB" sql:"DEFAULT:'{}'"`

	// stripe
	SubscribedAt                *utils.NullTime `json:"subscribed_at,omitempty"`
	StripeSessionId             *string         `json:"stripe_session_id"`

	// user options
	CanUpdateOlderEntries *bool `json:"can_update_older_entries" gorm:"default:false"`

	// extra flags to set and return
	IsEmailConfirmed bool `json:"is_email_confirmed" gorm:"-"`
}

// PublicUserData is the data that's being returned with public API
type PublicUserData struct {
	ID     uuid.UUID       `json:"-"`
	Login  string          `json:"login"`
	Public *datatypes.JSON `json:"public"`
}

type UserRepository interface {
	IsEmailConfirmed(user User) bool

	FindByID(transaction Transaction, id uuid.UUID) (*User, error)
	FindByEmail(transaction Transaction, email string) (*User, error)
	FindByCredentials(transaction Transaction, email string, password string) (*User, error)
	FindPublicDataByName(transaction Transaction, name string) (*PublicUserData, error)

	Insert(transaction Transaction, newUser NewUser) (*User, error)
	Patch(transaction Transaction, patchedUser PatchedUser, previousUser User) (*User, error)
	PatchStripeData(transaction Transaction, patchedUser PatchedUser) (*User, error)
}
