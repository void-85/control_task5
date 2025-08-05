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
			log.Printf("\033[32mPASSED %s 	(worked %s)\n", in_file, time_elapsed)
		}

	}

	log.Printf("\033[33m[ FINISHED ]\n\n")
}

func test_func() {

	inp := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var datasets int
	fmt.Fscanln(inp, &datasets)

	//log.Printf("there are %d datasets\n", datasets)

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
				/* if y == y1-1 && x == x1-1 && y == y2-1 && x == x2-1 {
					s += "B" + " "
				} else if y == y1-1 && x == x1-1 {
					s += "F" + " "
				} else if y == y2-1 && x == x2-1 {
					s += "T" + " "
				} else { */
				s += string(char_table[y][x]) + " "
				//}
			}
			s += "\n"
		}
		log.Printf("-------------------------------------\nN == %d   M == %d\nloaded char_table:\n%s\nCHECK PATH: (%d;%d) --> (%d;%d)", n, m, s, y1, x1, y2, x2)
		//---- PRINT char_table --------------------------------------------------

		if y1 == y2 && x1 == x2 {
			log.Printf("same coords ==> inside same hex, breaking dataset computations...")
			fmt.Fprintf(out, "0\n")
		} else {

			first_hex_skipped, hex_height, hex_width := get_hexagon_height_width(&char_table, n, m)

			num_of_y_hexes := (n - 1) / (hex_height + hex_height)
			num_of_x_hexes := (m - hex_height) / (hex_height + hex_width)

			log.Printf("reading %d x %d hex field (first hex is skipped : %t)", num_of_y_hexes, num_of_x_hexes, first_hex_skipped)

			table := make([][]byte, num_of_y_hexes*2)
			for i := range num_of_y_hexes * 2 {
				table[i] = make([]byte, num_of_x_hexes)
			}
			for i := range num_of_y_hexes * 2 {
				for j := range num_of_x_hexes {
					table[i][j] = '.'
				}
			}

			map_hex_table(&table, num_of_y_hexes*2, num_of_x_hexes, hex_height, hex_width, &char_table, first_hex_skipped)

			//---- PRINT table --------------------------------------------------

			s = ""
			for i := range num_of_y_hexes * 2 {
				for j := range num_of_x_hexes {
					s += string(table[i][j]) + "  "
				}
				s += "\n"
			}
			log.Printf("\nloaded table:\n%s", s)
			//---- PRINT table --------------------------------------------------

			fmt.Fprintf(out, "9999\n")
		}
	}

	//fmt.Fprintf(out, "%s\n", "0\n2\n2\n5")
}

