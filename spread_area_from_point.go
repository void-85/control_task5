package main

type delta struct {
	y int
	x int
}

func spread_area_from_point_helper(
	table *[][]byte,
	borders_table *[][]int,
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
					table,
					borders_table,
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
	table *[][]byte,
	borders_table *[][]int,
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
					table,
					borders_table,
					table_height,
					table_width,
					current_element_borders_reached,
				)
			}
		} else {

			switch (*table)[from_y][from_x] {
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

				(*table)[from_y][from_x] = replace_with
				(*borders_table)[from_y][from_x] = current_element_borders_reached

				/* 				//log.Printf(
					"replaced %c -> %c at (%d;%d):",
					replace_what,
					replace_with,
					from_y,
					from_x,
				) */
				//print_table(table, table_height, table_width)

				spread_area_from_point_helper(
					table,
					borders_table,
					table_height,
					table_width,
					from_y,
					from_x,
					current_element_borders_reached,
				)

			default:

				/*if (*table)[check_y][check_x] != replace_with &&
					(*table)[check_y][check_x] != ' ' {
					spread_start_y = check_y
					spread_start_x = check_x
					//log.Printf("### spread_start set to (%d;%d)", spread_start_y, spread_start_x)
				}*/
			}
		}
	}
}

func go_sailing_spread_from_border(
	table *[][]byte,
	borders_table *[][]int,
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
					if (*table)[i][j] == '~' {
						spread_area_from_point_helper(
							table,
							borders_table,
							table_height,
							table_width,
							i,
							j,
							current_element_borders_reached,
						)
					}
				}

			}
		}

		//log.Printf("~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~")
	}
}
