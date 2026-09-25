package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mortgage-loan-injection/internal/app/api/dto"
	"mortgage-loan-injection/internal/app/api/mapper"
	"mortgage-loan-injection/internal/core/domain"
)

func TestToDomainDatosCredito(t *testing.T) {
	num := "OPE-123"
	numSol := 123
	antiguedad := 1
	ejecutivo := "Juan"
	producto := 1
	objetivo := 1
	destino := 1
	monto := 100.0
	valor := 200.0
	fecha := "2024-01-01"
	valorC := 10.0
	plazo := 30
	meses := 0
	tasa := 3.5
	spread := 1.0

	d := &dto.DatosCreditoDTO{
		NumeroOperacionOriginal: &num,
		NumeroSolicitud: &numSol,
		AntiguedadVivienda: &antiguedad,
		EjecutivoComercial: &ejecutivo,
		Producto: &producto,
		Objetivo: &objetivo,
		Destino: &destino,
		MontoAprobado: &monto,
		ValorPropiedad: &valor,
		FechaAprobacion: &fecha,
		ValorContado: &valorC,
		Plazo1: &plazo,
		MesesGracia: &meses,
		Tasa1: &tasa,
		Spread1: &spread,
	}
	res := mapper.ToDomainDatosCredito(d)
	assert.NotNil(t, res)
	assert.Equal(t, "OPE-123", *res.NumeroOperacionOriginal)
}

func TestToDomainParticipantes(t *testing.T) {
	rut := "1-9"
	tipo := 1
	nombre := "Test"
	paterno := "P"
	materno := "M"
	fecha := "2000-01-01"

	d := []dto.ParticipanteDTO{
		{Rut: &rut, TipoParticipacion: &tipo, Nombre: &nombre, Paterno: &paterno, Materno: &materno, FechaNacimiento: &fecha},
	}
	res := mapper.ToDomainParticipantes(d)
	assert.Len(t, res, 1)
	assert.Equal(t, "Test", res[0].Nombre)
}

func TestToDomainPropiedades(t *testing.T) {
	comuna := 1
	tipo := 1
	antiguedad := 1
	direccion := "Dir"
	numero := 123

	d := []dto.PropiedadDTO{
		{Comuna: &comuna, TipoInmueble: &tipo, Antiguedad: &antiguedad, Direccion: &direccion, Numero: &numero},
	}
	res := mapper.ToDomainPropiedades(d)
	assert.Len(t, res, 1)
	assert.Equal(t, 1, res[0].Comuna)
}

func TestToDomainInyeccionRequest(t *testing.T) {
	num := "OPE-123"
	numSol := 123
	antiguedad := 1
	ejecutivo := "Juan"
	producto := 1
	objetivo := 1
	destino := 1
	monto := 100.0
	valor := 200.0
	fecha := "2024-01-01"
	valorC := 10.0
	plazo := 30
	meses := 0
	tasa := 3.5
	spread := 1.0

	rut := "1-9"
	tipo := 1
	nombre := "P1"
	paterno := "P"
	materno := "M"
	
	comuna := 1
	direccion := "Dir"
	numero := 123
	
	d := &dto.InyeccionRequestDTO{
		DatosCredito: &dto.DatosCreditoDTO{
			NumeroOperacionOriginal: &num, NumeroSolicitud: &numSol, AntiguedadVivienda: &antiguedad,
			EjecutivoComercial: &ejecutivo, Producto: &producto, Objetivo: &objetivo, Destino: &destino,
			MontoAprobado: &monto, ValorPropiedad: &valor, FechaAprobacion: &fecha, ValorContado: &valorC,
			Plazo1: &plazo, MesesGracia: &meses, Tasa1: &tasa, Spread1: &spread,
		},
		Participantes: []dto.ParticipanteDTO{{Nombre: &nombre, Rut: &rut, TipoParticipacion: &tipo, Paterno: &paterno, Materno: &materno, FechaNacimiento: &fecha}},
		Propiedades:   []dto.PropiedadDTO{{Comuna: &comuna, TipoInmueble: &tipo, Antiguedad: &antiguedad, Direccion: &direccion, Numero: &numero}},
	}
	res := mapper.ToDomainInyeccionRequest(d)
	assert.NotNil(t, res.DatosCredito)
	assert.Len(t, res.Participantes, 1)
	assert.Len(t, res.Propiedades, 1)
}

func TestToDTOInyeccionResponse(t *testing.T) {
	status := 200
	msg := "Success"
	num := 123
	d := domain.InyeccionResponse{
		StatusCode:      &status,
		Mensaje:         &msg,
		NumeroSolicitud: &num,
	}
	res := mapper.ToDTOInyeccionResponse(d)
	assert.Equal(t, 200, *res.StatusCode)
	assert.Equal(t, "Success", *res.Mensaje)
	assert.Equal(t, 123, *res.NumeroSolicitud)
}
