package services

import (
	"errors"
	"net/http"
	"os"
	"pos-services/app/models"
	"pos-services/pkg/utils"
	"pos-services/platform/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrEmailAlreadyRegistered   = errors.New("email already registered")
	ErrAccountWithEmailNotFound = errors.New("record not found, please check email you entered")
	ErrPasswordMismatch         = errors.New("wrong password")
)

func CreateAccount(a *models.Users) ServiceReturn {
	// validate input
	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return ServiceReturn{
			HttpStatusCode: http.StatusForbidden,
			Err:            err,
		}
	}

	// check if user is already registered
	db, err := database.OpenDBConnection()
	if err != nil {
		return ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}
	if _, err = db.GetByEmail(a.Email); err == nil {
		return ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            ErrEmailAlreadyRegistered,
		}
	}

	// generate hashed password
	p := &utils.Params{
		Memory:      64 * 1024,
		Iter:        3,
		Parallelism: 2,
		LenSalt:     16,
		LenKey:      32,
	}

	encodedHash, err := utils.GenerateHashedPassword(a.Password, p)
	if err != nil {
		return ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}
	a.Password = encodedHash

	// save to DB
	if err = db.CreateUser(a); err != nil {
		return ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	return ServiceReturn{
		HttpStatusCode: http.StatusCreated,
		Err:            nil,
	}
}

func SignIn(email, pwd string) (models.Users, ServiceReturn) {
	// check if user is already registered
	db, err := database.OpenDBConnection()
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	// check if email not registered
	user, err := db.GetByEmail(email)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusNotFound,
			Err:            ErrAccountWithEmailNotFound,
		}
	}

	// check if password not match
	match, err := utils.CompareSavedWithIncomingPassword(pwd, user.Password)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	if !match {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusForbidden,
			Err:            ErrPasswordMismatch,
		}
	}

	// check if already login
	session, err := db.GetSessionByEmail(email)
	if err == nil {
		return *user, ServiceReturn{
			HttpStatusCode: http.StatusOK,
			Err:            nil,
			Payload:        session,
		}
	}

	// create new session when not
	accessTokenDuration := os.Getenv("ACCESS_TOKEN_DURATION")
	timeDuration, err := time.ParseDuration(accessTokenDuration)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	p, err := utils.NewPayloadSession(email, timeDuration)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	if err = db.CreateSession(p); err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	symetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")
	pst, err := utils.NewPasetoMaker(symetricKey)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	token, err := pst.CreateToken(email, timeDuration)
	if err != nil {
		return models.Users{}, ServiceReturn{
			HttpStatusCode: http.StatusInternalServerError,
			Err:            err,
		}
	}

	// clear password
	user.Password = ""

	return *user, ServiceReturn{
		HttpStatusCode: http.StatusOK,
		Err:            nil,
		Payload: fiber.Map{
			"token": token,
		},
	}
}
