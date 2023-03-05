package models

import (
	"order-service/global/utils/model"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID             int        `json:"id,omitempty"`
	Email          string     `json:"email,omitempty"`
	EmailVerified  NullString `json:"email_verified,omitempty"`
	Password       string     `json:"password,omitempty"`
	LastLogin      NullString `json:"last_login,omitempty"`
	FirstName      NullString `json:"first_name,omitempty"`
	LastName       NullString `json:"last_name,omitempty"`
	RoleID         NullString `json:"role_id,omitempty"`
	MobileNumber   string     `json:"mobile_number,omitempty"`
	MobileVerified *time.Time `json:"mobile_verified,omitempty"`
	WhatsappID     NullString `json:"whatsapp_id,omitempty"`
	Status         string     `json:"status,omitempty"`
	IsAdmin        int        `json:"is_admin,omitempty"`
	GroupContentID NullString `json:"group_content_id,omitempty"`
	IsOwner        NullInt32  `json:"is_owner,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type UserChan struct {
	User     *User
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type UsersChan struct {
	Users    []*User
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type UserRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Users struct {
	Users []*User `json:"users,omitempty"`
	Total int64   `json:"total,omitempty"`
}

type UserClaims struct {
	jwt.StandardClaims
	UserID           int    `json:"user_id"`
	AgentID          int    `json:"agent_id"`
	UserEmail        string `json:"user_email"`
	UserRoleSlug     string `json:"user_role_slug"`
	UserRoleCategory string `json:"user_role_category"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
}
