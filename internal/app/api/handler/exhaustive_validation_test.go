package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"mortgage-loan-injection/internal/app/api/handler"
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

func setupExhaustiveRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
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
	h := handler.NewInjectionHandler(uc)
	router.POST("/injections", h.InjectLoanApplication)
	return router
}

func executeRequest(router *gin.Engine, bodyStr string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/injections", bytes.NewBuffer([]byte(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func getValidJSON() map[string]interface{} {
	return map[string]interface{}{
		"datos_credito": map[string]interface{}{
			"NumeroSolicitud":    192269,
			"AntiguedadVivienda": 1,
			"EjecutivoComercial": "PSOTO",
			"Producto":           1,
			"Objetivo":           1,
			"Destino":            1,
			"FechaAprobacion":    "2024-01-15",
			"MontoAprobado":      1688.09,
			"ValorPropiedad":     2300.00,
			"ValorContado":       261.91,
			"Plazo1":             30,
			"MesesGracia":        0,
			"Tasa1":              3.99,
			"Spread1":            0.0,
		},
		"participantes": []map[string]interface{}{
			{
				"Rut":               "18.015.051-K",
				"TipoParticipacion": 1,
				"Nombre":            "Pedro",
				"Paterno":           "Perez",
				"Materno":           "Perez",
				"FechaNacimiento":   "04-07-1992",
				"Email":             "pedro.perez@email.com",
			},
		},
	}
}

// ---------------------------------------------------------
// 1. DatosCreditoValidationTest
// ---------------------------------------------------------

func TestLegacy_DatosCredito_Valid(t *testing.T) {
	router := setupExhaustiveRouter()
	validBody, _ := json.Marshal(getValidJSON())
	w := executeRequest(router, string(validBody))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLegacy_DatosCredito_MissingNumeroSolicitud(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	delete(body["datos_credito"].(map[string]interface{}), "NumeroSolicitud")
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Missing NumeroSolicitud field")
}

func TestLegacy_DatosCredito_EjecutivoComercialTooLong(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	body["datos_credito"].(map[string]interface{})["EjecutivoComercial"] = "EJECUTIVO_COMERCIAL_MUY_LARGO"
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "EjecutivoComercial with invalid value")
}

func TestLegacy_DatosCredito_MontoAprobadoZeroOrNegative(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	
	// Zero
	body["datos_credito"].(map[string]interface{})["MontoAprobado"] = 0.0
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "MontoAprobado with invalid value")

	// Negative
	body["datos_credito"].(map[string]interface{})["MontoAprobado"] = -100.0
	b, _ = json.Marshal(body)
	w = executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "MontoAprobado with invalid value")
}

// ---------------------------------------------------------
// 2. BusinessRulesValidationTest
// ---------------------------------------------------------

func TestLegacy_Business_MontoAprobadoGreaterThanValorPropiedad(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	body["datos_credito"].(map[string]interface{})["MontoAprobado"] = 3000.0 // > 2300.0
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "MontoAprobado with invalid value")
}

func TestLegacy_Business_PlazoTooLowOrHigh(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	
	// Low
	body["datos_credito"].(map[string]interface{})["Plazo1"] = 3
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Plazo1 with invalid value")

	// High
	body["datos_credito"].(map[string]interface{})["Plazo1"] = 35
	b, _ = json.Marshal(body)
	w = executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Plazo1 with invalid value")

	// Not multiple of 5
	body["datos_credito"].(map[string]interface{})["Plazo1"] = 17
	b, _ = json.Marshal(body)
	w = executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Plazo1 with invalid value")
}

func TestLegacy_Business_TasaInvalid(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	
	// Low
	body["datos_credito"].(map[string]interface{})["Tasa1"] = 0.05
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Tasa1 with invalid value")

	// High
	body["datos_credito"].(map[string]interface{})["Tasa1"] = 55.0
	b, _ = json.Marshal(body)
	w = executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Tasa1 with invalid value")
}

// ---------------------------------------------------------
// 3. ParticipanteValidationTest
// ---------------------------------------------------------

func TestLegacy_Participante_RutInvalidFormat(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	
	invalidRuts := []string{"12345678-9", "12.345.6789", "A.345.678-9", "12.345.678-Z"}
	for _, rut := range invalidRuts {
		body["participantes"].([]map[string]interface{})[0]["Rut"] = rut
		b, _ := json.Marshal(body)
		w := executeRequest(router, string(b))
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Legacy squash
		assert.Contains(t, w.Body.String(), "Rut with invalid format")
	}
}

func TestLegacy_Participante_EmailInvalidFormat(t *testing.T) {
	router := setupExhaustiveRouter()
	body := getValidJSON()
	body["participantes"].([]map[string]interface{})[0]["Email"] = "correo-invalido"
	b, _ := json.Marshal(body)
	w := executeRequest(router, string(b))
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Legacy squash
	assert.Contains(t, w.Body.String(), "Email with invalid format")
}
