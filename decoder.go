package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

func ResetFleet() {
	for i, _ := range fleets {
		fleets[i].Distance = 0
		fleets[i].Message = []string{}
	}
}
func GetLocation(distances []Satellite) (Point, error) {
	// Se recorren los Satellites para matchear los datos de coordenadas (Dato conocido) y la distancia (Dato proveniente del emisor)
	ResetFleet()
	for _, transmitter := range distances {
		for key, sat := range fleets {
			// Se valida que la informacion de los satelites este completa
			if transmitter.Name == sat.Name {
				fleets[key].Distance = transmitter.Distance
			}
		}
	}
	// Una vez asignada la distancia y la posicion de los satelites se calculan las coordenadas usando triletaracion
	point, err := getTrilateration(fleets)

	return point, err
	//return MatchPoints(areas)
}

func GetMessage(messages []Satellite, mode string) (string, error) {
	normMessages, max := normaliceMessage(messages)

	words := make(map[int]string)
	for _, item := range normMessages {
		for index, mess := range item.Message {
			if mess != "" {
				_, ok := words[index]
				if !ok {
					words[index] = mess
				}
			}

		}
	}
	flatString := []string{}
	for i := 0; i < max; i++ {
		if words[i] != "" {
			flatString = append(flatString, words[i])
		} else if mode == "strict" {
			return "Mensaje no determinado", fmt.Errorf("El mensaje no se pudo determinar")
		}
	}

	str_concat := strings.Join(flatString, " ")
	return str_concat, nil
}
func normaliceMessage(messages []Satellite) ([]Satellite, int) {

	var max, proto, key int
	var val bool
	for index, m := range messages {
		if proto == 0 {
			proto = len(m.Message)
		}
		if len(m.Message) > max {
			max = len(m.Message)
			key = index
		}
		if proto != len(m.Message) {
			val = false
		}
	}
	if !val {
		for index, m := range messages {
			if len(m.Message) != len(messages[key].Message) {
				var keyDif int
				var normMess = make(map[int]string)
				for inCtr, ctr := range messages[key].Message {
					if ctr != "" {
						for inCom, w := range m.Message {
							if ctr == w {
								keyDif = inCtr - inCom
							}
						}
					}
				}
				for k, w := range m.Message {
					if k+keyDif >= 0 {
						normMess[k+keyDif] = w
					}
				}
				flatString := []string{}

				for i := 0; i < len(normMess)+keyDif; i++ {
					flatString = append(flatString, normMess[i])
				}

				messages[index].Message = flatString
			}
		}
	}
	return messages, max
}
func getTrilateration(transmitters []Satellite) (Point, error) {
	for _, sat := range fleets {
		// Se valida que la informacion de los satelites este completa
		if sat.Distance == 0 {
			err := errors.New(fmt.Sprintf("Distancia imposible de determinada en el satelite: %v", sat.Name))
			return Point{}, err
		}
	}
	// Se declaran las variables mat.dense para los ejes x,y
	var subX, subY, divX, divY, P3P1iex, triPT mat.Dense

	// Se declaran las matrices para operar con la libreria mat
	P1 := mat.NewDense(2, 1, []float64{transmitters[0].Position.X, transmitters[0].Position.Y})
	P2 := mat.NewDense(2, 1, []float64{transmitters[1].Position.X, transmitters[1].Position.Y})
	P3 := mat.NewDense(2, 1, []float64{transmitters[2].Position.X, transmitters[2].Position.Y})

	// Utilizando formulas de algebra lineal se resuelven las ecuaciones de los circulos
	subX.Sub(P2, P1)
	normX := mat.NewDense(2, 1, []float64{mat.Norm(&subX, 2), mat.Norm(&subX, 2)})
	divX.DivElem(&subX, normX)
	ex := blas64.Vector{N: 2, Data: []float64{divX.RawMatrix().Data[0], divX.RawMatrix().Data[1]}, Inc: 1}

	subY.Sub(P3, P1)
	miniusP3P1 := blas64.Vector{N: 2, Data: []float64{subY.RawMatrix().Data[0], subY.RawMatrix().Data[1]}, Inc: 1}
	i := blas64.Dot(ex, miniusP3P1)
	blas64.Scal(i, ex)

	iex := mat.NewDense(2, 1, []float64{ex.Data[0], ex.Data[1]})
	P3P1iex.Sub(&subY, iex)

	normY := mat.NewDense(2, 1, []float64{mat.Norm(&P3P1iex, 2), mat.Norm(&P3P1iex, 2)})
	divY.DivElem(&P3P1iex, normY)
	ey := blas64.Vector{N: 2, Data: []float64{divY.RawMatrix().Data[0], divY.RawMatrix().Data[1]}, Inc: 1}

	d := mat.Norm(&subX, 2)
	j := blas64.Dot(ey, miniusP3P1)
	// Se Calculan las coordendas usando la trilateracion 2D
	x := (math.Pow(transmitters[0].Distance, 2) - math.Pow(transmitters[1].Distance, 2) + math.Pow(d, 2)) / (2 * d)
	y := ((math.Pow(transmitters[0].Distance, 2) - math.Pow(transmitters[2].Distance, 2) + math.Pow(i, 2) + math.Pow(j, 2)) / (2 * j)) - ((i / j) * x)
	ex = blas64.Vector{N: 2, Data: []float64{divX.RawMatrix().Data[0], divX.RawMatrix().Data[1]}, Inc: 1}
	blas64.Scal(x, ex)
	blas64.Scal(y, ey)

	iex = mat.NewDense(2, 1, []float64{ex.Data[0], ex.Data[1]})
	iey := mat.NewDense(2, 1, []float64{ey.Data[0], ey.Data[1]})
	// Se obtienen las coordenadas del punto
	triPT.Add(P1, iex)
	triPT.Add(&triPT, iey)

	result := Point{(math.Round(triPT.RawMatrix().Data[0]*10) / 10), (math.Round(triPT.RawMatrix().Data[1]*10) / 10)}
	return result, nil
}
