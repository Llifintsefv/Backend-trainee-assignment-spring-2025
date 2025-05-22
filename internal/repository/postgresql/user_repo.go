package postgres

import (
	"Backend-trainee-assignment-spring-2025/internal/domain/models"
	"Backend-trainee-assignment-spring-2025/internal/repository"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepo(db *pgxpool.Pool, logger *slog.Logger) repository.UserRepository {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func (u *UserRepo) SaveUser(ctx context.Context, email string, password string, role models.Role) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash password", "error", err)
		return err
	}

	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)`

	_, err = u.db.Exec(ctx, query, email, passwordHash, role)

	if err != nil {
		u.logger.Error("failed to create user", "error", err)
		return err
	}

	return nil
}
func (u *UserRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`

	var user models.User

	err := u.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err.Error() == "no rows in result set" {
			u.logger.Warn("user not found", "email", email)
			return models.User{}, models.ErrUserNotFound
		}
		u.logger.Error("failed to get user", "error", err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	query := `SELECT id, email, password_hash, role FROM users WHERE id = $1`

	var user models.User

	err := u.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err.Error() == "no rows in result set" {
			u.logger.Warn("user not found", "id", id)
			return models.User{}, models.ErrUserNotFound
		}
		u.logger.Error("failed to get user", "error", err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, email string) error {
	query := `DELETE FROM users WHERE email = $1`

	_, err := u.db.Exec(ctx, query, email)

	if err != nil {
		u.logger.Error("failed to delete user", "error", err)
		return err
	}

	return nil
}

func (u *UserRepo) ComparePassword(ctx context.Context, password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
