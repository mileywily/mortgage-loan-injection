package handler

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"mortgage-loan-injection/internal/app/api/dto"
	"mortgage-loan-injection/internal/app/api/mapper"
	"mortgage-loan-injection/internal/core/usecases"
)

type InjectionHandler struct {
	useCase usecases.InjectionUseCase
}

func NewInjectionHandler(uc usecases.InjectionUseCase) *InjectionHandler {
	return &InjectionHandler{useCase: uc}
}

func (h *InjectionHandler) InjectLoanApplication(c *gin.Context) {
	// Replicando log exacto del controller legado
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	rawRequest := string(bodyBytes)
	slog.Info("Received injection request: " + rawRequest, "logger", "cl.bancofalabella.mortgage.injection.controller.InjectionController", "thread", "http-nio-8080-exec-1")
	
	// Restaurar el body para Gin
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqDTO dto.InyeccionRequestDTO

	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		h.handleValidationError(c, err)
		return
	}

	// Mapear a dominio
	domainReq := mapper.ToDomainInyeccionRequest(&reqDTO)

	// Ejecutar caso de uso
	domainRes, err := h.useCase.ProcessInjection(c.Request.Context(), domainReq)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	resDTO := mapper.ToDTOInyeccionResponse(domainRes)
	
	solicitud := "null"
	if resDTO.NumeroSolicitud != nil {
		solicitud = fmt.Sprintf("%d", *resDTO.NumeroSolicitud)
	}
	statusLog := "null"
	if resDTO.StatusCode != nil {
		statusLog = fmt.Sprintf("%d", *resDTO.StatusCode)
	}
	msgLog := "null"
	if resDTO.Mensaje != nil {
		msgLog = *resDTO.Mensaje
	}

	slog.Info(fmt.Sprintf("Injection processed for solicitud: %s with status: %s and message: %s", solicitud, statusLog, msgLog), "logger", "cl.bancofalabella.mortgage.injection.controller.InjectionController", "thread", "http-nio-8080-exec-1")
	
	// Si el status code viene del response y es >= 400
	if resDTO.StatusCode != nil && *resDTO.StatusCode >= 400 {
		c.JSON(*resDTO.StatusCode, resDTO)
		return
	}

	c.JSON(http.StatusOK, resDTO)
}

func (h *InjectionHandler) handleValidationError(c *gin.Context, err error) {
	// HttpMessageNotReadableException equivalent
	if err.Error() == "EOF" || regexp.MustCompile(`cannot unmarshal`).MatchString(err.Error()) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDTO{
			Code:    "missing_field",
			Message: "Request body is required",
		})
		return
	}

	var fieldName string
	var code, message string
	var status int

	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		firstErr := errs[0]
		fieldName = firstErr.Field()
		tag := firstErr.Tag()

		if tag == "customRutPattern" || tag == "customFechaPattern" || tag == "email" {
			status = http.StatusPaymentRequired // 402
			code = "invalid_format"
			message = fmt.Sprintf("%s with invalid format", fieldName)
		} else if tag == "gt" || tag == "gte" || tag == "min" || tag == "max" {
			status = http.StatusNotAcceptable // 406
			code = "invalid_value"
			message = fmt.Sprintf("%s with invalid value", fieldName)
		} else {
			status = http.StatusBadRequest // 400
			code = "missing_field"
			message = fmt.Sprintf("Missing %s field", fieldName)
		}
	} else {
		status = http.StatusBadRequest
		code = "bad_request"
		message = "Invalid request format"
	}
	
	// REPLICATE LEGACY BUG: Squash to 500 with the exact String that Java produced when failing to parse FinnFlow's error
	statusText := ""
	switch status {
	case 400: statusText = "400 BAD_REQUEST"
	case 402: statusText = "402 PAYMENT_REQUIRED"
	case 406: statusText = "406 NOT_ACCEPTABLE"
	default: statusText = fmt.Sprintf("%d ERROR", status)
	}

	legacyMensaje := fmt.Sprintf(`FinnFlow injection failed: %s - {"code":"%s","message":"%s"}`, statusText, code, message)
	
	var nilSolicitud *int
	legacyStatus := 500
	
	c.JSON(http.StatusInternalServerError, dto.InyeccionResponseDTO{
		StatusCode:      &legacyStatus,
		Mensaje:         &legacyMensaje,
		NumeroSolicitud: nilSolicitud,
	})
}

func (h *InjectionHandler) handleDomainError(c *gin.Context, err error) {
	msg := err.Error()
	
	var status int
	var code string

	if msg == "Missing participantes field" || msg == "MontoAprobado with invalid value" || msg == "Plazo1 with invalid value" || msg == "Tasa1 with invalid value" {
		if msg == "Missing participantes field" {
			status = http.StatusBadRequest
			code = "bad_request"
		} else if regexp.MustCompile(`invalid value`).MatchString(msg) {
			status = http.StatusNotAcceptable // 406
			code = "invalid_value"
		} else if regexp.MustCompile(`invalid format`).MatchString(msg) {
			status = http.StatusPaymentRequired // 402
			code = "invalid_format"
		}

		statusText := ""
		switch status {
		case 400: statusText = "400 BAD_REQUEST"
		case 402: statusText = "402 PAYMENT_REQUIRED"
		case 406: statusText = "406 NOT_ACCEPTABLE"
		default: statusText = fmt.Sprintf("%d ERROR", status)
		}
		
		legacyMensaje := fmt.Sprintf(`FinnFlow injection failed: %s - {"code":"%s","message":"%s"}`, statusText, code, msg)
		var nilSolicitud *int
		legacyStatus := 500
		
		c.JSON(http.StatusInternalServerError, dto.InyeccionResponseDTO{
			StatusCode:      &legacyStatus,
			Mensaje:         &legacyMensaje,
			NumeroSolicitud: nilSolicitud,
		})
		return
	}
	
	if msg == "token_not_valid" {
		legacyMensaje := "Token is not valid"
		legacyStatus := 401
		c.JSON(http.StatusUnauthorized, dto.InyeccionResponseDTO{
			StatusCode: &legacyStatus,
			Mensaje: &legacyMensaje,
		})
		return
	}

	legacyMensaje := "Internal server error"
	legacyStatus := 500
	c.JSON(http.StatusInternalServerError, dto.InyeccionResponseDTO{
		StatusCode: &legacyStatus,
		Mensaje: &legacyMensaje,
	})
}

