package internal

import "errors"

var ErrEmptyID = errors.New("empty user ID")
var ErrNotFound = errors.New("venta no encontrada")

type LocalStorage struct {
	ventasStore map[string]*Venta
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		ventasStore: map[string]*Venta{},
	}
}

func (storage *LocalStorage) Set(venta *Venta) error {
	if venta.ID == "" {
		return ErrEmptyID
	}

	storage.ventasStore[venta.ID] = venta
	return nil
}

func (storage *LocalStorage) Read(id string) (*Venta, error) {
	venta, exists := storage.ventasStore[id]
	if !exists {
		return nil, ErrNotFound
	}

	return venta, nil
}

func (storage *LocalStorage) Delete(id string) error {
	_, err := storage.Read(id)
	if err != nil {
		return err
	}

	delete(storage.ventasStore, id)
	return nil
}

func (storage *LocalStorage) Search(userID string, status string) ([]*Venta, error) {
	var ventasEncontradas []*Venta

	for _, venta := range storage.ventasStore {
		if venta.UserID == userID {
			if status == "" || venta.Status == status {
				ventasEncontradas = append(ventasEncontradas, venta)
			}
		}
	}
	return ventasEncontradas, nil
}
