package schemas

import "gorm.io/gorm"

type ShortedUrl struct {
	gorm.Model
	ShortId      string `json:"shortId"`
	EncryptedUrl string `json:"shortUrl"`
}
