package main

import (
	// "encoding/binary"
	"bufio"
	// "container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var numOfRows int
var rowSize int
var hamInts []uint64

type Vertex struct {
	val       uint64
	leader    *Vertex
	followers []*Vertex
}

func (v Vertex) String() string {

	var followers []uint64

	for _, x := range v.followers {
		followers = append(followers, x.val)
	}

	return fmt.Sprintf("\nval:\t%v\nLeader:\t%v\nFollowerr:\t%v\n\n\n",
		v.val, v.leader.val, followers)
}

var Vertices = make(map[uint64]*Vertex, numOfRows)
var LeaderMap = make(map[uint64]*Vertex, numOfRows)

// SetShortestEdge sets the shortest edge for v
func (v *Vertex) ConsumeFollowers(w *Vertex) {

	w.leader = v

	for _, x := range w.followers {
		x.leader = v
		v.followers = append(v.followers, x)
	}

	w.followers = nil

	delete(LeaderMap, w.val)
}

func main() {

	readFile(os.Args[1])

	TwoDegreeHamInts(rowSize)

	cluster()

	// fmt.Println(Vertices)

	fmt.Println(len(LeaderMap))

}

func readFile(filename string) {

	file, err := os.Open(filename) //should read in file named in CLI

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan first line
	if scanner.Scan() {

		firstLine := strings.Fields(scanner.Text())

		numOfRows, err = strconv.Atoi(firstLine[0])

		rowSize, err = strconv.Atoi(firstLine[1])

		if err != nil {
			log.Fatalf("couldn't convert number: %v\n", err)
		}

	}

	for scanner.Scan() {

		//remove spaces
		thisLine := strings.Replace(scanner.Text(), " ", "", -1)

		//convert to uint
		l, err := strconv.ParseUint(thisLine, 2, rowSize)

		if err != nil {
			log.Fatal(err)
		}

		//create Vertex
		v := &Vertex{l, nil, math.Inf(1), true, nil, []*Vertex{}, -1}
		v.shortestEdge = v
		v.leader = v
		v.followers = append(v.followers, v)

		//add Vertex to Vertices
		_, ok := Vertices[l]

		if !ok {
			Vertices[l] = v
			LeaderMap[l] = v
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func cluster() {
	for _, x := range Vertices {
		for _, h := range hamInts {

			w, ok := Vertices[x.val^h]

			if ok {
				if x.leader != w.leader {
					if len(x.leader.followers) > len(w.leader.followers) {
						x.leader.ConsumeFollowers(w.leader)
					} else {
						w.leader.ConsumeFollowers(x.leader)
					}
				}
			}
		}
	}
}

func TwoDegreeHamInts(s int) {

	for i := 0; i < s; i++ {
		hamInts = append(hamInts, uint64(math.Pow(2, float64(i))))

		for j := s - 1; j >= 0; j-- {
			hamInts = append(hamInts, uint64(math.Pow(2, float64(i))+math.Pow(2, float64(j))))
		}
	}
}
