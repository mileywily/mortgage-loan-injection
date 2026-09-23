package usecases

import (
	"context"
	"errors"
	"log/slog"
	"mortgage-loan-injection/internal/core/domain"
)

type InjectionUseCase interface {
	ProcessInjection(ctx context.Context, req domain.InyeccionRequest) (domain.InyeccionResponse, error)
}

type injectionUseCaseImpl struct {
	finnflowClient interface {
		Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error)
	}
}

func NewInjectionUseCaseImpl(client interface {
	Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error)
}) InjectionUseCase {
	return &injectionUseCaseImpl{
		finnflowClient: client,
	}
}

func (u *injectionUseCaseImpl) ProcessInjection(ctx context.Context, req domain.InyeccionRequest) (domain.InyeccionResponse, error) {
	slog.Debug("Starting injection process", "logger", "cl.bancofalabella.mortgage.injection.service.InjectionService", "thread", "http-nio-8080-exec-1")
	if err := u.validateBusinessRules(req); err != nil {
		return domain.InyeccionResponse{}, err
	}

	// Call the external FinnFlow API
	resp, err := u.finnflowClient.Inject(req)
	if err == nil {
		msgLog := "null"
		if resp.Mensaje != nil {
			msgLog = *resp.Mensaje
		}
		slog.Debug("Injection completed with response: "+msgLog, "logger", "cl.bancofalabella.mortgage.injection.service.InjectionService", "thread", "http-nio-8080-exec-1")
	}
	return resp, err
}

func (u *injectionUseCaseImpl) validateBusinessRules(req domain.InyeccionRequest) error {
	// Validar que haya al menos un participante (ya cubierto por validation de DTO, pero lo replicamos a nivel dominio)
	if len(req.Participantes) == 0 {
		return errors.New("Missing participantes field")
	}

	// Validar que el monto aprobado no sea mayor al valor de la propiedad
	if req.DatosCredito.MontoAprobado > req.DatosCredito.ValorPropiedad {
		return errors.New("MontoAprobado with invalid value")
	}

	// Validar que el plazo sea válido (5, 10, 15, 20, 25, 30 años)
	plazo := req.DatosCredito.Plazo1
	if plazo < 5 || plazo > 30 || plazo%5 != 0 {
		return errors.New("Plazo1 with invalid value")
	}

	// Validar que la tasa sea razonable (entre 0.1% y 50%)
	tasa1 := req.DatosCredito.Tasa1
	if tasa1 < 0.1 || tasa1 > 50.0 {
		return errors.New("Tasa1 with invalid value")
	}

	return nil
}
