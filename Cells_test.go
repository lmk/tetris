package main

import (
	"fmt"
	"testing"
)

func TestGetFloorType(t *testing.T) {

	for i := range SHAPES {
		block := NewBlock(i + 1)
		floorType := GetFloorType(block.Shape)

		switch i {
		case 0:
			if floorType != "1111" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 0")
			}

		case 1:
			if floorType != "2120" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 1")
			}

		case 2:
			if floorType != "0110" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 2")
			}

		case 3:
			if floorType != "0210" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 3")
			}

		case 4:
			if floorType != "0110" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 4")
			}

		case 5:
			if floorType != "0110" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 5")
			}

		case 6:
			if floorType != "0120" {
				fmt.Println(block.Shape)
				fmt.Println(floorType)
				t.Error("TestTrimShape 6")
			}
		}
	}

	cells := [][]int{
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	}

	floorType := GetFloorType(cells)

	if floorType != "0310" {
		fmt.Println(cells)
		fmt.Println(floorType)
		t.Error("TestTrimShape 7")
	}
}

func TestTrimShape(t *testing.T) {
	for i := range SHAPES {
		block := NewBlock(i + 1)
		trimedCells, margin := TrimShape(block.Shape)

		switch i {
		case 0:
			if SameCells(trimedCells, [][]int{{1, 1, 1, 1}}) == false || !margin.Same(&Margin{0, 1, 0, 2}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 1:
			if SameCells(trimedCells, [][]int{
				{1, 1, 1},
				{0, 1, 0},
			}) == false || !margin.Same(&Margin{0, 1, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 2:
			if SameCells(trimedCells, [][]int{
				{1, 1},
				{1, 1},
			}) == false || !margin.Same(&Margin{1, 1, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 3:
			if SameCells(trimedCells, [][]int{
				{1, 0},
				{1, 1},
				{0, 1},
			}) == false || !margin.Same(&Margin{1, 0, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 4:
			if SameCells(trimedCells, [][]int{
				{1, 0},
				{1, 0},
				{1, 1},
			}) == false || !margin.Same(&Margin{1, 0, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 5:
			if SameCells(trimedCells, [][]int{
				{0, 1},
				{0, 1},
				{1, 1},
			}) == false || !margin.Same(&Margin{1, 0, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}

		case 6:
			if SameCells(trimedCells, [][]int{
				{0, 1},
				{1, 1},
				{1, 0},
			}) == false || !margin.Same(&Margin{1, 0, 1, 1}) {
				t.Log(block.Shape)
				t.Log(trimedCells)
				t.Log(margin)
				t.Error("TestTrimShape", i)
			} else {
				t.Log("pass", i)
			}
		}
	}
}