func map_hex_table(table *[][]byte, total_y_hexes, total_x_hexes, hex_h, hex_w int, char_table *[][]byte, first_hex_skipped bool) {

	/* // table params according to first _hex_skipped
	start_y_shift := 0
	order_y_shift := +1

	if first_hex_skipped{
		start_y_shift := 1
		order_y_shift := -1
	} */

	for y := range total_y_hexes / 2 {
		for x := range total_x_hexes {

			lower_position := (x%2 == 1)
			if first_hex_skipped {
				lower_position = !lower_position
			}

			//is_ground, num_of_src_dst_points :=
			test_coords_for_ground_hex_and_check_for_src_dst(
				char_table,
				y,
				x,
				hex_h,
				hex_w,
				lower_position,
			)

		}
	}

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

	//----------------------------------------------------------------------------
	passed_chars := 0
	passed_lines := 0

	// H line TOP+BOTTOM
	//-----------------------------------------
	passed_chars = 0
	for i := range hex_w {
		if (*ct)[ct_y][ct_x+hex_h+i] == '_' && (*ct)[ct_y+hex_h+hex_h][ct_x+hex_h+i] == '_' {
			passed_chars++
		}
	}
	if passed_chars == hex_w {
		passed_lines += 2
	}
	//-----------------------------------------

	/* 	// DIAG line LEFT TOP+BOTTOM
	   	//-----------------------------------------
	   	passed_chars = 0
	   	for i := range hex_h {
	   		if (*char_table)[char_table_y+hex_h-i][char_table_x+i] == '/' && (*char_table)[char_table_y+hex_h+i+1][char_table_x+i] == '\\' {
	   			passed_chars++
	   		}
	   	}
	   	if passed_chars == hex_h {
	   		passed_lines += 2
	   	}
	   	//-----------------------------------------

	   	// DIAG line RIGHT BOTTOM+TOP
	   	//-----------------------------------------
	   	passed_chars = 0
	   	for i := range hex_h {
	   		if (*char_table)[char_table_y+hex_h+hex_h-i][char_table_x+hex_h+i+hex_w] == '/' && (*char_table)[char_table_y+1+i][char_table_x+hex_h+i+hex_w] == '\\' {
	   			passed_chars++
	   		}
	   	}
	   	if passed_chars == hex_h {
	   		passed_lines += 2
	   	}
	   	//-----------------------------------------
	*/

	// DIAG line LEFT TOP+BOTTOM + RIGHT BOTTOM+TOP
	//-----------------------------------------
	passed_chars = 0
	for i := range hex_h {

		if (*ct)[ct_y+hex_h-i][ct_x+i] == '/' && (*ct)[ct_y+hex_h+i+1][ct_x+i] == '\\' {
			passed_chars++
		}

		if (*ct)[ct_y+hex_h+hex_h-i][ct_x+hex_h+i+hex_w] == '/' && (*ct)[ct_y+1+i][ct_x+hex_h+i+hex_w] == '\\' {
			passed_chars++
		}

		log.Printf(
			"edges check coords: (%d;%d) (%d;%d) \033[32m(%d;%d) (%d;%d)\033[35m(%d;%d) (%d;%d)",
			ct_y+hex_h-i,
			ct_x+i,
			ct_y+hex_h+i+1,
			ct_x+i,

			ct_y+hex_h+hex_h-i,
			ct_x+hex_h+i+hex_w,
			ct_y+1+i,
			ct_x+hex_h+i+hex_w,

			ct_y+hex_h-i,
			ct_x+hex_h+hex_h+hex_w-i-1,
			ct_y+hex_h+i+1,
			ct_x+hex_h+hex_h+hex_w-i-1,
		)
	}

	if passed_chars == 2*hex_h {
		passed_lines += 4
	}
	//-----------------------------------------

	if passed_lines == 6 {
		is_ground = true
	}
	//----------------------------------------------------------------------------

	log.Printf(
		"testing table(%d;%d) --> char_table(%d, %d)(is lower: %t) \t==> is ground:%t\tPOINTS:%d",
		test_y,
		test_x,
		ct_y,
		ct_x,
		lower_position,
		is_ground,
		num_of_src_dst_points,
	)

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
						//-----------------------------------------
						passed_chars = 0
						for i := range w {
							if (*char_table)[y][x+h+i] == '_' && (*char_table)[y+h+h][x+h+i] == '_' {
								passed_chars++
							}
						}
						if passed_chars == w {
							passed_lines += 2
						}
						//-----------------------------------------

						// DIAG line LEFT TOP+BOTTOM
						//-----------------------------------------
						passed_chars = 0
						for i := range h {
							if (*char_table)[y+h-i][x+i] == '/' && (*char_table)[y+h+i+1][x+i] == '\\' {
								passed_chars++
							}
						}
						if passed_chars == h {
							passed_lines += 2
						}
						//-----------------------------------------

						// DIAG line RIGHT BOTTOM+TOP
						//-----------------------------------------
						passed_chars = 0
						for i := range h {
							if (*char_table)[y+h+h-i][x+h+i+w] == '/' && (*char_table)[y+1+i][x+h+i+w] == '\\' {
								passed_chars++
							}
						}
						if passed_chars == h {
							passed_lines += 2
						}
						//-----------------------------------------

						/* if passed_lines > 3 {
							log.Printf("for w==%d h==%d only %d of 6 are passed", w, h, passed_lines)
						} */

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
