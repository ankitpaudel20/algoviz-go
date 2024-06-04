package main

import (
	"image/color"
)

const (
	outer_padding = 10
)

type CellState int

const (
	CellStateNotVisited CellState = iota
	CellStateVisited
	CellStateQueued
)

var (
	CellStateVisitedColor color.RGBA = color.RGBA{160, 32, 240, 255}
)

type GridCell struct {
	col color.RGBA
}
type GridCanvasState struct {
	prev_State *GridCanvasState
	next_State *GridCanvasState
	changes    []GridCell
}
type GridCanvas struct {
	sizeX     int
	sizeY     int
	cellSizeX int
	cellSizeY int
	Cells     []GridCell
}
