package main

import "log"

func find_coords_touching_revealed_areas(
	t *[][]byte,
	bt *[][]int,
	table_height,
	table_width,
	start_y,
	start_x int,
) (
	found_y,
	found_x,
	neighbor_element_borders_reached int,
) {

	found_y, found_x, neighbor_element_borders_reached = -1, -1, -1

	max_diag := table_height
	if max_diag < table_width {
		max_diag = table_width
	}

	var found bool
	for step := range max_diag {

		for x := start_x - step; x <= start_x+step; x++ {

			y := start_y - step
			found,
				neighbor_element_borders_reached =
				validate_touch_position(t, bt, table_height, table_width, y, x)

			if found {
				return y, x, neighbor_element_borders_reached + 1
			}

			y = start_y + step
			found,
				neighbor_element_borders_reached =
				validate_touch_position(t, bt, table_height, table_width, y, x)

			if found {
				return y, x, neighbor_element_borders_reached + 1
			}
		}

		for y := start_y - step + 1; y <= start_y+step-1; y++ {

			x := start_x - step
			found,
				neighbor_element_borders_reached =
				validate_touch_position(t, bt, table_height, table_width, y, x)

			if found {
				return y, x, neighbor_element_borders_reached + 1
			}

			x = start_x + step
			found,
				neighbor_element_borders_reached =
				validate_touch_position(t, bt, table_height, table_width, y, x)

			if found {
				return y, x, neighbor_element_borders_reached + 1
			}
		}

	}

	return found_y, found_x, neighbor_element_borders_reached
}

func validate_touch_position(
	t *[][]byte,
	bt *[][]int,
	table_height,
	table_width,
	check_y,
	check_x int,
) (
	found bool,
	neighbor_element_borders_reached int,
) {

	if validate_in_boundaries_positon(table_height, table_width, check_y, check_x) {
		if (*t)[check_y][check_x] == 'G' ||
			(*t)[check_y][check_x] == '~' {

			deltas := []delta{
				{-2, +0},
				{+2, +0},
				{-1, -1},
				{-1, +1},
				{+1, -1},
				{+1, +1},
			}
			for _, delta := range deltas {
				if validate_in_boundaries_positon(table_height, table_width, check_y+delta.y, check_x+delta.x) {
					if (*t)[check_y+delta.y][check_x+delta.x] == '*' {

						if (*t)[check_y][check_x] == 'G' {
							replace_what = 'G'
						} else {
							replace_what = '~'
						}

						log.Printf(
							"FOUND REVEALED AREAS TOUCHING COORD at (%d;%d) with NEIGHBOR CROSSED BORDERS==%d",
							check_y+delta.y,
							check_x+delta.x,
							(*bt)[check_y+delta.y][check_x+delta.x],
						)
						return true, (*bt)[check_y+delta.y][check_x+delta.x]

					}
				}
			}

		}
	}

	return false, -1
}

func validate_in_boundaries_positon(
	table_height,
	table_width,
	check_y,
	check_x int,
) (
	valid bool,
) {
	if 0 <= check_y && check_y < table_height &&
		0 <= check_x && check_x < table_width {
		return true
	}
	return false
}
