package genetics

import (
	"github.com/ataboo/gotravel/atamath"
	"log"
	"math"
	"math/rand"
	"time"
)

var randGen = rand.New(rand.NewSource(time.Now().Unix()))


type RoadMap struct {
	Order []int
	magnitude float64
	bestCost float64
	dirty bool
	cities []City
}

type City struct {
	Pos atamath.Vect2
	Order int
}


func (rm *RoadMap) Cost() float64 {
	if len(rm.Order) != len(rm.cities) {
		log.Fatal("Mismatch between RoadMap and Cities")
	}

	if rm.dirty || rm.magnitude == 0.0 {
		cost := 0.0

		ForEachCity(rm.OrderedCities(), func(a *City, b *City) {
			cost += b.Pos.Sub(a.Pos).Mag()
		})

		rm.magnitude = cost
		rm.dirty = false
	}

	return rm.magnitude
}

func (rm *RoadMap) Solved() bool {
	return math.Abs(rm.BestCost() - rm.Cost()) < 1e-3
}

func (rm RoadMap) Clone() RoadMap {
	newMap := RoadMap{
		cities: rm.cities,
		dirty:  rm.dirty,
		magnitude: rm.magnitude,
		Order: make([]int, len(rm.Order)),
	}

	copy(newMap.Order, rm.Order)

	return newMap
}

func (rm RoadMap) BestCost() float64 {
	if rm.bestCost == 0 {
		rm.bestCost = rm.cities[1].Pos.Sub(rm.cities[0].Pos).Mag() * float64(len(rm.cities))
	}
	return rm.bestCost
}

func (rm *RoadMap) Shuffle() *RoadMap {
	rm.dirty = true
	rm.Order = randGen.Perm(len(rm.Order))

	return rm
}

func (rm *RoadMap) OrderedCities() []*City {
	ordered := make([]*City, len(rm.Order))

	for i, o := range rm.Order {
		ordered[i] = &rm.cities[o]
	}

	return ordered
}

func MakeCircleCities(count int) []City {
	cities := make([]City, count)

	angleStep := 2.0 * math.Pi / float64(count)
	angle := 0.0
	rad := 1.0
	for i:=0; i<count; i++ {
		pos := atamath.Vect2{X: rad * math.Cos(angle), Y: rad * math.Sin(angle)}

		cities[i] = City{Pos: pos, Order: i}
		angle += angleStep
	}

	return cities
}

func MakeRandoCities(count int) []City {
	cities := make([]City, count)

	for i:=0; i<count; i++ {
		cities[i] = City{
			Pos: atamath.RandVectInRange(atamath.Vect2{-1, -1}, atamath.Vect2{2, 2}),
			Order: i,
		}
	}

	return cities
}

func RandomRoadmap(cities []City) RoadMap {
	perm := randGen.Perm(len(cities))
	rm := RoadMap{dirty:true, Order:perm, cities:cities}

	return rm
}

type Neighbours func(a *City, b *City)

func ForEachCity(c []*City, callback Neighbours) {
	for i, a := range c {
		var b *City
		if i == len(c) - 1 {
			b = c[0]
		} else {
			b = c[i + 1]
		}

		callback(a, b)
	}
}