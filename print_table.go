package main

import (
	"fmt"
	"log"
)

func print_table(
	t *[][]byte,
	bt *[][]int,
	table_height,
	table_width int) {

	s := ""
	for i := range table_height {
		for j := range table_width {
			s += string((*t)[i][j])
			if (*t)[i][j] == '*' {
				s += fmt.Sprintf("%d ", (*bt)[i][j])
			} else {
				s += "  "
			}
		}
		s += "\n"
	}
	log.Printf("\nloaded table:\n%s", s)
}
