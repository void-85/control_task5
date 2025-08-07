package main

func map_hex_table(
	table *[][]byte,
	total_y_hexes,
	total_x_hexes,
	hex_h,
	hex_w int,
	char_table *[][]byte,
	first_hex_skipped bool,
) (
	found_src_dst_in_one_hex bool,
	spread_start_y,
	spread_start_x int,
) {

	for y := range total_y_hexes / 2 {
		for x := range total_x_hexes {

			lower_position := (x%2 == 1)
			if first_hex_skipped {
				lower_position = !lower_position
			}

			is_ground, num_of_src_dst_points :=
				test_coords_for_ground_hex_and_check_for_src_dst(
					char_table,
					y,
					x,
					hex_h,
					hex_w,
					lower_position,
				)

			if num_of_src_dst_points == 2 {
				return true, 0, 0
			}

			lower_shift := 0
			if lower_position {
				lower_shift = 1
			}

			setting_char := 'G'
			if is_ground {

				if num_of_src_dst_points == 1 {
					setting_char = 'X'
					spread_start_y = y*2 + lower_shift
					spread_start_x = x
				}

			} else {
				setting_char = '~'
			}

			(*table)[y*2+lower_shift][x] = byte(setting_char)
		}
	}

	return false, spread_start_y, spread_start_x
}
