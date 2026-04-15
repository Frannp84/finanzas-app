package main

type Movimiento struct {
	Tipo  string
	Monto float64

	Detalle   string
	Categoria string

	Anio int
	Mes int
	Dia int

	Cuotas int

	MedioPago string

	Cuenta  int
	Tarjeta int
}
type Tarjeta struct {
	Nombre string

	LimiteTotal      float64
	LimiteDisponible float64

	Cierre      int
	Vencimiento int

	CuentaAsociada int
}
