package mapper

import (
	"mortgage-loan-injection/internal/app/api/dto"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToDomainDatosCredito(t *testing.T) {
	num := 123
	ant := 5
	ejec := "Juan"
	prod := 1
	obj := 2
	dest := 3
	monto := 100.5
	val := 200.5
	fecha := "2023-01-01"
	cont := 100.0
	plazo1 := 10
	meses := 2
	tasa := 3.5
	spread := 1.0

	dtoData := &dto.DatosCreditoDTO{
		NumeroSolicitud:    &num,
		AntiguedadVivienda: &ant,
		EjecutivoComercial: &ejec,
		Producto:           &prod,
		Objetivo:           &obj,
		Destino:            &dest,
		MontoAprobado:      &monto,
		ValorPropiedad:     &val,
		FechaAprobacion:    &fecha,
		ValorContado:       &cont,
		Plazo1:             &plazo1,
		MesesGracia:        &meses,
		Tasa1:              &tasa,
		Spread1:            &spread,
	}

	domainData := ToDomainDatosCredito(dtoData)
	assert.Equal(t, num, domainData.NumeroSolicitud)
	assert.Equal(t, ant, domainData.AntiguedadVivienda)
}

func TestToDomainDatosCredito_Nil(t *testing.T) {
	domainData := ToDomainDatosCredito(nil)
	assert.Equal(t, 0, domainData.NumeroSolicitud)
}
