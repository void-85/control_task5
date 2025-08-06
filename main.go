package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	test_dir = "tests"
)

var (
	already_sailed               = false
	destination_found            = false
	spread_start_y               = -1
	spread_start_x               = -1
	element_borders_reached      = 0
	replace_what            byte = '_'
	replace_with            byte = '_'

	tests_passed = 0
	tests_total  = 0
)

func main() {

	log.Printf("\033[33m[ STARTED ]")

	test_input_files, err := filepath.Glob(filepath.Join(test_dir, "*"))
	if err != nil {
		log.Printf("\033[31mfailed to list input files: %v", err)
	}

	for _, in_file := range test_input_files {

		file_info, err := os.Stat(in_file)
		if err != nil || !file_info.Mode().IsRegular() || strings.Contains(file_info.Name(), ".a") {
			continue
		}

		log.Printf("\033[37mtesting file \"%s\"\n", in_file)
		tests_total++
		out_file := in_file + ".a"

		input, err := os.ReadFile(in_file)
		if err != nil {
			log.Printf("\033[31mfailed to read input file %s: %v\n", in_file, err)
			continue
		}

		output, err := os.ReadFile(out_file)
		if err != nil {
			log.Printf("\033[31mfailed to read output file %s: %v\n", out_file, err)
			continue
		}

		std__in := os.Stdin
		std_out := os.Stdout
		read_in_pipe, write_in_pipe, _ := os.Pipe()

		go func() {
			defer write_in_pipe.Close()
			write_in_pipe.Write(input)
		}()

		os.Stdin = read_in_pipe
		read_out_pipe, write_out_pipe, _ := os.Pipe()
		os.Stdout = write_out_pipe
		start_time := time.Now()

		test_func()

		write_out_pipe.Close()
		time_elapsed := time.Since(start_time)
		var buf bytes.Buffer
		io.Copy(&buf, read_out_pipe)
		os.Stdin = std__in
		os.Stdout = std_out

		actual_output := strings.TrimSpace(buf.String())
		expected_output := strings.TrimSpace(string(output))

		if actual_output != expected_output {

			actual_output_splitted := strings.Split(actual_output, "\n")
			expected_output_splitted := strings.Split(string(output), "\n")

			log.Printf("\033[31mFAILED %s (worked %s)\033[34m", in_file, time_elapsed)
			log.Printf("\033[35mEXP\tACT\tLINE #\033[34m")
			for i := range int(math.Max(float64(len(actual_output_splitted)), float64(len(expected_output_splitted)))) {

				a, b, color := "", "", ""
				if i < len(expected_output_splitted) {
					a = expected_output_splitted[i]
				}
				if i < len(actual_output_splitted) {
					b = actual_output_splitted[i]
				}

				if a == b {
					color = "\033[32m"
				} else {
					color = "\033[31m"
				}

				if a != "" && b != "" {
					log.Printf("%s%s\t%s\t(#%d)\033[34m", color, a, b, i+1)
				}
			}

			/* if len(actual_output_splitted) == len(expected_output_splitted) {
				log.Printf("\n\033[31mFAILED %s (worked %s):", time_elapsed)
				for i := range len(actual_output_splitted) {
					log.Printf("%s\t\t<-must be--\t\t%s", actual_output_splitted[i], expected_output_splitted[i])
				}
			} else {
				log.Printf("\n\033[31mFAILED %s (worked %s)\nExpected:\n%s\nGot:\n%s\n", in_file, time_elapsed, expected_output, actual_output)
			} */

		} else {
			tests_passed++
			log.Printf("\033[32mPASSED %s 	(worked %s)\n", in_file, time_elapsed)
		}

	}

	log.Printf("\033[33m[ FINISHED ]\n\n")
	log.Printf("TESTS : %d / %d passed", tests_passed, tests_total)
}

