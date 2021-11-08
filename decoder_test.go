package main

import (
	"math"
	"strings"
	"testing"
)

var satTest = []Satellite{
	{
		Name:     "kenobi",
		Distance: 0,
		Message:  []string{"", "este", "es", "un", "mensaje"},
		Position: Point{X: -500.0, Y: -200.0},
	},
	{
		Name:     "skywalker",
		Distance: 0,
		Message:  []string{"este", "", "un", "mensaje"},
		Position: Point{X: 100.0, Y: -100.0},
	},
	{
		Name:     "sato",
		Distance: 0,
		Message:  []string{"", "", "es", "", "", "mensaje", ""},
		Position: Point{X: 500.0, Y: 100.0},
	},
}

func TestGetLocation(t *testing.T) {
	// Coordenadas del punto que se necesita encontrar
	control := Point{X: 400.0, Y: -200.0}
	for key, sat := range satTest {
		// Se calcula la distancia del punto a los satelites
		satTest[key].Distance = math.Sqrt(math.Pow(sat.Position.X-control.X, 2) + math.Pow(sat.Position.Y-control.Y, 2))

	}
	// Se llama la funcion GetLocation para evaluar la efectividad de la funcion
	point, err := GetLocation(satTest)
	// Validaciones de excepciones
	if err != nil {
		t.Errorf(err.Error())
	}
	// Validaciones de errores de logica
	if point.X != control.X || point.Y != control.Y {
		t.Errorf("La Triletaracion no se calculo de forma correcta, se esperaban las coordenadas: %v, %v y se obtuvo %v, %v", control.X, control.Y, point.X, point.Y)
	}

}
func TestGetMessage(t *testing.T) {
	var control string
	control = "este es un mensaje"

	message1, err1 := GetMessage(satTest, "flexible")
	if err1 != nil {
		t.Errorf(err1.Error())
	}
	if message1 != control {
		t.Errorf("No se pudo completar el mensaje, se esperaba: %s y se obtuvo %s", control, message1)
	}

	control = "Hola mi nombre es diego fernando prieto y estoy testeando el correcto funcionamiento de esta aplicacion"
	contArray := strings.Split(control, " ")
	distortion := 1
	for key, _ := range satTest {
		var badMessage []string
		for _, word := range contArray {
			if distortion == 2 {
				//distortion = 0
				badMessage = append(badMessage, word)
			} else {
				distortion = distortion + 1
				badMessage = append(badMessage, "")
			}
		}
		// Se calcula la distancia del punto a los satelites
		satTest[key].Message = badMessage
	}

	message2, err2 := GetMessage(satTest, "strict")
	if err2 != nil {
		t.Errorf(err2.Error())
	}
	if message2 != control {
		t.Errorf("No se pudo completar el mensaje, se esperaba: %s y se obtuvo %s", control, message2)
	}

}
