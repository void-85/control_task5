package main

type delta struct {
	y int
	x int
}

func spread_area_from_point_helper(
	t *[][]byte,
	vt *[][]bool,
	table_height,
	table_width,
	from_y,
	from_x int,
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
					vt,
					table_height,
					table_width,
					from_y+delta.y,
					from_x+delta.x,
				)
			}
		}
	}
}

func spread_area_from_point(
	t *[][]byte,
	vt *[][]bool,
	table_height,
	table_width,
	from_y,
	from_x int,
) {

	if !destination_found {

		if 0 <= from_y && from_y <= table_height-1 &&
			0 <= from_x && from_x <= table_width-1 {

			if !(*vt)[from_y][from_x] {

				switch (*t)[from_y][from_x] {
				case 'X':

					destination_found = true

				case replace_what:

					(*t)[from_y][from_x] = replace_with
					(*vt)[from_y][from_x] = true

					spread_area_from_point_helper(
						t,
						vt,
						table_height,
						table_width,
						from_y,
						from_x,
					)

				default:
				}

				(*vt)[from_y][from_x] = true
			}
		}
	}
}
