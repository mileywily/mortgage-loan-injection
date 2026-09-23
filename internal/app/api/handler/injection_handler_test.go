package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"mortgage-loan-injection/internal/core/domain"
	"mortgage-loan-injection/internal/core/usecases"
)

type mockFinnFlowClient struct{}

func (m *mockFinnFlowClient) Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error) {
	status := 200
	msg := "Application injected successfully"
	solicitud := req.DatosCredito.NumeroSolicitud
	return domain.InyeccionResponse{
		StatusCode:      &status,
		Mensaje:         &msg,
		NumeroSolicitud: &solicitud,
	}, nil
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customRutPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`, fl.Field().String())
			return matched
		})
		v.RegisterValidation("customFechaPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{2}-[0-9]{2}-[0-9]{4}$`, fl.Field().String())
			return matched
		})
	}
	uc := usecases.NewInjectionUseCaseImpl(&mockFinnFlowClient{})
	h := NewInjectionHandler(uc)
	router.POST("/injections", h.InjectLoanApplication)
	return router
}

func TestInjectLoanApplication_Success(t *testing.T) {
	router := setupRouter()

	jsonBody := []byte(`{
		"datos_credito": {
			"NumeroSolicitud": 123,
			"AntiguedadVivienda": 5,
			"EjecutivoComercial": "Juan",
			"Producto": 1,
			"Objetivo": 2,
			"Destino": 3,
			"MontoAprobado": 100.0,
			"ValorPropiedad": 200.0,
			"FechaAprobacion": "2023-01-01",
			"ValorContado": 100.0,
			"Plazo1": 20,
			"MesesGracia": 2,
			"Tasa1": 3.5,
			"Spread1": 1.0
		},
		"participantes": [
			{
				"Rut": "12.345.678-9",
				"TipoParticipacion": 1,
				"Nombre": "A",
				"Paterno": "B",
				"Materno": "C",
				"FechaNacimiento": "01-01-1990"
			}
		]
	}`)

	req, _ := http.NewRequest("POST", "/injections", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInjectLoanApplication_MissingField(t *testing.T) {
	router := setupRouter()
	jsonBody := []byte(`{
		"participantes": []
	}`)

	req, _ := http.NewRequest("POST", "/injections", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Replicates 500 legacy bug
}

func TestInjectLoanApplication_InvalidFormat(t *testing.T) {
	router := setupRouter()
	jsonBody := []byte(`{
		"datos_credito": {
			"NumeroSolicitud": 123,
			"AntiguedadVivienda": 5,
			"EjecutivoComercial": "Juan",
			"Producto": 1,
			"Objetivo": 2,
			"Destino": 3,
			"MontoAprobado": 100.0,
			"ValorPropiedad": 200.0,
			"FechaAprobacion": "2023-01-01",
			"ValorContado": 100.0,
			"Plazo1": 20,
			"MesesGracia": 2,
			"Tasa1": 3.5,
			"Spread1": 1.0
		},
		"participantes": [
			{
				"Rut": "invalid-rut",
				"TipoParticipacion": 1,
				"Nombre": "A",
				"Paterno": "B",
				"Materno": "C",
				"FechaNacimiento": "01-01-1990"
			}
		]
	}`)

	req, _ := http.NewRequest("POST", "/injections", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Replicates 500 legacy bug
}
