package rendering

import (
	"fmt"
	"github.com/ataboo/gotravel/atamath"
	gen "github.com/ataboo/gotravel/genetics"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"image/color"
	"math"
	"time"
)

var canvasSize atamath.Vect2
var canvasCenter atamath.Vect2
var canvasScale float64
var resizingWindow = false
var resizeStart = make(chan int)
var da *gtk.DrawingArea

var visibleMap *gen.RoadMap


func Start(d *gtk.DrawingArea) {
	da = d

	go resizeWatcher()
}

func ShowRoadmap(rm *gen.RoadMap) {
	visibleMap = rm
	da.QueueDraw()
}

func Draw(da *gtk.DrawingArea, cr *cairo.Context) {
	cr.SetSourceRGB(0.9, 0.9, 0.9)
	cr.Rectangle(0, 0, float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight()))
	cr.Fill()

	if !resizingWindow {
		if visibleMap != nil {
			drawRoadmap(visibleMap, da, cr)
		}
	}
}

// Should be called when the window is re-sized.
// Triggers the resizeWatcher to prevent rendering for a time.
func OnResize(_ *gtk.Window) bool {
	resizeStart <- 0

	return false
}

// When the resizeStart channel is called, this toggles the resizingWindow flag 'true', delays, then 'false' again.
// Triggering the resizeStart channel while it is already 'true' will reset the timer.
func resizeWatcher() {
	var expiry <- chan time.Time

	for {
		select {
		case <-resizeStart:
			resizingWindow = true
			expiry = time.After(time.Millisecond * 100.0)
		case <- expiry:
			resizingWindow = false
			doneResize()
		}
	}
}

// Trigger draw events after the window has been re-sized and the timer has elapsed
func doneResize() {
	updateCanvas()
	da.QueueDraw()
}

func updateCanvas() {
	canvasSize = atamath.Vect2{float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight())}

	canvasScale = math.Min(canvasSize.X, canvasSize.Y)
	canvasCenter = atamath.Vect2{canvasSize.X / 2, canvasSize.Y / 2}.Scale(1/canvasScale)
}

func drawRoadmap(rm *gen.RoadMap, da *gtk.DrawingArea, cr *cairo.Context) {

	cr.Scale(canvasScale, canvasScale)

	cr.SetSourceRGB(0.9, 0.9, 0.9)
	cr.Rectangle(0, 0, 2, 2)
	cr.Fill()


	gen.ForEachCity(rm.OrderedCities(), func(a *gen.City, b *gen.City) {
		drawCity(a, cr, color.RGBA{0, 0, 200, 255})
	})

	gen.ForEachCity(rm.OrderedCities(), func(a *gen.City, b *gen.City) {
		drawRoad(a, b, cr, color.RGBA{200, 0, 0, 255})
	})

	cr.Stroke()


	cr.Scale(1/canvasScale, 1/canvasScale)
	costText := fmt.Sprintf("Cost: %.2f", rm.Cost())
	cr.SelectFontFace("sans-serif", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(16)

	extents := cr.TextExtents(costText)

	cr.SetSourceRGB(0.0, 0.0, 0.0)
	cr.MoveTo(canvasSize.X / 2.0 - extents.Width / 2.0, canvasSize.Y - 20.0)
	//cr.MoveTo(canvasCenter.X - extents.Width/2.0, canvasCenter.Y)
	cr.ShowText(costText)
}

func drawCity(c *gen.City, cr *cairo.Context, col color.RGBA) {
	circPos := c.Pos.Scale(0.4).Add(canvasCenter)

	cr.SetSourceRGBA(SplitRGBA(col))
	cr.Arc(circPos.X, circPos.Y, 0.01, 0, 360)
	cr.Fill()
}

func drawRoad(start *gen.City, end *gen.City, cr *cairo.Context, col color.RGBA) {
	startPos := start.Pos.Scale(0.4).Add(canvasCenter)
	destoPos := end.Pos.Scale(0.4).Add(canvasCenter)


	cr.SetLineWidth(0.005)
	cr.SetSourceRGBA(SplitRGBA(col))
	cr.MoveTo(startPos.X, startPos.Y)
	cr.LineTo(destoPos.X, destoPos.Y)
}

func SplitRGBA(col color.RGBA) (float64, float64, float64, float64) {
	return float64(col.R)/255.0, float64(col.G)/255.0, float64(col.B)/255.0, float64(col.A)/255.0
}