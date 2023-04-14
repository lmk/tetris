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

func (b *Block) ShapeFrom(shape [][]int) {
	b.Shape = make([][]int, BLOCK_ROW)
	for i := 0; i < BLOCK_ROW; i++ {
		b.Shape[i] = make([]int, BLOCK_COLUMN)
		copy(b.Shape[i], shape[i])
	}
}

func (b *Block) ShapeTo() [][]int {
	shape := make([][]int, BLOCK_ROW)
	for i := 0; i < BLOCK_ROW; i++ {
		shape[i] = make([]int, BLOCK_COLUMN)
		copy(shape[i], b.Shape[i])
	}

	return shape
}

func (b *Block) Clone(from Block) {
	b.Row = from.Row
	b.Col = from.Col
	b.ShapeIndex = from.ShapeIndex
	b.ShapeFrom(from.Shape)
}

func NewBlock(shapeIndex int) Block {
	block := Block{
		Row:        0,
		Col:        BOARD_CENTER,
		ShapeIndex: shapeIndex,
	}
	block.SetShape(shapeIndex)

	return block
}

func (b *Block) SetShape(shapeIndex int) {
	b.ShapeIndex = shapeIndex
	b.ShapeFrom(SHAPES[shapeIndex-1])
}

func (b *Block) Rotate() {
	newShape := make([][]int, BLOCK_ROW)
	for i := 0; i < BLOCK_ROW; i++ {
		newShape[i] = make([]int, BLOCK_COLUMN)
		for j := 0; j < BLOCK_COLUMN; j++ {
			newShape[i][j] = b.Shape[BLOCK_COLUMN-j-1][i]
		}
	}
	b.ShapeFrom(newShape)
}
