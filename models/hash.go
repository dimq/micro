package models

import "github.com/jinzhu/gorm"

type HashTable struct {
	Hash string `gorm:"PRIMARY_KEY;INDEX;UNIQUE"`
}

// HashUpsert /
func HashUpsert(db *gorm.DB, hash string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var tempHash HashTable

		errCheck := tx.Where("hash = ?", hash).First(&tempHash).Error
		if errCheck == gorm.ErrRecordNotFound {
			if errInsert := tx.Create(&HashTable{Hash: hash}).Error; errInsert != nil {
				return errInsert
			}
		} else {
			return errCheck
		}
		return nil
	})
}
