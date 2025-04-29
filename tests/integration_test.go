package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sales-api/api"
	"sales-api/internal"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateAndGet(t *testing.T) {
	app := gin.Default()

	httpClient := resty.New()

	// Crear una respuesta simulada para la llamada GET /users/1234
	user := User{ID: "1f438b44-3f69-4e54-a8f6-5a3814910e5a"}
	userJSON, _ := json.Marshal(user)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBuffer(userJSON)),
	}

	// Usamos MockRoundTripper para simular las respuestas de HTTP
	httpClient.SetTransport(&MockRoundTripper{
		response: mockResponse,
		err:      nil, // No hay error
	})

	localStorage := internal.NewLocalStorage()

	api.InitRoutes(app, httpClient, localStorage)

	// POST
	req, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewBufferString(`{
		"user_id": "1f438b44-3f69-4e54-a8f6-5a3814910e5a",
  		"amount": 150.0
	}`))

	res := fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.Code)

	var resVenta *internal.Venta
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resVenta))
	require.Equal(t, "1f438b44-3f69-4e54-a8f6-5a3814910e5a", resVenta.UserID)
	require.Equal(t, 150.0, resVenta.Amount)
	require.Equal(t, 1, resVenta.Version)
	require.NotEmpty(t, resVenta.ID)
	require.NotEmpty(t, resVenta.Status)
	require.NotEmpty(t, resVenta.CreatedAt)
	require.NotEmpty(t, resVenta.UpdatedAt)

	venta, _ := localStorage.Read(resVenta.ID)
	venta.Status = "pending"
	localStorage.Set(venta)

	ventaId := resVenta.ID
	venta, _ = localStorage.Read(ventaId)

	// PATCH
	req, _ = http.NewRequest(http.MethodPatch, "/sales/"+venta.ID, bytes.NewBufferString(`{
		"status": "approved"
	}`))

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resVenta))
	require.Equal(t, "approved", resVenta.Status)

	// GET
	req, _ = http.NewRequest(http.MethodGet, "/sales/"+resVenta.ID, nil)

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, ventaId, resVenta.ID)
	require.Equal(t, "1f438b44-3f69-4e54-a8f6-5a3814910e5a", resVenta.UserID)
	require.Equal(t, 150.0, resVenta.Amount)
	require.Equal(t, 2, resVenta.Version)
	require.Equal(t, "approved", resVenta.Status)
}

func fakeRequest(e *gin.Engine, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	return w
}

// User es la estructura que simula el retorno de un usuario
type User struct {
	ID string `json:"id"`
}

// MockRoundTripper es una estructura que implementa RoundTripper para mockear las respuestas HTTP
type MockRoundTripper struct {
	// response se usa para devolver una respuesta simulada
	response *http.Response
	// err se usa para devolver un error simulado
	err error
}

// RoundTrip es el método que implementa el RoundTripper y es llamado cuando Resty hace la solicitud
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}
