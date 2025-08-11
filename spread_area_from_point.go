package main

type delta struct {
	y int
	x int
}

func spread_area_from_point_helper(
	t *[][]byte,
	bt *[][]int,
	vt *[][]bool,
	table_height,
	table_width,
	from_y,
	from_x,
	current_element_borders_reached int,
) {

	if !destination_found {

		//log.Printf("SPREAD STARTED AT (%d;%d) replacing %c -> %c", from_y, from_x, replace_what, replace_with)

		deltas := []delta{
			{+0, +0},
			{-2, +0},
			{+2, +0},
			{-1, -1},
			{-1, +1},
			{+1, -1},
			{+1, +1},
		}

		for _, delta := range deltas {

			if !destination_found {

				spread_area_from_point(
					t,
					bt,
					vt,
					table_height,
					table_width,
					from_y+delta.y,
					from_x+delta.x,
					current_element_borders_reached,
				)
			}
		}
	}
}

func spread_area_from_point(
	t *[][]byte,
	bt *[][]int,
	vt *[][]bool,
	table_height,
	table_width,
	from_y,
	from_x int,
	current_element_borders_reached int,
) {

	if !destination_found {
		if from_y < 0 || table_height-1 < from_y ||
			from_x < 0 || table_width-1 < from_x {

			if !already_sailed {
				go_sailing_spread_from_border(
					t,
					bt,
					vt,
					table_height,
					table_width,
					current_element_borders_reached,
				)
			}
		} else {

			if !(*vt)[from_y][from_x] {

				switch (*t)[from_y][from_x] {
				case 'X':

					element_borders_reached = current_element_borders_reached
					if replace_what == '~' {
						//log.Printf("FOUND DESTINATION FROM WATER!!!")
						element_borders_reached++
					}
					destination_found = true
					/* 				//log.Printf(
					"######################\nDESTINATION FOUND at (%d;%d)!!!",
					from_y, from_x) */

				case replace_what:

					(*t)[from_y][from_x] = replace_with
					(*bt)[from_y][from_x] = current_element_borders_reached

					(*vt)[from_y][from_x] = true

					spread_area_from_point_helper(
						t,
						bt,
						vt,
						table_height,
						table_width,
						from_y,
						from_x,
						current_element_borders_reached,
					)

				default:
				}

				(*vt)[from_y][from_x] = true
			}
		}
	}
}

func go_sailing_spread_from_border(
	t *[][]byte,
	bt *[][]int,
	vt *[][]bool,
	table_height,
	table_width,
	current_element_borders_reached int,
) {

	if !destination_found && replace_what == '~' {

		//log.Printf("~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~")

		already_sailed = true

		for i := range table_height {
			for j := range table_width {

				if i == 0 || j == 0 ||
					i == table_height-1 || j == table_width-1 ||
					(i == 1) {
					if (*t)[i][j] == '~' && !(*vt)[i][j] {
						spread_area_from_point_helper(
							t,
							bt,
							vt,
							table_height,
							table_width,
							i,
							j,
							current_element_borders_reached,
						)

						(*vt)[i][j] = true
					}
				}

			}
		}

		//log.Printf("~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~")
	}
}
