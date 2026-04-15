package main

func nombreMes(mes int) string {

    meses := []string{
        "Enero",
        "Febrero",
        "Marzo",
        "Abril",
        "Mayo",
        "Junio",
        "Julio",
        "Agosto",
        "Septiembre",
        "Octubre",
        "Noviembre",
        "Diciembre",
    }

    if mes < 1 || mes > 12 {

        return "Mes invalido"

    }

    return meses[mes-1]

}
func calcularMesResumen(compraMes int, diaCompra int, cierre int) int {

    if diaCompra > cierre {
        return compraMes + 1
    }

    return compraMes
}