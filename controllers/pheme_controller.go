package controllers

import (
	"time"

	"github.com/feserr/pheme-user/models"
	"github.com/gofiber/fiber/v2"
)

// SecretKey the secret key.
const SecretKey = "secret"

// GetAllPhemes godoc
// @Summary      Retrieve all phemes
// @Description  get all phemes
// @Tags         phemes
// @Produce      json
// @Success      200  {object}  []models.Pheme
// @Failure      401  {object}  models.Message
// @Router       /pheme [get]
func GetAllPhemes(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	phemes, err := models.FetchAllPhemes(user.ID)
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "No phemes found for the user",
		})
	}

	return c.JSON(phemes)
}

// GetUserPhemes godoc
// @Summary      Retrieve the user phemes
// @Description  get the user phemes
// @Tags         phemes
// @Produce      json
// @Success      200  {object}  []models.Pheme
// @Failure      401  {object}  models.Message
// @Router       /pheme/mine [get]
func GetUserPhemes(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	phemes, err := models.FetchUserPhemes(user.ID, byte(models.PRIVATE))
	if err != nil {
		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "No phemes found for the user",
		})
	}

	return c.JSON(phemes)
}

// GetPheme godoc
// @Summary      Retrieve the pheme
// @Description  get the pheme
// @Tags         phemes
// @Produce      json
// @Param        id   path      int  true  "Pheme ID"
// @Success      200  {object}  models.Pheme
// @Failure      204  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /pheme/{id} [get]
func GetPheme(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	var paramsPhemeID models.PhemeParamsID
	if err := c.ParamsParser(&paramsPhemeID); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	phemes, err := models.FetchPheme(paramsPhemeID.ID, user.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "No phemes found",
		})
	}

	return c.JSON(phemes)
}

// PostPheme godoc
// @Summary      Post a pheme to the user
// @Description  post a user pheme
// @Tags         phemes
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.PhemeParamsID
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /pheme [post]
func PostPheme(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	var body models.PhemeParamsPost
	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid JSON body",
		})
	}

	if err := validate.Struct(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		println(err.Error())
		return c.JSON(fiber.Map{
			"message": "Wrong JSON params",
		})
	}

	pheme := models.Pheme{}
	pheme.Version = models.PhemeVersion()
	pheme.CreatedAt = time.Now()
	pheme.Visibility = byte(body.Visibilty)
	pheme.Category = body.Category
	pheme.Text = body.Text
	pheme.CreatedBy = user.ID
	pheme.UserID = body.UserID

	id, err := models.CreatePheme(pheme)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to insert pheme",
		})
	}

	if id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Cannot create phemes for non-friends users",
		})
	}

	return c.JSON(models.PhemeParamsID{ID: id})
}

// DeletePheme godoc
// @Summary      Delete a pheme from the user
// @Description  delete a user pheme
// @Tags         phemes
// @Produce      json
// @Param        id   path      int  true  "Pheme ID"
// @Success      200  {object}  models.PhemeParamsID
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /pheme/{id} [delete]
func DeletePheme(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	var paramsDelete models.PhemeParamsID
	if err := c.ParamsParser(&paramsDelete); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	id, err := models.DeletePheme(paramsDelete.ID, user.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to delete pheme",
		})
	}

	return c.JSON(models.PhemeParamsID{ID: id})
}

// UpdatePheme godoc
// @Summary      Update a pheme to the user
// @Description  update a user pheme
// @Tags         phemes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Pheme ID"
// @Success      200  {object}  models.Pheme
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /pheme/{id} [put]
func UpdatePheme(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	var paramsUpdate models.PhemeParamsID
	if err := c.ParamsParser(&paramsUpdate); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong parameters",
		})
	}

	var pheme models.PhemeParamsPost
	if err := c.BodyParser(&pheme); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid JSON body",
		})
	}

	if err := validate.Struct(pheme); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong JSON params",
		})
	}

	updatedPheme, err := models.UpdatePheme(pheme, paramsUpdate.ID, user.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to update pheme",
		})
	}

	return c.JSON(updatedPheme)
}
