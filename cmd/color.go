package main

type SColor uint8

const (
    WHITE = 0
    LGRAY = 1
    GRAY = 2
    BLACK = 3
    KIRBY = 4
    RED = 5
    ORANGE = 6
    BROWN = 7
    YELLOW = 8
    lGREEN = 9
    GREEN = 10
    LBLUE = 11
    TURQUOISE = 12
    BLUE = 13
    PINK = 14
    VIOLET = 15
)

func isValidColor(color SColor) bool {
    return color >= WHITE && color <= VIOLET
}
