package main

import "fmt"

type Margin struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

// m과 m2가 같은지 비교한다.
func (m *Margin) Same(m2 *Margin) bool {
	return m.Left == m2.Left && m.Right == m2.Right && m.Top == m2.Top && m.Bottom == m2.Bottom
}

// GetFloorType 바닥에 닫는 셀별 인덱스를 문자열로 반환한다.
// ex, ㅜ 모양이면 212: 바닥에 닫는 셀이 첫번째 셀은 2칸 높이, 두번째 셀은 1칸 높이, 세번째 셀은 2칸 높이
// ㅗ 모양이면 111: 바닥에 닫는 셀이 모두 1칸 높이
// ⌈ 모양이면 13: 바닥에 닫는 셀이 첫번째 셀은 1칸 높이, 두번째 셀은 3칸 높이
func GetFloorType(cells [][]int) string {

	top, bottom := -1, -1
	for y, row := range cells {
		for _, cell := range row {
			if cell == 1 {
				if top == -1 || y < top {
					top = y
				}
				if bottom == -1 || y > bottom {
					bottom = y
				}
			}
		}
	}

	cells = cells[top : bottom+1]

	ret := ""

	depth := len(cells)
	for x := 0; x < len(cells[0]); x++ {
		idx := -1
		for y := depth - 1; y >= 0; y-- {
			if cells[y][x] == 1 {
				idx = y + 1
				break
			}
		}

		if idx == -1 {
			ret += "0"
		} else {
			ret += fmt.Sprintf("%d", (depth-idx)+1)
		}
	}

	return ret
}

// cells이 결합 가능여부를 반환한다.
func CanCombine(cells1, cells2 [][]int, offsetX, offsetY int) bool {

	for y := 0; y < len(cells1); y++ {
		for x := 0; x < len(cells1[0]); x++ {
			if (cells1[y][x] != 0 && cells2[y+offsetY][x+offsetX] != 0) || (cells1[y][x] == 0 && cells2[y+offsetY][x+offsetX] == 0) {
				Trace.Printf("CanCombine false y:%d x:%d y+offsetY:%d x+offsetX:%d", y, x, y+offsetY, x+offsetX)
				return false
			}
		}
	}

	Trace.Printf("CanCombine true")

	return true
}

// cells1과 cells2가 같은지 비교한다.
func SameCells(cells1, cells2 [][]int) bool {

	for y := 0; y < len(cells1); y++ {
		for x := 0; x < len(cells1[0]); x++ {
			if cells1[y][x] != cells2[y][x] {
				return false
			}
		}
	}

	return true
}

// left, top, right, bottom의 비어있는 셀을 제거한다.
func TrimShape(cells [][]int) ([][]int, Margin) {

	left, top, right, bottom := -1, -1, -1, -1
	for y, row := range cells {
		for x, cell := range row {
			if cell == 1 {
				if left == -1 || x < left {
					left = x
				}
				if right == -1 || x > right {
					right = x
				}
				if top == -1 || y < top {
					top = y
				}
				if bottom == -1 || y > bottom {
					bottom = y
				}
			}
		}
	}

	trimed := make([][]int, bottom-top+1)
	copy(trimed, cells[top:bottom+1])
	for i := top; i <= bottom; i++ {
		trimed[i-top] = make([]int, right-left+1)
		copy(trimed[i-top], cells[i][left:right+1])
	}

	margin := Margin{
		Left:   left,
		Right:  len(cells[0]) - right - 1,
		Top:    top,
		Bottom: len(cells) - bottom - 1,
	}

	if right == -1 {
		margin.Right = 0
	}
	if bottom == -1 {
		margin.Bottom = 0
	}

	return trimed, margin
}

// cells의 가장 낮은 셀의 x좌표를 반환한다.
func findLowest(cells [][]int, v int) int {

	col := 0
	row := 0

	for y := len(cells) - 1; y >= 0; y-- {
		for x := 0; x < len(cells[0]); x++ {
			if cells[y][x] == v && row < y {
				row = y
				col = x
			}
		}
	}

	return col
}

// cell을 위로 이동했을때 막히는 cell을 채워준다.
func fillTailToUp(cells [][]int) [][]int {

	// copy
	tailed := make([][]int, len(cells))
	for i := range tailed {
		tailed[i] = make([]int, len(cells[0]))
		copy(tailed[i], cells[i])
	}

	for x := 0; x < len(tailed[0]); x++ {
		sy := len(tailed)
		for y := 0; y < len(tailed); y++ {
			if tailed[y][x] == 1 {
				sy = y
				break
			}
		}

		for y := sy; y < len(tailed); y++ {
			tailed[y][x] = 1
		}
	}

	return tailed
}

// cell을 아래 이동했을때 막히는 cell을 채워준다.
func fillTailToDown(cells [][]int) [][]int {

	// copy
	tailed := make([][]int, len(cells))
	for i := range tailed {
		tailed[i] = make([]int, len(cells[0]))
		copy(tailed[i], cells[i])
	}

	for x := 0; x < len(tailed[0]); x++ {
		sy := 0
		for y := len(tailed) - 1; y >= 0; y-- {
			if tailed[y][x] == 1 {
				sy = y
				break
			}
		}

		for y := sy - 1; y >= 0; y-- {
			tailed[y][x] = 1
		}
	}

	return tailed
}

func cellsToString(cells [][]int) string {

	ret := ""
	for _, row := range cells {
		ret += fmt.Sprintf("%v\n", row)
	}

	return ret
}
