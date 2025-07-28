package main

import (
	"fmt"

	"github.com/kahnaisehC/sajcredez/pkg"
)

func main() {
	s := board.CreateSajcredez()
	fmt.Println(s.GetBoardString())
}
