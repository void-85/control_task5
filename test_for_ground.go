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

	// H line TOP+BOTTOM
	for i := range hex_w {
		if (*ct)[ct_y][ct_x+hex_h+i] != '_' || (*ct)[ct_y+hex_h+hex_h][ct_x+hex_h+i] != '_' {
			return false, 0
		}
	}

	// DIAG line LEFT TOP+BOTTOM + RIGHT BOTTOM+TOP
	for i := range hex_h {

		if (*ct)[ct_y+hex_h-i][ct_x+i] != '/' || (*ct)[ct_y+hex_h+i+1][ct_x+i] != '\\' {
			return false, 0
		}

		if (*ct)[ct_y+hex_h-i][ct_x+hex_h+hex_h+hex_w-i-1] != '\\' || (*ct)[ct_y+hex_h+i+1][ct_x+hex_h+hex_h+hex_w-i-1] != '/' {
			return false, 0
		}
	}

	for i := range hex_h {
		for j := ct_x + i + 1; j < ct_x+hex_h+hex_h+hex_w-i-1; j++ {
			if (*ct)[ct_y+hex_h-i][j] == 'X' {
				num_of_src_dst_points++
			}

			if (*ct)[ct_y+hex_h+i+1][j] == 'X' {
				num_of_src_dst_points++
			}
		}
	}

	return true, num_of_src_dst_points
}
