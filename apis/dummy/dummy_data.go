package dummy

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/naufalziyad/RESTful-APIs/apis/models"
)

var users = []models.User{
	models.User{
		UserName: "naufal",
		FullName: "Naufal Ziyad L",
		Email:    "naufal.ziyad@detik.com",
		Password: "rahasianih",
	},
	models.User{
		UserName: "shanaya",
		FullName: "Shanaya Sekar S",
		Email:    "shanaya@gmail.com",
		Password: "rahasia2nih",
	},
}

var ads = []models.Ads{
	models.Ads{
		Title:   "Beriklan di Adsmart",
		Content: "Silahkan beriklan disini, dapat point untung sendiri",
		AdsLink: "https://adsmart.detik.com",
	},
	models.Ads{
		Title:   "Belanja di Shanaya store",
		Content: "Rebut hadiahnya , dapat untungnya",
		AdsLink: "http://shanaya.store",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Ads{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("Cannt drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Ads{}).Error
	if err != nil {
		log.Fatalf("Cannt migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Ads{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot add dummy in users table: %v", err)
		}
		ads[i].OwnerID = users[i].ID

		err = db.Debug().Model(&models.Ads{}).Create(&ads[i]).Error
		if err != nil {
			log.Fatalf("cannot add dummy in ads table: %v", err)
		}
	}

}
