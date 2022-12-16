package models

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// UserVersion returns the current version of the user schema.
func UserVersion() uint {
	return 1
}

func init() {
	Db.AutoMigrate(User{})
}

// User model info
// @Description User account
type User struct {
	ID           uint      `json:"id"`
	Version      uint      `json:"version" gorm:"not null"`
	Name         string    `json:"userName" gorm:"not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     []byte    `json:"-"  gorm:"not null"`
	PasswordDate time.Time `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"not null"`
	Followers    []User    `json:"-" gorm:"many2many:followship;association_jointable_foreignkey:follow_id"`
	Friends      []User    `json:"-" gorm:"many2many:friendship;association_jointable_foreignkey:friend_id"`
}

// GetUser returns the logged user.
func GetUser(c *fiber.Ctx, secretKey string) (User, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	var user User
	if err != nil {
		return user, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	Db.Where("id = ?", claims.Issuer).First(&user)

	return user, nil
}

// DeleteByID deletes the user by the ID.
func DeleteByID(userID uint) error {
	if err := Db.Delete(&User{}, userID); err.Error != nil {
		println(err.Error)
		return err.Error
	}

	return nil
}

// FindByID returns the user from the ID.
func FindByID(userID uint) (*User, error) {
	user := &User{}
	if err := Db.First(user, userID).Error; err != nil {
		println(err)
		return user, err
	}

	return user, nil
}

// FindByName returns the users that contains the name.
func FindByName(userName string) ([]User, error) {
	users := []User{}
	usersByName := Db.Model(&User{}).Select("id, name").Order("created_at desc").Find(&users, "name LIKE ?", "%"+userName+"%")
	if usersByName.Error != nil {
		println(usersByName.Error)
		return users, usersByName.Error
	}

	return users, nil
}

// IsFriend returns if it is friend or not.
func IsFriend(userID uint, friendID uint) (bool, error) {
	friend := User{}
	friend.ID = friendID

	user := User{}
	user.ID = userID

	isFriend := Db.Model(&user).Association("Friends").Find(&friend)
	if isFriend != nil {
		println(isFriend)
		return false, isFriend
	}

	if friend.Email == "" {
		return false, nil
	}

	return true, nil
}

// GetFriends returns the friends of a user.
func GetFriends(userID uint) (*[]uint, error) {
	friends := &[]uint{}
	allFriends := Db.Table("friendship").Select("friend_id").Find(&friends, "user_id = ?", userID)
	if allFriends.Error != nil {
		println(allFriends.Error)
		return friends, allFriends.Error
	}

	return friends, nil
}

// GetFollowers returns the followers of a user.
func GetFollowers(userID uint) (*[]uint, error) {
	followers := &[]uint{}
	allFollowers := Db.Table("followship").Select("follower_id").Find(&followers, "user_id = ?", userID)
	if allFollowers.Error != nil {
		println(allFollowers.Error)
		return followers, allFollowers.Error
	}

	return followers, nil
}

// AddFriend adds a friends to a user.
func AddFriend(userID uint, friendID uint) error {
	user := User{}

	friend, err := FindByID(friendID)
	if err != nil {
		println(err)
		return err
	}

	Db.Preload("Friends").First(&user, "id = ?", userID)
	err = Db.Model(&user).Association("Friends").Append(friend)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// AddFollower adds a follower to a user.
func AddFollower(userID uint, followerID uint) error {
	user := User{}

	follower, err := FindByID(followerID)
	if err != nil {
		println(err)
		return err
	}

	Db.Preload("Followers").First(&user, "id = ?", userID)
	err = Db.Model(&user).Association("Followers").Append(follower)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RemoveFriend removes a friends for a user.
func RemoveFriend(userID uint, friendID uint) error {
	user := User{}

	friend, err := FindByID(friendID)
	if err != nil {
		println(err)
		return err
	}

	Db.Preload("Friends").First(&user, "id = ?", userID)
	err = Db.Model(&user).Association("Friends").Delete(friend)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RemoveFollower removes a follower for a user.
func RemoveFollower(userID uint, followerID uint) error {
	user := User{}

	follower, err := FindByID(followerID)
	if err != nil {
		println(err)
		return err
	}

	Db.Preload("Followers").First(&user, "id = ?", userID)
	err = Db.Model(&user).Association("Followers").Delete(follower)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
