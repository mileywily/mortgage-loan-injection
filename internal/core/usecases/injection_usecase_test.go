package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mortgage-loan-injection/internal/core/domain"
	"mortgage-loan-injection/internal/core/usecases"
)

type mockClient struct {
	shouldFail bool
}

func (m *mockClient) Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error) {
	if m.shouldFail {
		return domain.InyeccionResponse{}, errors.New("mock error")
	}
	msg := "Success"
	status := 200
	return domain.InyeccionResponse{
		Mensaje:    &msg,
		StatusCode: &status,
	}, nil
}

func TestProcessInjection_Valid(t *testing.T) {
	uc := usecases.NewInjectionUseCaseImpl(&mockClient{})
	req := domain.InyeccionRequest{
		DatosCredito: domain.DatosCredito{
			MontoAprobado:  100,
			ValorPropiedad: 200,
			Plazo1:         30,
			Tasa1:          3.5,
		},
		Participantes: []domain.Participante{
			{Nombre: "Test"},
		},
	}
	resp, err := uc.ProcessInjection(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestProcessInjection_FailClient(t *testing.T) {
	uc := usecases.NewInjectionUseCaseImpl(&mockClient{shouldFail: true})
	req := domain.InyeccionRequest{
		DatosCredito: domain.DatosCredito{
			MontoAprobado:  100,
			ValorPropiedad: 200,
			Plazo1:         30,
			Tasa1:          3.5,
		},
		Participantes: []domain.Participante{
			{Nombre: "Test"},
		},
	}
	_, err := uc.ProcessInjection(context.Background(), req)
	assert.Error(t, err)
}

func TestValidateBusinessRules(t *testing.T) {
	uc := usecases.NewInjectionUseCaseImpl(&mockClient{})

	tests := []struct {
		name        string
		req         domain.InyeccionRequest
		expectedErr string
	}{
		{
			name:        "Missing Participantes",
			req:         domain.InyeccionRequest{},
			expectedErr: "Missing participantes field",
		},
		{
			name: "MontoAprobado > ValorPropiedad",
			req: domain.InyeccionRequest{
				Participantes: []domain.Participante{{Nombre: "A"}},
				DatosCredito:  domain.DatosCredito{MontoAprobado: 200, ValorPropiedad: 100},
			},
			expectedErr: "MontoAprobado with invalid value",
		},
		{
			name: "Plazo Invalido",
			req: domain.InyeccionRequest{
				Participantes: []domain.Participante{{Nombre: "A"}},
				DatosCredito:  domain.DatosCredito{MontoAprobado: 100, ValorPropiedad: 200, Plazo1: 17},
			},
			expectedErr: "Plazo1 with invalid value",
		},
		{
			name: "Tasa Invalida",
			req: domain.InyeccionRequest{
				Participantes: []domain.Participante{{Nombre: "A"}},
				DatosCredito:  domain.DatosCredito{MontoAprobado: 100, ValorPropiedad: 200, Plazo1: 20, Tasa1: 60},
			},
			expectedErr: "Tasa1 with invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.ProcessInjection(context.Background(), tt.req)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}
