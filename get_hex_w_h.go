package main

/* func get_hexagon_height_width(
	char_table *[][]byte,
	n,
	m int,
) (
	first_hex_skipped bool,
	height,
	width int,
) {

	first_hex_skipped = false
	//log.Printf("examining table %dx%d", n, m)


	// each must be +1
	for h := range n {
		for w := range m {



		}
	}

	return first_hex_skipped, height, width

}

func test_hexagon_size(
	char_table *[][]byte,
	test_height,
	test_width int,
) (
	first_hex_skipped bool,
	height,
	width int,
) {

	ct := *char_table

	for y := range n - h - h  {
		for x := range m - w - h - h + 1 {

			passed_chars := 0
			passed_lines := 0

			// H line TOP+BOTTOM
			passed_chars = 0
			for i := range w {
				if ct[y][x+h+i] == '_' && ct[y+h+h][x+h+i] == '_' {
					passed_chars++
				}
			}
			if passed_chars == w {
				passed_lines += 2
			}

			// DIAG line LEFT TOP+BOTTOM
			passed_chars = 0
			for i := range h {
				if ct[y+h-i][x+i] == '/' && ct[y+h+i+1][x+i] == '\\' {
					passed_chars++
				}
			}
			if passed_chars == h {
				passed_lines += 2
			}

			// DIAG line RIGHT BOTTOM+TOP
			passed_chars = 0
			for i := range h {
				if ct[y+h+h-i][x+h+i+w] == '/' && ct[y+1+i][x+h+i+w] == '\\' {
					passed_chars++
				}
			}
			if passed_chars == h {
				passed_lines += 2
			}

			if passed_lines == 6 {

				//log.Printf("found hex(%dx%d) at point(%d;%d)", h, w, y, x)
				if y & (h+h) == 0 && x & (h+w+h+w-1) == 0 {
					first_hex_skipped = false
				} else {
					first_hex_skipped = true
				}

				height = h
				width = w

				return first_hex_skipped, height, width
			}
		}
	}

} */

func get_hexagon_height_width(
	ct *[][]byte,
	n,
	m int,
) (
	firstHexSkipped bool,
	height,
	width int,
) {

	t := *ct

MAIN_LOOP:
	for h := 1; h <= n; h++ {
		hh := h + h
		for w := 1; w <= m; w++ {

			for y := 0; y <= n-hh-1; y++ {
				for x := 0; x <= m-w-hh; x++ {

					// H line TOP+BOTTOM
					ok := true
					for i := 0; i < w; i++ {
						if t[y][x+h+i] != '_' || t[y+hh][x+h+i] != '_' {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}

					// DIAG line LEFT TOP+BOTTOM
					ok = true
					for i := 0; i < h; i++ {
						if t[y+h-i][x+i] != '/' || t[y+h+i+1][x+i] != '\\' {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}

					// DIAG line RIGHT BOTTOM+TOP
					ok = true
					for i := 0; i < h; i++ {
						if t[y+hh-i][x+h+i+w] != '/' || t[y+1+i][x+h+i+w] != '\\' {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}

					// Found match
					if y%(1+hh) == 0 && x%(h+w+h+w) == 0 {
						firstHexSkipped = false
					} else {
						firstHexSkipped = true
					}

					height = h
					width = w
					break MAIN_LOOP
				}
			}
		}
	}
	return
}
