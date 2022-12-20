package controllers

import (
	"github.com/feserr/pheme-user/models"
	"github.com/gofiber/fiber/v2"
)

// GetUserByID godoc
// @Summary      Retrieve the user by ID
// @Description  get the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name path      id  true  "User ID"
// @Success      200  {object}  models.UserPublicData
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	var paramsName models.UserParamsID
	err := c.ParamsParser(&paramsName)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	user, err := models.FindByID(paramsName.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "No user found for that ID",
		})
	}

	userPublicData := models.UserToUserPublicData(user)

	return c.JSON(userPublicData)
}

// GetUsersByName godoc
// @Summary      Retrieve the users by name
// @Description  get the users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name path      string  true  "User name"
// @Success      200  {object}  []models.User
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/{name} [get]
func GetUsersByName(c *fiber.Ctx) error {
	var paramsName models.UserParamsName
	err := c.ParamsParser(&paramsName)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	users, err := models.FindByName(paramsName.Name)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "No user found for that name",
		})
	}

	usersPublicData := models.UsersToUsersPublicData(users)

	return c.JSON(usersPublicData)
}

// GetFriends godoc
// @Summary      Retrieve the user friends
// @Description  get the user friends
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/friends [get]
func GetFriends(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	friends, err := models.GetFriends(user.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "Failed to get the friends",
		})
	}

	return c.JSON(friends)
}

// GetFollowers godoc
// @Summary      Retrieve the user followers
// @Description  get the user followers
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name path      id  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/followers/{id} [get]
func GetFollowers(c *fiber.Ctx) error {
	var paramsName models.UserParamsID
	err := c.ParamsParser(&paramsName)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	followers, err := models.GetFollowers(paramsName.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "Failed to get the friends",
		})
	}

	return c.JSON(followers)
}

// AddFriend godoc
// @Summary      Add a friends to the user
// @Description  put a friend to the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path      string  true  "Friend ID"
// @Success      200  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/friend/{id} [put]
func AddFriend(c *fiber.Ctx) error {
	var paramsID models.UserParamsID
	err := c.ParamsParser(&paramsID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	if user.ID == paramsID.ID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User ID and friend ID are the same",
		})
	}

	friend, err := models.AddFriend(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fail to add the friend",
		})
	}

	return c.JSON(friend)
}

// AddFollower godoc
// @Summary      Add a follower to the user
// @Description  put a follower to the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path      string  true  "Follower ID"
// @Success      200  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/follower/{id} [put]
func AddFollower(c *fiber.Ctx) error {
	var paramsID models.UserParamsID
	err := c.ParamsParser(&paramsID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	if user.ID == paramsID.ID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User ID and follower ID are the same",
		})
	}

	follower, err := models.AddFollower(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fail to add the follower",
		})
	}

	return c.JSON(follower)
}

// DeleteFriend godoc
// @Summary      Delete a friend of the user
// @Description  delete a friend of the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path      string  true  "Friend ID"
// @Success      200  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/friend/{id} [delete]
func DeleteFriend(c *fiber.Ctx) error {
	var paramsID models.UserParamsID
	err := c.ParamsParser(&paramsID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	if user.ID == paramsID.ID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User ID and friend ID are the same",
		})
	}

	err = models.RemoveFriend(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "Fail to delete the friend",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

// DeleteFollower godoc
// @Summary      Delete a follower of the user
// @Description  Delete a follower of the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path      string  true  "Follower ID"
// @Success      200  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user/follower/{id} [delete]
func DeleteFollower(c *fiber.Ctx) error {
	var paramsID models.UserParamsID
	err := c.ParamsParser(&paramsID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	if user.ID == paramsID.ID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User ID and follower ID are the same",
		})
	}

	err = models.RemoveFollower(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "Fail to delete the follower",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}
