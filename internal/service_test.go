package internal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

// --- mockStorage que cumple con la interfaz Storage ---
type mockStorage struct {
	SetFunc    func(*Venta) error
	ReadFunc   func(string) (*Venta, error)
	SearchFunc func(string, string) ([]*Venta, error)
}

func (m *mockStorage) Set(v *Venta) error {
	if m.SetFunc != nil {
		return m.SetFunc(v)
	}
	return nil
}

func (m *mockStorage) Read(id string) (*Venta, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(id)
	}
	return nil, nil
}

func (m *mockStorage) Search(userID string, status string) ([]*Venta, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(userID, status)
	}
	return nil, nil
}

// --- TEST: Create devuelve error si usuario no existe ---
func TestService_Create_UserNotFound(t *testing.T) {
	// logger falso
	logger := zap.NewNop()

	// storage mock
	storage := &mockStorage{
		SetFunc: func(v *Venta) error {
			t.Fatal("Set no debería ser llamado si el usuario no existe")
			return nil
		},
	}

	userService := NewMockUserService(false, nil)

	service := NewService(storage, userService, logger)

	venta := &Venta{
		UserID: "1234",
		Amount: 100,
	}

	err := service.Create(venta)

	require.ErrorIs(t, err, ErrUserNotFound)
}

// NewMockUserService devuelve una implementación de UserService que siempre retorna
// los valores pasados como parámetros.
func NewMockUserService(exists bool, err error) UserService {
	return &mockUserServiceImpl{
		exists: exists,
		err:    err,
	}
}

// Implementación privada del UserService
type mockUserServiceImpl struct {
	exists bool
	err    error
}

func (m *mockUserServiceImpl) userExist(userID string) (bool, error) {
	return m.exists, m.err
}
