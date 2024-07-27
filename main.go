package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {

	readFileIntoMap("./measurements.txt")

}

type CityTemp struct {
	min     float64
	max     float64
	mean    float64
	visited int
	sum     float64
}

type CityMap = map[string]CityTemp

func readFileIntoMap(filepath string) CityMap {
	mapOfTemp := make(CityMap)

	f, err := os.Open(filepath)

	if err != nil {
		log.Fatal("error to read file", err)
	}

	r := bufio.NewReader(f)

	i := 0
	for {
		i++
		line, err := r.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				return mapOfTemp
			}
			log.Fatal("some error occured", err)
		}

		index := strings.Index(line, ";")
		city := line[:index]
		temp := line[index+1 : len(line)-1]

		num, err := strconv.ParseFloat(temp, 64)

		if i%10000 == 0 {
			log.Println("at count", i)
		}

		if err != nil {
			log.Fatal("string conversion failed", err)

		}

		if _, ok := mapOfTemp[city]; ok {

			current := mapOfTemp[city]

			mapOfTemp[city] = CityTemp{
				min:     math.Min(current.min, num),
				max:     math.Max(current.max, num),
				visited: current.visited + 1,
				mean:    current.sum/float64(current.visited) + 1,
				sum:     current.sum + num,
			}

		} else {
			mapOfTemp[city] = CityTemp{
				min:     math.Inf(1),
				max:     math.Inf(-1),
				mean:    num,
				visited: 1,
				sum:     num,
			}
		}
	}

}
