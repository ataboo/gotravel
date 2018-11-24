package genetics

import (
	"github.com/ataboo/gotravel/atamath"
	"math"
	"sort"
	"time"
)


type GeneCfg struct {
	CityCount int
	PopCap int
	MaxGenerations int
	CullRate float64
	CullReprieve float64
	MutateRate float64
	MutateDeviation float64
	Delay time.Duration
	StatPeriod int
	RandomCityPos bool
}

type GeneStats struct {
	Generation int
	BestMap RoadMap
}

type Population []RoadMap

func RunGenetic(cfg GeneCfg) (statChan chan GeneStats, stopChan chan int) {
	statChan = make(chan GeneStats)
	stopChan = make(chan int)

	go func() {
		pop := makePopulation(cfg)
		pop = pop.rank()

		for i := 0; i < cfg.MaxGenerations; i++ {
			pop = pop.cull(cfg)
			pop = pop.shuffle().breedBySplicing(cfg).mutate(cfg).rank()

			stats := GeneStats{
				i,
				pop[0].Clone(),
			}

			if cfg.StatPeriod == 0 || i%cfg.StatPeriod == 0 {
				select {
				case statChan <- stats:
				case <-stopChan:
					return
					//default:
					//	fmt.Print("Skipped msg send\n")
				}
			}

			time.Sleep(cfg.Delay)
		}

		statChan<-GeneStats{Generation: -1}
	}()

	return statChan, stopChan
}

func makePopulation(cfg GeneCfg) Population {
	pop := make(Population, cfg.PopCap)
	var cities []City
	if cfg.RandomCityPos {
		cities = MakeRandoCities(cfg.CityCount)
	} else {
		cities = MakeCircleCities(cfg.CityCount)
	}
	for i:=0; i<len(pop); i++ {
		pop[i] = RandomRoadmap(cities)
	}

	return pop
}

func (p Population) cityCount() int {
	return len(p[0].Order)
}

func (p Population) rank() Population {
	sort.Slice(p, func(i, j int) bool {
		return p[i].Cost() < p[j].Cost()
	})

	return p
}

func (p Population)cull(cfg GeneCfg) Population {
	cutoffIndex := int((1 - cfg.CullRate) * float64(cfg.PopCap))

	keeper := p[0:cutoffIndex]
	remainder := p[cutoffIndex:]

	reprieveCount := int(cfg.CullReprieve * float64(cfg.PopCap))
	reprievePerm := randGen.Perm(len(remainder))
	reprieve := make([]RoadMap, reprieveCount)
	for i:=0; i<reprieveCount; i++ {
		reprieve[i] = remainder[reprievePerm[i]]
	}

	keeper = append(keeper, reprieve...)

	return keeper
}

func (p Population)shuffle() Population {
	randGen.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})

	return p
}

func (p Population)breedBySplicing(cfg GeneCfg) Population {
	children := make([]RoadMap, 0)

	for i := 0; i<len(p)-1; i+=2 {
		a := p[i]
		b := p[i+1]

		ab, ba := splice(a.Order, b.Order)

		abMap := RoadMap{
			Order:normalize(ab),
			cities:a.cities,
			dirty:true,
		}
		baMap := RoadMap{
			Order:normalize(ba),
			cities:a.cities,
			dirty:true,
		}

		children = append(children, abMap, baMap)
	}

	return append(p, children...)
}

func splice(a []int, b []int) ([]int, []int){
	ab := make([]int, len(a))
	ba := make([]int, len(a))

	split := randGen.Intn(len(a))

	copy(ab, a)
	copy(ba, b)
	ab = append(ab[0:split], b[split:]...)
	ba = append(ba[0:split], a[split:]...)

	return ab, ba
}

// Remove duplicates randomly to make a a sequence from 0 to len(a) while keeping order
func normalize(a []int) []int {
	instances := make([][]int, len(a))
	missing := make([]int, 0)
	duplicated := make([][]int, 0)

	for i, n := range a {
		instances[n] = append(instances[n], i)
	}

	for i, inst := range instances {
		if len(inst) == 0 {
			missing = append(missing, i)
		} else if len(inst) > 1 {
			duplicated = append(duplicated, inst)
		}
	}

	for _, m := range missing {
		iDup := 0
		if len(duplicated) > 1 {
			iDup = randGen.Intn(len(duplicated) - 1)
		}
		dupRow := duplicated[iDup]

		jDup := randGen.Intn(len(dupRow) - 1)
		replacedIndex :=dupRow[jDup]

		a[replacedIndex] = m

		if len(dupRow) == 2 {
			if iDup == len(duplicated) - 1 {
				duplicated = duplicated[0:iDup]
			} else {
				duplicated = append(duplicated[0:iDup], duplicated[iDup+1:]...)
			}
		} else {
			if jDup == len(dupRow) - 1 {
				duplicated[iDup] = dupRow[0:jDup]
			} else {
				duplicated[iDup] = append(dupRow[0:jDup], dupRow[jDup+1:]...)
			}
		}
	}

	return a
}

func (p Population)mutate(cfg GeneCfg) Population {
	newMaps := make(Population, cfg.PopCap - len(p))

	for i:=0; i<cfg.PopCap - len(p); i++ {
		newMap := RandomRoadmap(p[0].cities)
		copy(newMap.Order, p[randGen.Intn(len(p))].Order)
		newMaps[i] = newMap

		rate := (randGen.Float64() * 2 - 1) * cfg.MutateDeviation + cfg.MutateRate
		swapCount := int(math.Round(rate * float64(cfg.CityCount / 2)))
		swapCount = atamath.MinInt(cfg.CityCount / 2, swapCount)

		if swapCount < 1 {
			continue
		}

		perms := randGen.Perm(cfg.CityCount)
		for i:=0; i<swapCount*2; i+=2 {
			newMap.Order[perms[i]], newMap.Order[perms[i+1]] = newMap.Order[perms[i+1]], newMap.Order[perms[i]]
		}
	}

	p = append(p[0:], newMaps...)

	return p
}