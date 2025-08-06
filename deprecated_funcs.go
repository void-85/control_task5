package main

import "log"

func spread_from_point(
	table *[][]byte,
	table_height,
	table_width,
	from_y,
	from_x int,
) {

	if !destination_found {

		log.Printf("SPREAD STARTED AT (%d;%d) replacing %c -> %c", from_y, from_x, replace_what, replace_with)

		var dx, dy int

		dy, dx = 0, 0
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = -2, 0
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = +2, 0
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = -1, -1
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = -1, +1
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = +1, -1
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}

		dy, dx = +1, +1
		if !destination_found {
			spread_from_point_check_replce(table, table_height, table_width, from_y+dy, from_x+dx)
		}
	}
}

func spread_from_point_check_replce(
	table *[][]byte,
	table_height,
	table_width,
	check_y,
	check_x int,
) {

	if !destination_found {
		if check_y < 0 || table_height-1 < check_y ||
			check_x < 0 || table_width-1 < check_x {

			if !already_sailed {
				go_sailing_spread_from_border(table, table_height, table_width)
			}
		} else {

			switch (*table)[check_y][check_x] {
			case 'X':
				if replace_what == '~' {
					log.Printf("FOUND DESTINATION FROM WATER!!!")
					element_borders_reached++
				}
				destination_found = true
				log.Printf(
					"######################\nDESTINATION FOUND at (%d;%d)!!!",
					check_y, check_x)

			case replace_what:

				(*table)[check_y][check_x] = replace_with

				log.Printf(
					"replaced %c -> %c at (%d;%d):",
					replace_what,
					replace_with,
					check_y,
					check_x,
				)
				//print_table(table, table_height, table_width)

				spread_from_point(table, table_height, table_width, check_y, check_x)

			default:

				/*if (*table)[check_y][check_x] != replace_with &&
					(*table)[check_y][check_x] != ' ' {
					spread_start_y = check_y
					spread_start_x = check_x
					log.Printf("### spread_start set to (%d;%d)", spread_start_y, spread_start_x)
				}*/
			}
		}
	}
}

func go_sailing_spread_from_border(
	table *[][]byte,
	table_height,
	table_width int,
) {

	if !destination_found && replace_what == '~' {

		log.Printf("~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~ ~~~ SAILING!!! ~~~")

		already_sailed = true

		for i := range table_height {
			for j := range table_width {

				if i == 0 || j == 0 || i == table_height-1 || j == table_width-1 {
					if (*table)[i][j] == '~' {
						spread_from_point(table, table_height, table_width, i, j)
					}
				}

			}
		}

		log.Printf("~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~ ~~~ END SAILING ~~~")
	}
}
