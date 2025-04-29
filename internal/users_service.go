package internal

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type UserService interface {
	userExist(userID string) (bool, error)
}

type userService struct {
	client *resty.Client
}

func NewUserService(client *resty.Client) *userService {
	return &userService{
		client: client,
	}
}

func (userService *userService) userExist(userID string) (bool, error) {
	resp, err := userService.client.R().
		SetHeader("Content-Type", "application/json").
		Get("http://localhost:8080/users/" + userID)

	if err != nil {
		return false, err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		fmt.Println("Error en la respuesta de usuarios:", resp.Status())
		fmt.Println("Cuerpo de la respuesta:", resp.Body())
		return false, nil
	}

	return true, nil
}
