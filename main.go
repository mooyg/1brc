package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {

	m := readFileIntoMap("./measurements.txt")
	log.Println(len(m))

}

type CityTemp struct {
	min     float64
	max     float64
	mean    float64
	visited int
	sum     float64
}

type CityMap = map[string]CityTemp

type ResultChanTemp struct {
	city string
	temp float64
}

func readFileIntoMap(filepath string) CityMap {
	mapOfTemp := make(chan CityMap, 1)
	mapOfTemp <- make(CityMap)

	f, err := os.Open(filepath)

	if err != nil {
		log.Fatal("error to read file", err)
	}

	r := bufio.NewReader(f)

	var wg sync.WaitGroup

	i := 0

	lineChan := make(chan string)

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				city, temp := parseLine(line)
				m := <-mapOfTemp
				updateMap(m, city, temp)
				mapOfTemp <- m
			}
		}()
	}
	for {
		i++

		line, err := r.ReadString('\n')

		if i%10000 == 0 {
			log.Println(i)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("some error occured", err)
		}

		lineChan <- line
	}
	close(lineChan)

	wg.Wait()
	return <-mapOfTemp

}

func parseLine(line string) (string, float64) {
	index := strings.Index(line, ";")
	city := line[:index]
	temp := line[index+1 : len(line)-1]

	num, err := strconv.ParseFloat(temp, 64)

	if err != nil {
		log.Fatal("string conversion failed", err)

	}
	return city, num
}

func updateMap(mapOfTemp CityMap, city string, num float64) {
	if current, ok := mapOfTemp[city]; ok {
		mapOfTemp[city] = CityTemp{
			min:     math.Min(current.min, num),
			max:     math.Max(current.max, num),
			visited: current.visited + 1,
			mean:    (current.sum + num) / float64(current.visited+1),
			sum:     current.sum + num,
		}
	} else {
		mapOfTemp[city] = CityTemp{
			min:     num,
			max:     num,
			mean:    num,
			visited: 1,
			sum:     num,
		}
	}
}
