package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

var (
	ErrInvalidAmount     = errors.New("el monto debe ser mayor a cero")
	ErrUserNotFound      = errors.New("el usuario no existe")
	ErrVentaNotFound     = errors.New("la venta no existe")
	ErrInvalidTransition = errors.New("solo se puede cambiar el estado si está en 'pending'")
	ErrInvalidStatus     = errors.New("estado inválido")
)

type Service struct {
	storage *LocalStorage
	client  *resty.Client
}

func NewService(storage *LocalStorage, client *resty.Client) *Service {
	client.SetBaseURL("http://localhost:8081")

	return &Service{
		storage: storage,
		client:  client,
	}
}

func (service *Service) Create(venta *Venta) error {
	if venta.Amount < 0 {
		return ErrInvalidAmount
	}

	resp, err := service.client.R().
		SetHeader("Content-Type", "application/json").
		Get("http://localhost:8080/users/" + venta.UserID)

	if err != nil {
		return err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		fmt.Println("Error en la respuesta de usuarios:", resp.Status())
		fmt.Println("Cuerpo de la respuesta:", resp.Body())
		return ErrUserNotFound
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	estados := []string{"pending", "approved", "rejected"}
	venta.Status = estados[rnd.Intn(len(estados))]

	venta.ID = uuid.NewString()
	now := time.Now()
	venta.CreatedAt = now
	venta.UpdatedAt = now
	venta.Version = 1

	return service.storage.Set(venta)
}

func (service *Service) UpdateStatus(id string, status string) (*Venta, error) {
	venta, err := service.storage.Read(id)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrVentaNotFound
		}
		return nil, err
	}

	if venta.Status != "pending" {
		return nil, ErrInvalidTransition
	}

	if status != "approved" && status != "rejected" {
		return nil, ErrInvalidStatus
	}

	venta.Status = status
	venta.UpdatedAt = time.Now()
	venta.Version++

	err = service.storage.Set(venta)
	if err != nil {
		return nil, err
	}

	return venta, nil
}

func (service *Service) SearchVentas(userID string, status string) ([]*Venta, error) {
	return service.storage.Search(userID, status)
}
