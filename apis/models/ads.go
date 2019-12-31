package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Ads struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	AdsLink   string    `gorm:"size:255;not null;" json:"ads_link"`
	Owner     User      `json:"owner"`
	OwnerID   uint32    `gorm:"not null" json:"owner_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (a *Ads) Prepare() {
	a.ID = 0
	a.Title = html.EscapeString(strings.TrimSpace(a.Title))
	a.Content = html.EscapeString(strings.TrimSpace(a.Content))
	a.Owner = User{}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Ads) Validate() error {
	if a.Title == "" {
		return errors.New("Please Input Title")
	}
	if a.Content == "" {
		return errors.New("Please Input Content")
	}
	if a.AdsLink == "" {
		return errors.New("Please Input URL Links for your Ads")
	}
	if a.OwnerID < 1 {
		return errors.New("Please Input Owner")
	}
	return nil
}

func (a *Ads) SaveAds(db *gorm.DB) (*Ads, error) {
	var err error
	err = db.Debug().Model(&Ads{}).Create(&a).Error
	if err != nil {
		return &Ads{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Ads{}, err
		}
	}
	return a, nil
}

func (a *Ads) FindAllAds(db *gorm.DB) (*[]Ads, error) {
	var err error
	ads := []Ads{}
	err = db.Debug().Model(&Ads{}).Limit(100).Find(&ads).Error
	if err != nil {
		return &[]Ads{}, err
	}
	if len(ads) > 0 {
		for i, _ := range ads {
			err := db.Debug().Model(&User{}).Where("id = ?", ads[i].OwnerID).Take(&ads[i].Owner).Error
			if err != nil {
				return &[]Ads{}, err
			}
		}
	}
	return &ads, nil
}

func (a *Ads) FindAdsByID(db *gorm.DB, pid uint64) (*Ads, error) {
	var err error
	err = db.Debug().Model(&Ads{}).Where("id = ?", pid).Take(&a).Error

	if gorm.IsRecordNotFoundError(err) {
		return &Ads{}, errors.New("Ads Not Found")
	}

	if err != nil {
		return &Ads{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Ads{}, err
		}
	}
	return a, nil
}

func (a *Ads) UpdateAds(db *gorm.DB) (*Ads, error) {
	var err error
	err = db.Debug().Model(&Ads{}).Where("id = ?", a.ID).Updates(Ads{
		Title:     a.Title,
		Content:   a.Content,
		AdsLink:   a.AdsLink,
		UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Ads{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Ads{}, err
		}
	}
	return a, nil
}

func (a *Ads) DeleteAds(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Ads{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
