package main

var (
	SHAPES = [...][][]int{
		{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 0, 0},
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
	}
)

type Block struct {
	Row        int     `json:"row"`
	Col        int     `json:"col"`
	Shape      [][]int `json:"shape"`
	ShapeIndex int     `json:"shapeIndex"`
}

func CloneShape(shape [][]int) [][]int {
	newShape := make([][]int, len(shape))
	for i, c := range shape {
		newShape[i] = make([]int, len(c))
		copy(newShape[i], c)
	}

	return newShape
}

func (b *Block) Clone(from *Block) {
	b.Row = from.Row
	b.Col = from.Col
	b.ShapeIndex = from.ShapeIndex
	b.Shape = CloneShape(from.Shape)
}

// NewBlock
// shapeIndex : 1 ~ 7
func NewBlock(shapeIndex int) *Block {
	block := Block{
		Row:        0,
		Col:        BOARD_CENTER,
		ShapeIndex: shapeIndex,
	}
	block.SetShape(shapeIndex)

	return &block
}

func (b *Block) SetShape(shapeIndex int) {
	b.ShapeIndex = shapeIndex
	b.Shape = CloneShape(SHAPES[shapeIndex-1])
}

func (b *Block) Rotate() {
	newShape := make([][]int, BLOCK_ROW)
	for i := 0; i < BLOCK_ROW; i++ {
		newShape[i] = make([]int, BLOCK_COLUMN)
		for j := 0; j < BLOCK_COLUMN; j++ {
			newShape[i][j] = b.Shape[BLOCK_COLUMN-j-1][i]
		}
	}
	b.Shape = CloneShape(newShape)
}

func (b *Block) inBlock(row, col int) bool {
	return b.Row <= row && row < b.Row+BLOCK_ROW && b.Col <= col && col < b.Col+BLOCK_COLUMN
}
