package models

import (
	"errors"
	"log"
	"time"
)

// PhemeVersion returns the version of the Pheme schema.
func PhemeVersion() uint {
	return 1
}

func init() {
	err := Db.AutoMigrate(Pheme{})
	if err != nil {
		panic("Couldn't migrate DB")
	}
}

// Pheme model info
// @Description Pheme content
type Pheme struct {
	ID         uint      `json:"id"`
	Version    uint      `json:"version" gorm:"not null" validate:"required"`
	CreatedAt  time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Visibility byte      `json:"visibility" sql:"visibility" gorm:"not null" validate:"required"`
	Category   string    `json:"category" gorm:"not null" validate:"required"`
	Text       string    `json:"text" gorm:"not null" validate:"required"`
	CreatedBy  uint      `json:"createdBy" gorm:"not null"`
	UserID     uint      `json:"userId" gorm:"not null" validate:"required"`
}

// FetchAllPhemes returns all the phemes of a user, friends and followers with equal or higher visibility.
func FetchAllPhemes(userID uint) (*[]Pheme, error) {
	phemes := &[]Pheme{}
	allUserPhemes := Db.Model(&Pheme{}).Order("created_at desc").Find(&phemes, "user_id = ? and visibility >= ?", userID, byte(PRIVATE))
	if allUserPhemes.Error != nil {
		println(allUserPhemes.Error)
		return phemes, allUserPhemes.Error
	}

	friends, err := GetFriends(userID)
	if err == nil && len(*friends) > 0 {
		friendsPhemes := &[]Pheme{}
		allFriendsPhemes := Db.Model(&Pheme{}).Order("created_at desc").Find(friendsPhemes, "user_id in ? and visibility >= ?", *friends, byte(PROTECTED))
		if allFriendsPhemes.Error == nil {
			*phemes = append(*phemes, *friendsPhemes...)
		}
	}

	followers, err := GetFollowers(userID)
	if err == nil && len(*followers) > 0 {
		followersPhemes := &[]Pheme{}
		allFollowersPhemes := Db.Model(&Pheme{}).Order("created_at desc").Find(followersPhemes, "user_id in ? and visibility >= ?", *followers, byte(PUBLIC))
		if allFollowersPhemes.Error == nil {
			*phemes = append(*phemes, *followersPhemes...)
		}
	}

	return phemes, nil
}

// FetchUserPhemes returns all the phemes of the logged user with equal or higher visibility.
func FetchUserPhemes(userID uint, visibility byte) (*[]Pheme, error) {
	phemes := &[]Pheme{}
	allPhemes := Db.Model(&Pheme{}).Order("created_at desc").Find(&phemes, "user_id = ? and visibility >= ?", userID, visibility)
	if allPhemes.Error != nil {
		println(allPhemes.Error)
		return phemes, allPhemes.Error
	}

	return phemes, nil
}

// FetchPheme returns the pheme if is visible to the user.
func FetchPheme(phemeID uint, userID uint) (*Pheme, error) {
	pheme := &Pheme{}
	thePheme := Db.Model(&Pheme{}).Find(&pheme, phemeID)
	if thePheme.Error != nil {
		println(thePheme.Error)
		return pheme, thePheme.Error
	}

	if pheme.CreatedBy != userID || pheme.UserID != userID {
		return nil, errors.New("pheme not visible for the user")
	}

	return pheme, nil
}

// FetchPhemes returns the phemes if are visible to the user.
func FetchPhemes(phemeID uint, userID uint, visibility byte) (*[]Pheme, error) {
	phemes := &[]Pheme{}
	thePhemes := Db.Model(&Pheme{}).Find(&phemes, "user_id = ? and visibility >= ?", phemeID, visibility)
	if thePhemes.Error != nil {
		println(thePhemes.Error)
		return phemes, thePhemes.Error
	}

	visiblePhemes := &[]Pheme{}
	for _, pheme := range *phemes {
		if pheme.CreatedBy == userID || pheme.UserID == userID {
			*visiblePhemes = append(*visiblePhemes, pheme)
		}
	}

	return visiblePhemes, nil
}

// CreatePheme adds a pheme to the DB.
func CreatePheme(pheme Pheme) (uint, error) {
	if pheme.CreatedBy != pheme.UserID {
		res, err := IsFriend(pheme.CreatedBy, pheme.UserID)
		if err != nil {
			println(err)
			return 0, err
		}

		if !res {
			return 0, nil
		}
	}

	createdPheme := Db.Create(&pheme)
	if createdPheme.Error != nil {
		log.Println(createdPheme.Error)
		return pheme.ID, createdPheme.Error
	}

	return pheme.ID, nil
}

// DeletePheme removes a pheme from a user.
func DeletePheme(phemeID uint, userID uint) (uint, error) {
	deletedPheme := Db.Unscoped().Delete(Pheme{}, "id = ? AND user_id = ?", phemeID, userID)
	if deletedPheme.Error != nil {
		log.Println(deletedPheme.Error)
		return phemeID, deletedPheme.Error
	}

	if deletedPheme.RowsAffected < 1 {
		return phemeID, errors.New("couldn't delete because it don't exist")
	}

	return phemeID, nil
}

// UpdatePheme updates the data of a pheme.
func UpdatePheme(pheme PhemeParamsPost, phemeID uint, userID uint) (Pheme, error) {
	oldPheme := Pheme{}
	updatedPost := Db.First(&oldPheme, "id = ? AND created_by = ?", phemeID, userID)
	if updatedPost.Error != nil {
		log.Println(updatedPost.Error)
		return oldPheme, updatedPost.Error
	}

	oldPheme.Version = PhemeVersion()
	oldPheme.UpdatedAt = time.Now()
	oldPheme.Visibility = pheme.Visibilty
	oldPheme.Category = pheme.Category
	oldPheme.Text = pheme.Text
	updatedPost = Db.Save(&oldPheme)
	if updatedPost.Error != nil {
		log.Println(updatedPost.Error)
		return oldPheme, updatedPost.Error
	}

	return oldPheme, nil
}
