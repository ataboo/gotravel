package main

import (
	"fmt"
	gen "github.com/ataboo/gotravel/genetics"
	"github.com/ataboo/gotravel/rendering"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

//var cities []gen.City

//var testMap gen.RoadMap

func main() {
	//cities = gen.MakeCircleCities(24)
	//testMap = gen.RandomRoadmap(cities)

	gtk.Init(nil)
	initWindow()
	gtk.Main()
}

func initWindow() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	win.SetSizeRequest(800, 600)


	win.SetTitle("Go Genetic Traveller")
	win.Connect("destroy", gtk.MainQuit)

	da, _ := gtk.DrawingAreaNew()
	da.SetVExpand(true)
	da.SetHExpand(true)
	rendering.Start(da)

	button, _ := gtk.ButtonNewFromIconName("view-refresh", gtk.ICON_SIZE_BUTTON)
	horiz, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 40)
	horiz.SetHAlign(gtk.ALIGN_END)
	vert, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 40)
	vert.Add(horiz)
	horiz.Add(button)
	vert.Add(da)
	win.Add(vert)

	button.Connect("clicked", clickRefresh)

	da.Connect("draw", rendering.Draw)
	win.Connect("check-resize", rendering.OnResize)

	win.ShowAll()

}

func clickRefresh() {
	fmt.Println("Clicked Refresh!")
	//bestMap.Shuffle()
	//da.QueueDraw()
	go runGenes()
}

func runGenes() {
	cfg := gen.GeneCfg{
		CityCount: 80,
		PopCap: 200,
		MaxGenerations: 10000000,
		CullRate: 0.8,
		MutateRate: 0.10,
		MutateDeviation: 0.5,
		StatPeriod: 100,
		//Delay: time.Millisecond * 100.0,
	}

	stats, _ := gen.RunGenetic(cfg)
	var stat gen.GeneStats
	chanFor:
	for {
		select {
		case stat = <-stats:
			if stat.Generation < 0 {
				break chanFor
			}

			fmt.Printf("%d  |  %.2f\n", stat.Generation, stat.BestMap.Cost())
			rendering.ShowRoadmap(&stat.BestMap)
			//if stat.BestMap.Solved() {
			//	fmt.Printf("Solved in %d!\n", stat.Generation)
			//	stop <- 0
			//	break chanFor
			//}
		}
	}
}

