package mapper

import (
	"mortgage-loan-injection/internal/app/api/dto"
	"mortgage-loan-injection/internal/core/domain"
)

func ToDomainDatosCredito(dtoDatos *dto.DatosCreditoDTO) domain.DatosCredito {
	if dtoDatos == nil {
		return domain.DatosCredito{}
	}
	return domain.DatosCredito{
		NumeroSolicitud:    *dtoDatos.NumeroSolicitud,
		AntiguedadVivienda: *dtoDatos.AntiguedadVivienda,
		EjecutivoComercial: *dtoDatos.EjecutivoComercial,
		Producto:           *dtoDatos.Producto,
		Objetivo:           *dtoDatos.Objetivo,
		Destino:            *dtoDatos.Destino,
		MontoAprobado:      *dtoDatos.MontoAprobado,
		ValorPropiedad:     *dtoDatos.ValorPropiedad,
		FechaAprobacion:    *dtoDatos.FechaAprobacion,
		ValorContado:       *dtoDatos.ValorContado,
		Plazo1:             *dtoDatos.Plazo1,
		Plazo2:             dtoDatos.Plazo2,
		MesesGracia:        *dtoDatos.MesesGracia,
		Tasa1:              *dtoDatos.Tasa1,
		Spread1:            *dtoDatos.Spread1,
	}
}

func ToDomainParticipantes(dtoPart []dto.ParticipanteDTO) []domain.Participante {
	if dtoPart == nil {
		return nil
	}
	var partes []domain.Participante
	for _, p := range dtoPart {
		partes = append(partes, domain.Participante{
			Rut:               *p.Rut,
			TipoParticipacion: *p.TipoParticipacion,
			Nombre:            *p.Nombre,
			Paterno:           *p.Paterno,
			Materno:           *p.Materno,
			FechaNacimiento:   *p.FechaNacimiento,
			Email:             p.Email,
		})
	}
	return partes
}

func ToDomainPropiedades(dtoProp []dto.PropiedadDTO) []domain.Propiedad {
	if dtoProp == nil {
		return nil
	}
	var props []domain.Propiedad
	for _, p := range dtoProp {
		props = append(props, domain.Propiedad{
			TipoInmueble: *p.TipoInmueble,
			Antiguedad:   *p.Antiguedad,
			Direccion:    *p.Direccion,
			Numero:       *p.Numero,
			Depto:        p.Depto,
			Comuna:       *p.Comuna,
		})
	}
	return props
}

func ToDomainInyeccionRequest(req *dto.InyeccionRequestDTO) domain.InyeccionRequest {
	return domain.InyeccionRequest{
		DatosCredito:  ToDomainDatosCredito(req.DatosCredito),
		Participantes: ToDomainParticipantes(req.Participantes),
		Propiedades:   ToDomainPropiedades(req.Propiedades),
	}
}

func ToDTOInyeccionResponse(res domain.InyeccionResponse) dto.InyeccionResponseDTO {
	return dto.InyeccionResponseDTO{
		StatusCode:      res.StatusCode,
		Mensaje:         res.Mensaje,
		NumeroSolicitud: res.NumeroSolicitud,
	}
}
