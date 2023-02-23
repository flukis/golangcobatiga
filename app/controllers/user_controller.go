package controllers

import (
	"net/http"
	"pos-services/app/models"
	"pos-services/app/services"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	user := &models.Users{}
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	serviceReturn := services.CreateAccount(user)

	if serviceReturn.Err != nil {
		return c.Status(serviceReturn.HttpStatusCode).JSON(fiber.Map{
			"status": false,
			"msg":    serviceReturn.Err.Error(),
		})
	}

	return c.SendStatus(http.StatusCreated)
}

func SignIn(c *fiber.Ctx) error {
	payload := &services.SignInPayload{}
	if err := c.BodyParser(payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	user, serviceReturn := services.SignIn(payload.Email, payload.Password)

	if serviceReturn.Err != nil {
		return c.Status(serviceReturn.HttpStatusCode).JSON(fiber.Map{
			"status": false,
			"msg":    serviceReturn.Err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": true,
		"msg":    "success login",
		"data": fiber.Map{
			"email":  user.Email,
			"id":     user.ID,
			"avatar": user.Avatar,
			"name":   user.Name,
			"token":  serviceReturn.Payload,
		},
	})
}

func AddUserAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	addr := &models.UserLocation{}
	if err := c.BodyParser(addr); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	res, serviceReturn := services.AddAddress(id, addr)
	if serviceReturn.Err != nil {
		return c.Status(serviceReturn.HttpStatusCode).JSON(fiber.Map{
			"status": false,
			"msg":    serviceReturn.Err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": true,
		"msg":    "success login",
		"data":   res,
	})
}
