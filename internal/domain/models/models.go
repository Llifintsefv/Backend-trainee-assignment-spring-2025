package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type City string

const (
	CityMoscow          City = "Москва"
	CitySaintPetersburg City = "Санкт-Петербург"
	CityKazan           City = "Казань"
)

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

type ProductType string

const (
	TypeElectronics ProductType = "электроника"
	TypeClothing    ProductType = "одежда"
	TypeShoes       ProductType = "обувь"
)

type ReceptionStatus string

const (
	StatusInProgress ReceptionStatus = "in_progress"
	StatusClosed     ReceptionStatus = "close"
)

type PVZ struct {
	ID                uuid.UUID
	City              City
	RegistrationsData time.Time
}

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
}

type Reception struct {
	ID       uuid.UUID
	DataTime time.Time
	PVZID    uuid.UUID
	Status   ReceptionStatus
}

type Product struct {
	ID          uuid.UUID
	DataTime    time.Time
	Type        ProductType
	ReceptionID uuid.UUID
}

// Ответы API

type ProductInfo struct {
	Product Product
}

type ReceptionInfo struct {
	Reception Reception
	Product   []ProductInfo
}

type PVZInfo struct {
	PVZ       PVZ
	Reception []ReceptionInfo
}

// JWT

type AccessTokenClaims struct {
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip"`
	Role      Role   `json:"role"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        uuid.UUID
	JTI       string
	UserID    string
	TokenHash string
	IPAddress string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// request

type DummyLoginRequest struct {
	Role Role `json:"role" validate:"required,oneof=moderator employee"`
}

type RefreshRequest struct {
	RefreshToken string `json:"RefreshToken" validate:"required,refreshTokenFormat"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     Role   `json:"role" validate:"required,oneof=employee moderator"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type PVZRequest struct {
	City City `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
}

// Errors

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
