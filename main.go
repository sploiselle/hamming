package main

import (
	"bufio"
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
	val    uint64
	leader *Vertex
	rank   int
}

func (v Vertex) String() string {

	return fmt.Sprintf("\nval:\t%vLeader:\t%vRank:\t%v",
		v.val, v.leader.val, v.rank)
}

var Vertices = make(map[uint64]*Vertex, numOfRows)
var LeaderMap = make(map[uint64]*Vertex, numOfRows)

func main() {

	readFile(os.Args[1])

	TwoDegreeHamInts(rowSize)

	cluster()

	fmt.Println(len(LeaderMap))

}

func readFile(filename string) {
	file, err := os.Open(filename)

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

		thisLine := strings.Replace(scanner.Text(), " ", "", -1)
		l, err := strconv.ParseUint(thisLine, 2, rowSize)

		if err != nil {
			log.Fatal(err)
		}

		v := &Vertex{l, nil, 0}
		v.leader = v

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

// Cluster Vertices based on max hamming distance
// defined by inverse hamming values.
func cluster() {
	for _, x := range Vertices {
		for _, h := range hamInts {
			w, ok := Vertices[x.val^h]
			if ok {
				xl := FindAndUpdate(x)
				wl := FindAndUpdate(w)

				if xl != wl {
					Union(xl, wl)
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

// FindAndUpdate performs a Union-Find find operations
// while implementing Union by Rank path compression
func FindAndUpdate(x *Vertex) *Vertex {
	if x.leader == x {
		return x
	}
	x.leader = FindAndUpdate(x.leader)
	return x.leader
}

// Union performs Union-Find union operations
// implementing Union by Rank
func Union(x, y *Vertex) {
	switch {
	case x.rank < y.rank:
		x.leader = y
		delete(LeaderMap, x.val)
	case x.rank == y.rank:
		x.rank++
		fallthrough
	case x.rank > y.rank:
		y.leader = x
		delete(LeaderMap, y.val)
	}
}
