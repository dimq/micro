package models

import "github.com/jinzhu/gorm"

type HashTable struct {
	Hash string `gorm:"PRIMARY_KEY;INDEX;UNIQUE"`
}

// HashUpsert /
func HashUpsert(db *gorm.DB, hash string) (bool, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, err
	}
	var tempHash HashTable

	errCheck := tx.Where("hash = ?", hash).First(&tempHash).Error
	if errCheck == gorm.ErrRecordNotFound {
		if errInsert := tx.Create(&HashTable{Hash: hash}).Error; errInsert != nil {
			tx.Rollback()
			return false, errInsert
		}

	} else {
		return false, errCheck
	}
	return true, tx.Commit().Error
}
