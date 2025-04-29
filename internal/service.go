package internal

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidAmount     = errors.New("el monto debe ser mayor a cero")
	ErrUserNotFound      = errors.New("el usuario no existe")
	ErrVentaNotFound     = errors.New("la venta no existe")
	ErrInvalidTransition = errors.New("solo se puede cambiar el estado si está en 'pending'")
	ErrInvalidStatus     = errors.New("estado inválido")
	ErrInternalError     = errors.New("error interno")
)

// probando para test
type Storage interface {
	Set(*Venta) error
	Read(string) (*Venta, error)
	Search(userID string, status string) ([]*Venta, error)
}

type Service struct {
	storage     Storage // 👈 antes era *LocalStorage
	userService UserService
	logger      *zap.Logger
}

func NewService(storage Storage, userService UserService, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync()
	}

	return &Service{
		storage:     storage,
		userService: userService,
		logger:      logger,
	}
}

func (service *Service) Create(venta *Venta) error {
	if venta.Amount < 0 {
		return ErrInvalidAmount
	}

	userExist, err := service.userService.userExist(venta.UserID)
	if err != nil {
		return ErrInternalError
	}

	if !userExist {
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

	if err := service.storage.Set(venta); err != nil {
		service.logger.Error("failed to set venta", zap.Error(err), zap.Any("venta", venta))
		return err
	}

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

	//err = service.storage.Set(venta)
	//if err != nil {
	//	return nil, err
	//}

	if err := service.storage.Set(venta); err != nil {
		service.logger.Error("failed to update venta", zap.Error(err), zap.Any("venta", venta))
		return nil, err
	}

	return venta, nil
}

func (service *Service) GetVenta(ventaID string) (*Venta, error) {
	return service.storage.Read(ventaID)
}

func (service *Service) SearchVentas(userID string, status string) ([]*Venta, error) {
	return service.storage.Search(userID, status)
}
