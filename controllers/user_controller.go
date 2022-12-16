package controllers

import (
	"github.com/feserr/pheme-user/models"
	"github.com/gofiber/fiber/v2"
)

// GetUsersByName godoc
// @Summary      Retrieve the user phemes
// @Description  get the user phemes
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

	return c.JSON(users)
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

	err = models.AddFriend(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fail to add the friend",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
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

	err = models.AddFollower(user.ID, paramsID.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fail to add the follower",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
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
