package genetics

import (
	"fmt"
	"sort"
	"testing"
)

func TestMakePopulation(t *testing.T) {
	cfg := GeneCfg{
		CityCount: 4,
		PopCap: 4,
	}

	pop := makePopulation(cfg)

	if pop.cityCount() != 4 {
		t.Error("Expected different city count", pop.cityCount())
	}

	if len(pop) != 4 {
		t.Error("Unexpected population length", len(pop))
	}
}

func TestRank(t *testing.T) {
	cities := MakeCircleCities(4)
	map1 := RoadMap{
		Order: []int{0, 1, 2, 3},
		cities: cities,
	}

	map2 := RoadMap{
		Order: []int{1, 3, 0, 2},
		cities: cities,
	}

	pop := Population{
		map2,
		map1,
	}

	if map1.Cost() >= map2.Cost() {
		t.Error("Expected map 1 to have lower cost.", map1, map2)
	}

	pop = pop.rank()

	if pop[0].Cost() != map1.Cost() || pop[1].Cost() != map2.Cost() {
		t.Error("Unexpected ranking.", pop)
	}
}

func TestCull(t *testing.T) {
	cfg := GeneCfg{
		PopCap:10,
		CityCount:16,
		CullRate:0.20,
	}

	pop := makePopulation(cfg)

	pop = pop.cull(cfg)

	if len(pop) != 8 {
		t.Error("Unexpected population length", len(pop))
	}
}

func TestNormalize(t *testing.T) {
	table := [][]int {
		{0, 1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1, 0},
		{1, 2, 3, 2, 5, 0},
		{1, 1, 1, 1, 1, 1},
		{2, 3, 1, 4, 5, 1},
	}

	for _, row := range table {
		cp := make([]int, len(row))
		for i, v := range row {
			cp[i] = v
		}

		normalize(row)

		if len(cp) != len(row) {
			t.Error("Unexpected length", cp, row)
		}

		sort.Ints(row)
		for i, v := range row {
			if i != v {
				t.Error("Value missing (result has been sorted in testing)", cp, row)
			}
		}
	}

	checkOrder := []int {2, 1, 3, 3}

	normalize(checkOrder)

	if checkOrder[0] != 2 || checkOrder[1] != 1 {
		t.Error("Unexpected order", checkOrder)
	}
}

func TestBreed(t *testing.T) {
	cfg := GeneCfg{
		PopCap: 2,
		CityCount: 8,
		CullRate:0.5,
	}

	p := makePopulation(cfg)

	cfg.PopCap = 4

	old1:=make([]int, 8)
	copy(old1, p[0].Order)
	old2:=make([]int, 8)
	copy(old2, p[1].Order)

	p = p.breedBySplicing(cfg)

	if len(p) != 4 {
		t.Error("Unexpected population count", p)
	}

	for i:=0; i<8; i++ {
		if p[0].Order[i] != old1[i] {
			t.Error("Unexpected order value", p[0].Order, old1)
		}
		if p[1].Order[i] != old2[i] {
			t.Error("Unexpected order value", p[0].Order, old1)
		}
	}

	if len(p[2].Order) != 8 {
		t.Error("Unexpected order length in offspring", p[2])
	}
}

func TestMutate(t *testing.T) {
	cfg := GeneCfg{
		PopCap: 1,
		CityCount: 8,
		MutateRate: 0.25,
		MutateDeviation: 0.1,
	}

	p := makePopulation(cfg)
	before := make([]int, 8)
	copy(before, p[0].Order)
	cfg.PopCap = 4

	p = p.mutate(cfg)

	for i:=0; i<len(before); i++ {
		if p[0].Order[i] != before[i] {
			t.Error("Unexpected mismatch", p[0].Order, before)
		}
	}

	t.Log(p)

	if len(p) != cfg.PopCap {
		t.Error("Unexpected population", p, cfg.PopCap)
	}
}

func TestBlah(t *testing.T) {
	cfg := GeneCfg{
		CityCount: 200,
		PopCap: 1000,
		MaxGenerations: 20000,
		CullRate: 0.95,
		MutateRate: 0.75,
		MutateDeviation: 0.25,
	}

	stats, stop := RunGenetic(cfg)
	var stat GeneStats
	chanFor:
	for {
		select {
		case stat = <-stats:
			if stat.Generation < 0 {
				break chanFor
			}

			fmt.Printf("%d  |  %.2f of %.2f\n", stat.Generation, stat.BestMap.Cost(), stat.BestMap.BestCost())
			if stat.BestMap.Solved() {
				fmt.Printf("Solved in %d!\n", stat.Generation)
				stop <- 0
				break chanFor
			}
		}
	}

	//fmt.Printf("Best solution: %.2f, Perfect: %.2f\n", stat.BestMap.Cost(), stat.BestMap.BestCost())


}