func test_func() {

	inp := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var datasets int
	fmt.Fscanln(inp, &datasets)

	//log.Printf("there are %d datasets\n", datasets)

DATASETS_LOOP:
	for range datasets {

		var n, m int
		fmt.Fscanln(inp, &n, &m)

		//log.Printf("n== %d m == %d\n", n, m)

		char_table := make([][]byte, n)
		for i := range n {
			char_table[i] = make([]byte, m)
		}

		// adding data to table
		for y := range n {

			chars, _ := inp.ReadString('\n') ///FUCK spaces trailing!!!!
			//log.Printf("line is %d chars length\n", len(chars))

			for x := range m {

				//log.Printf("reading (%d;%d)\n", y, x)

				char_table[y][x] = chars[x]
			}
		}

		var y1, y2, x1, x2 int
		fmt.Fscanln(inp, &y1, &x1)
		fmt.Fscanln(inp, &y2, &x2)

		if y1 == y2 && x1 == x2 {
			char_table[y1-1][x1-1] = 'B'
		} else {
			char_table[y1-1][x1-1] = 'X'
			char_table[y2-1][x2-1] = 'X'
		}

		//---- PRINT char_table --------------------------------------------------
		s := ""
		for y := range n {
			for x := range m {
				s += string(char_table[y][x]) + " "
			}
			s += "\n"
		}
		log.Printf("---------- NEW DATASET ---------------------------\nN == %d   M == %d\nloaded char_table:\n%s\nCHECK PATH: (%d;%d) --> (%d;%d)", n, m, s, y1, x1, y2, x2)
		//---- PRINT char_table --------------------------------------------------

		if y1 == y2 && x1 == x2 {
			log.Printf("same coords inside same hex, breaking dataset computations...")
			fmt.Fprintf(out, "0\n")
			continue DATASETS_LOOP
		}

		first_hex_skipped, hex_height, hex_width := get_hexagon_height_width(&char_table, n, m)

		num_of_y_hexes := (n - 1) / (hex_height + hex_height)
		num_of_x_hexes := (m - hex_height) / (hex_height + hex_width)

		log.Printf("reading %d x %d hex field (first hex is skipped : %t)", num_of_y_hexes, num_of_x_hexes, first_hex_skipped)

		table_height := num_of_y_hexes * 2
		table_width := num_of_x_hexes

		table := make([][]byte, table_height)
		for i := range table_height {
			table[i] = make([]byte, table_width)
		}
		for i := range table_height {
			for j := range table_width {
				table[i][j] = ' '
			}
		}

		borders_table := make([][]uint8, table_height)
		for i := range table_height {
			borders_table[i] = make([]uint8, table_width)
		}

		destination_found,
			spread_start_y,
			spread_start_x =
			map_hex_table(
				&table,
				num_of_y_hexes*2,
				num_of_x_hexes,
				hex_height,
				hex_width,
				&char_table,
				first_hex_skipped,
			)

		if destination_found {
			log.Printf("src & dst inside same hex, breaking dataset computations...")
			fmt.Fprintf(out, "0\n")
			continue DATASETS_LOOP
		}

		log.Printf(
			"starting spreading from (%d;%d)",
			spread_start_y,
			spread_start_x,
		)

		table[spread_start_y][spread_start_x] = '*'
		borders_table[spread_start_y][spread_start_x] = 0
		print_table(&table, &borders_table, table_height, table_width)

		/*
			// destination_found is FALSE from "map_hex_table" func
			already_sailed = false
			element_borders_reached = 0
			replace_what = 'G'
			replace_with = '*'

			for !destination_found {

				spread_from_point(
					&table,
					table_height,
					table_width,
					spread_start_y,
					spread_start_x,
				)

				switch replace_what {
				case 'G':
					replace_what = '~'
				case '~':
					replace_what = 'G'
				default:
					replace_what = '_'
				}

				log.Printf("CHANGED : NOW replacing %c -> %c", replace_what, replace_with)

				if !destination_found {
					element_borders_reached++

					log.Printf(
						"at current table state INCREASING 'element_borders_reached' to %d\n#########################\n",
						element_borders_reached)

					print_table(&table, table_height, table_width)
					log.Printf("\n#########################\n")

					if element_borders_reached > 200 {
						log.Printf("ERROR !!! ERROR !!! ERROR !!! ERROR !!! ERROR !!! ERROR !!!")
						log.Printf("too much borders crossed, possibly error! check!")
						break DATASETS_LOOP
					}
				}
			} */

		fmt.Fprintf(out, "%d\n", element_borders_reached)
	}
}

func print_table(t *[][]byte, bt *[][]uint8, table_height, table_width int) {

	s := ""
	for i := range table_height {
		for j := range table_width {
			s += string((*t)[i][j])
			if (*t)[i][j] == '*' {
				s += fmt.Sprintf("|%d ", (*bt)[i][j])
			} else {
				s += "   "
			}
		}
		s += "\n"
	}
	log.Printf("\nloaded table:\n%s", s)

}

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

func get_hexagon_height_width(char_table *[][]byte, n, m int) (first_hex_skipped bool, height, width int) {

	first_hex_skipped = false
	log.Printf("examining table %dx%d", n, m)

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

							log.Printf("found hex(%dx%d) at point(%d;%d)", h, w, y, x)
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
		log.Printf("hex size not found!")
	} else {
		log.Printf("h==%d w==%d hex size found", height, width)
	}

	return first_hex_skipped, height, width

}
