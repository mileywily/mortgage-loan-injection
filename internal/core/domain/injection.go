package domain

type DatosCredito struct {
	NumeroSolicitud    int
	AntiguedadVivienda int
	EjecutivoComercial string
	Producto           int
	Objetivo           int
	Destino            int
	MontoAprobado      float64
	ValorPropiedad     float64
	FechaAprobacion    string
	ValorContado       float64
	Plazo1             int
	Plazo2             *int
	MesesGracia        int
	Tasa1              float64
	Spread1            float64

	NumeroOperacionOriginal    *string  `json:"NumeroOperacionOriginal,omitempty"`
	MontoAbono                 *float64 `json:"MontoAbono,omitempty"`
	IncluirGastosOperacionales *int     `json:"IncluirGastosOperacionales,omitempty"`
	TipoGarantia               *int     `json:"TipoGarantia,omitempty"`
	NumeroRenegociacion        *int     `json:"NumeroRenegociacion,omitempty"`
	SubsidioOriginal           *int     `json:"SubsidioOriginal,omitempty"`
	ReduccionInteresesYGC      *int     `json:"ReduccionInteresesYGC,omitempty"`
}

type Participante struct {
	Rut               string
	TipoParticipacion int
	Nombre            string
	Paterno           string
	Materno           string
	FechaNacimiento   string
	Email             *string
}

type Propiedad struct {
	TipoInmueble int
	Antiguedad   int
	Direccion    string
	Numero       int
	Depto        *string
	Comuna       int
}

type InyeccionRequest struct {
	DatosCredito  DatosCredito
	Participantes []Participante
	Propiedades   []Propiedad
}

type InyeccionResponse struct {
	StatusCode      *int
	Mensaje         *string
	NumeroSolicitud *int
}
