package usecases

import (
	"context"
	"mortgage-loan-injection/internal/core/domain"
	"testing"
	"github.com/stretchr/testify/assert"
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

func TestProcessInjection_Success(t *testing.T) {
	uc := NewInjectionUseCaseImpl(&mockFinnFlowClient{})
	req := domain.InyeccionRequest{
		DatosCredito: domain.DatosCredito{
			NumeroSolicitud: 12345,
			MontoAprobado: 100,
			ValorPropiedad: 200,
			Plazo1: 20,
			Tasa1: 3.5,
		},
		Participantes: []domain.Participante{
			{Rut: "1.234.567-8"},
		},
	}

	res, err := uc.ProcessInjection(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, 200, *res.StatusCode)
	assert.Equal(t, 12345, *res.NumeroSolicitud)
}

func TestProcessInjection_BusinessRulesFail(t *testing.T) {
	uc := NewInjectionUseCaseImpl(&mockFinnFlowClient{})
	
	// Fails: Plazo1 invalid
	req := domain.InyeccionRequest{
		DatosCredito: domain.DatosCredito{
			NumeroSolicitud: 12345,
			MontoAprobado: 100,
			ValorPropiedad: 200,
			Plazo1: 2, // invalid
			Tasa1: 3.5,
		},
		Participantes: []domain.Participante{
			{Rut: "1.234.567-8"},
		},
	}

	_, err := uc.ProcessInjection(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "Plazo1 with invalid value", err.Error())
}
