package atamath

import (
	"math"
	"math/rand"
)

type Vect2 struct {
	X float64
	Y float64
}

func (v Vect2) Scale(f float64) Vect2 {
	return Vect2{v.X * f, v.Y * f}
}

func (a Vect2) Add(b Vect2) Vect2 {
	return Vect2{a.X + b.X, a.Y + b.Y}
}

func (a Vect2) Sub(b Vect2) Vect2 {
	return Vect2{a.X - b.X, a.Y - b.Y}
}

func (a Vect2) AddInt(b IntVect2) Vect2 {
	return Vect2{a.X + float64(b.X), a.Y + float64(b.Y)}
}

func (a Vect2) MagSqr() float64 {
	return a.X * a.X + a.Y * a.Y
}

func (a Vect2) Mag() float64 {
	return math.Sqrt(a.MagSqr())
}

func ZeroVect() Vect2 {
	return Vect2{0, 0}
}

func RandVectInRange(s Vect2, d Vect2) Vect2 {
	x := rand.Float64() * d.X + s.X
	y := rand.Float64() * d.Y + s.Y

	return Vect2{x, y}
}

type IntVect2 struct {
	X int
	Y int
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}