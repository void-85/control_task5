package main

func test_coords_for_ground_hex_and_check_for_src_dst(
	ct *[][]byte,
	test_y,
	test_x,
	hex_h,
	hex_w int,
	lower_position bool,
) (
	is_ground bool,
	num_of_src_dst_points int,
) {

	ct_y := test_y * (hex_h + hex_h)
	ct_x := test_x * (hex_h + hex_w)

	if lower_position {
		ct_y += hex_h
	}

	// CRITICAL EDGE CASE - NO LOWER HEXES IN CHAR_TABLE but trying to move over the boundaries
	if ct_y+hex_h+hex_h >= len(*ct) {
		return false, 0
	}

	passed_chars := 0
	passed_lines := 0

	// H line TOP+BOTTOM
	passed_chars = 0
	for i := range hex_w {
		if (*ct)[ct_y][ct_x+hex_h+i] == '_' && (*ct)[ct_y+hex_h+hex_h][ct_x+hex_h+i] == '_' {
			passed_chars++
		}
	}
	if passed_chars == hex_w {
		passed_lines += 2
	}

	// DIAG line LEFT TOP+BOTTOM + RIGHT BOTTOM+TOP
	passed_chars = 0
	for i := range hex_h {

		if (*ct)[ct_y+hex_h-i][ct_x+i] == '/' && (*ct)[ct_y+hex_h+i+1][ct_x+i] == '\\' {
			passed_chars++
		}

		if (*ct)[ct_y+hex_h-i][ct_x+hex_h+hex_h+hex_w-i-1] == '\\' && (*ct)[ct_y+hex_h+i+1][ct_x+hex_h+hex_h+hex_w-i-1] == '/' {
			passed_chars++
		}

		for j := ct_x + i + 1; j < ct_x+hex_h+hex_h+hex_w-i-1; j++ {
			if (*ct)[ct_y+hex_h-i][j] == 'X' {
				num_of_src_dst_points++
			}

			if (*ct)[ct_y+hex_h+i+1][j] == 'X' {
				num_of_src_dst_points++
			}
		}
	}

	if passed_chars == 2*hex_h {
		passed_lines += 4
	}

	if passed_lines == 6 {
		is_ground = true
	}

	return is_ground, num_of_src_dst_points
}
