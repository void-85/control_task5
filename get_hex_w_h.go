package main

func get_hexagon_height_width(char_table *[][]byte, n, m int) (first_hex_skipped bool, height, width int) {

	first_hex_skipped = false
	//log.Printf("examining table %dx%d", n, m)

MAIN_LOOP:
	for h := range n + 1 {
		for w := range m + 1 {

			if 1 <= w && 1 <= h {
				for y := range n - h - h - 1 + 1 {
					for x := range m - w - h - h + 1 {

						passed_chars := 0
						passed_lines := 0

						// H line TOP+BOTTOM
						passed_chars = 0
						for i := range w {
							if (*char_table)[y][x+h+i] == '_' && (*char_table)[y+h+h][x+h+i] == '_' {
								passed_chars++
							}
						}
						if passed_chars == w {
							passed_lines += 2
						}

						// DIAG line LEFT TOP+BOTTOM
						passed_chars = 0
						for i := range h {
							if (*char_table)[y+h-i][x+i] == '/' && (*char_table)[y+h+i+1][x+i] == '\\' {
								passed_chars++
							}
						}
						if passed_chars == h {
							passed_lines += 2
						}

						// DIAG line RIGHT BOTTOM+TOP
						passed_chars = 0
						for i := range h {
							if (*char_table)[y+h+h-i][x+h+i+w] == '/' && (*char_table)[y+1+i][x+h+i+w] == '\\' {
								passed_chars++
							}
						}
						if passed_chars == h {
							passed_lines += 2
						}

						if passed_lines == 6 {

							//log.Printf("found hex(%dx%d) at point(%d;%d)", h, w, y, x)
							if y%(1+h+h) == 0 && x%(h+w+h+w) == 0 {
								first_hex_skipped = false
							} else {
								first_hex_skipped = true
							}

							height = h
							width = w
							break MAIN_LOOP
						}
					}
				}
			}
		}
	}

	if height == 0 && width == 0 {
		//log.Printf("hex size not found!")
	} else {
		//log.Printf("h==%d w==%d hex size found", height, width)
	}

	return first_hex_skipped, height, width

}
