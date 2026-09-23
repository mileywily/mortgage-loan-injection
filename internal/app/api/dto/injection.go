package dto

type DatosCreditoDTO struct {
	NumeroSolicitud    *int     `json:"NumeroSolicitud" binding:"required"`
	AntiguedadVivienda *int     `json:"AntiguedadVivienda" binding:"required"`
	EjecutivoComercial *string  `json:"EjecutivoComercial" binding:"required,max=20"`
	Producto           *int     `json:"Producto" binding:"required"`
	Objetivo           *int     `json:"Objetivo" binding:"required"`
	Destino            *int     `json:"Destino" binding:"required"`
	MontoAprobado      *float64 `json:"MontoAprobado" binding:"required,gt=0"`
	ValorPropiedad     *float64 `json:"ValorPropiedad" binding:"required,gt=0"`
	FechaAprobacion    *string  `json:"FechaAprobacion" binding:"required"`
	ValorContado       *float64 `json:"ValorContado" binding:"required,gt=0"`
	Plazo1             *int     `json:"Plazo1" binding:"required"`
	Plazo2             *int     `json:"Plazo2,omitempty"`
	MesesGracia        *int     `json:"MesesGracia" binding:"required"`
	Tasa1              *float64 `json:"Tasa1" binding:"required,gt=0"`
	Spread1            *float64 `json:"Spread1" binding:"required,gte=0"`
}

type ParticipanteDTO struct {
	Rut               *string `json:"Rut" binding:"required,customRutPattern"`
	TipoParticipacion *int    `json:"TipoParticipacion" binding:"required"`
	Nombre            *string `json:"Nombre" binding:"required,max=100"`
	Paterno           *string `json:"Paterno" binding:"required,max=100"`
	Materno           *string `json:"Materno" binding:"required,max=100"`
	FechaNacimiento   *string `json:"FechaNacimiento" binding:"required,customFechaPattern"`
	Email             *string `json:"Email,omitempty" binding:"omitempty,email,max=100"`
}

type PropiedadDTO struct {
	TipoInmueble *int    `json:"TipoInmueble" binding:"required"`
	Antiguedad   *int    `json:"Antiguedad" binding:"required"`
	Direccion    *string `json:"Direccion" binding:"required,max=100"`
	Numero       *int    `json:"Numero" binding:"required"`
	Depto        *string `json:"Depto,omitempty" binding:"omitempty,max=20"`
	Comuna       *int    `json:"Comuna" binding:"required"`
}

type InyeccionRequestDTO struct {
	DatosCredito *DatosCreditoDTO  `json:"datos_credito" binding:"required"`
	Participantes []ParticipanteDTO `json:"participantes" binding:"required,min=1,dive"`
	Propiedades   []PropiedadDTO    `json:"propiedades,omitempty" binding:"omitempty,dive"`
}

type ErrorResponseDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type InyeccionResponseDTO struct {
	StatusCode      *int    `json:"StatusCode"`
	Mensaje         *string `json:"Mensaje"`
	NumeroSolicitud *int    `json:"NumeroSolicitud"`
}
