package core

type SColor uint8

const (
	WHITE     SColor = 0
	LGRAY     SColor = 1
	GRAY      SColor = 2
	BLACK     SColor = 3
	KIRBY     SColor = 4
	RED       SColor = 5
	ORANGE    SColor = 6
	BROWN     SColor = 7
	YELLOW    SColor = 8
	lGREEN    SColor = 9
	GREEN     SColor = 10
	LBLUE     SColor = 11
	TURQUOISE SColor = 12
	BLUE      SColor = 13
	PINK      SColor = 14
	VIOLET    SColor = 15
)

func isValidColor(color SColor) bool {
	return color >= WHITE && color <= VIOLET
}
