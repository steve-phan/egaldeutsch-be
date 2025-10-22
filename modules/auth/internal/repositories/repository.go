package repositories

import (
	"time"

	authpkg "egaldeutsch-be/internal/auth"
	models "egaldeutsch-be/modules/auth/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertRefreshToken(tokenHash string, userID string, expiresAt int64, ip *string, userAgent *string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	rt := &models.RefreshToken{
		UserID:    uid,
		TokenHash: tokenHash,
		ExpiresAt: time.Unix(expiresAt, 0),
	}
	if ip != nil {
		rt.IP = ip
	}
	if userAgent != nil {
		rt.UserAgent = userAgent
	}
	return r.db.Create(rt).Error
}

// RotateRefreshToken creates a new refresh token row and marks the old one revoked.
// It returns the user id as string and whether reuse was detected.
func (r *Repository) RotateRefreshToken(oldHash, newHash string, newExpiresAt int64, ip *string, userAgent *string) (string, bool, error) {
	// Use transaction and FOR UPDATE semantics
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", false, tx.Error
	}

	var old models.RefreshToken
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token_hash = ?", oldHash).First(&old).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return "", false, authpkg.ErrInvalidRefreshToken
		}
		return "", false, err
	}

	// detect reuse
	if old.Revoked && old.ReplacedBy != nil {
		// revoke all for user
		if err := tx.Model(&models.RefreshToken{}).Where("user_id = ?", old.UserID).Updates(map[string]interface{}{"revoked": true}).Error; err != nil {
			tx.Rollback()
			return "", false, err
		}
		tx.Commit()
		return old.UserID.String(), true, nil
	}

	// create new token
	uid := old.UserID
	newRT := &models.RefreshToken{
		UserID:    uid,
		TokenHash: newHash,
		ExpiresAt: time.Unix(newExpiresAt, 0),
	}
	if ip != nil {
		newRT.IP = ip
	}
	if userAgent != nil {
		newRT.UserAgent = userAgent
	}

	if err := tx.Create(newRT).Error; err != nil {
		tx.Rollback()
		return "", false, err
	}

	// mark old revoked and set replaced_by
	if err := tx.Model(&models.RefreshToken{}).Where("token_hash = ?", oldHash).Updates(map[string]interface{}{"revoked": true, "replaced_by": newHash}).Error; err != nil {
		tx.Rollback()
		return "", false, err
	}

	if err := tx.Commit().Error; err != nil {
		return "", false, err
	}

	return uid.String(), false, nil
}

func (r *Repository) RevokeRefreshTokenByHash(hash string, replacedBy *string) error {
	updates := map[string]interface{}{"revoked": true}
	if replacedBy != nil {
		updates["replaced_by"] = *replacedBy
	}
	res := r.db.Model(&models.RefreshToken{}).Where("token_hash = ?", hash).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return authpkg.ErrInvalidRefreshToken
	}
	return nil
}

func (r *Repository) RevokeAllForUser(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return r.db.Model(&models.RefreshToken{}).Where("user_id = ?", uid).Updates(map[string]interface{}{"revoked": true}).Error
}

func (r *Repository) InsertPasswordReset(tokenHash string, userID string, expiresAt int64) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pr := &models.PasswordReset{
		UserID:    uid,
		TokenHash: tokenHash,
		ExpiresAt: time.Unix(expiresAt, 0),
	}
	return r.db.Create(pr).Error
}

func (r *Repository) VerifyAndMarkPasswordReset(tokenHash string) (string, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	var pr models.PasswordReset
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token_hash = ? AND used = false AND expires_at > now()", tokenHash).First(&pr).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return "", authpkg.ErrInvalidRefreshToken
		}
		return "", err
	}

	if err := tx.Model(&models.PasswordReset{}).Where("token_hash = ?", tokenHash).Update("used", true).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return pr.UserID.String(), nil
}
