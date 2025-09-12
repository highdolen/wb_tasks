package main

import "fmt"

func main() {
	temp := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	m := make(map[int][]float64)

	for _, v := range temp {
		switch {
		case v < -20:
			m[-20] = append(m[-20], v)
		case v > 10 && v < 20:
			m[10] = append(m[10], v)
		case v > 20 && v < 30:
			m[20] = append(m[20], v)
		case v > 30:
			m[30] = append(m[30], v)
		}
	}
	fmt.Println(m)
}
