package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"
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

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	failed_tests_names := ""

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

			failed_tests_names += fmt.Sprintf("failed: %s\n", in_file)

			//actual_output_splitted := strings.Split(actual_output, "\n")
			//expected_output_splitted := strings.Split(string(output), "\n")

			log.Printf("\033[31mFAILED %s (worked %s)\033[34m", in_file, time_elapsed)
			//log.Printf("\033[35mEXP\tACT\tLINE #\033[34m")
			/* for i := range int(math.Max(float64(len(actual_output_splitted)), float64(len(expected_output_splitted)))) {

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
					//log.Printf("%s%s\t%s\t(#%d)\033[34m", color, a, b, i+1)
				}
			} */

			/* if len(actual_output_splitted) == len(expected_output_splitted) {
				//log.Printf("\n\033[31mFAILED %s (worked %s):", time_elapsed)
				for i := range len(actual_output_splitted) {
					//log.Printf("%s\t\t<-must be--\t\t%s", actual_output_splitted[i], expected_output_splitted[i])
				}
			} else {
				//log.Printf("\n\033[31mFAILED %s (worked %s)\nExpected:\n%s\nGot:\n%s\n", in_file, time_elapsed, expected_output, actual_output)
			} */

		} else {
			tests_passed++
			log.Printf("\033[32mPASSED %s 	(worked %s)\n", in_file, time_elapsed)
		}

	}

	log.Printf("\033[33m[ FINISHED ]\n\n")
	log.Printf("TESTS : %d / %d passed", tests_passed, tests_total)
	log.Printf("%s", failed_tests_names)
}

func test_func() {

	inp := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var datasets int
	fmt.Fscanln(inp, &datasets)

	////log.Printf("there are %d datasets\n", datasets)

DATASETS_LOOP:
	for range /*cur_dataset := */ datasets {

		var n, m int
		fmt.Fscanln(inp, &n, &m)

		////log.Printf("n== %d m == %d\n", n, m)

		char_table := make([][]byte, n)
		for i := range n {
			char_table[i] = make([]byte, m)
		}

		// adding data to table
		for y := range n {

			chars, _ := inp.ReadString('\n') ///FUCK spaces trailing!!!!
			////log.Printf("line is %d chars length\n", len(chars))

			for x := range m {

				////log.Printf("reading (%d;%d)\n", y, x)

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
		/* 		s := ""
		   		for y := range n {
		   			for x := range m {
		   				s += string(char_table[y][x]) + " "
		   			}
		   			s += "\n"
		   		} */
		/* 		//log.Printf(
			"---------- %d DATASET ---------------------------\nN == %d   M == %d\nloaded char_table:\n%s\nCHECK PATH: (%d;%d) --> (%d;%d)",
			cur_dataset+1,
			n, m, s, y1, x1, y2, x2,
		) */
		//---- PRINT char_table --------------------------------------------------

		if y1 == y2 && x1 == x2 {
			//log.Printf("same coords inside same hex, breaking dataset computations...")
			fmt.Fprintf(out, "0\n")
			continue DATASETS_LOOP
		}

		first_hex_skipped, hex_height, hex_width := get_hexagon_height_width(&char_table, n, m)

		num_of_y_hexes := (n - 1) / (hex_height + hex_height)
		num_of_x_hexes := (m - hex_height) / (hex_height + hex_width)

		//log.Printf("reading %d x %d hex field (first hex is skipped : %t)", num_of_y_hexes, num_of_x_hexes, first_hex_skipped)

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

		borders_table := make([][]int, table_height)
		for i := range table_height {
			borders_table[i] = make([]int, table_width)
		}

		visited_table := make([][]bool, table_height)
		for i := range table_height {
			visited_table[i] = make([]bool, table_width)
		}
		reset_visited_table(&visited_table, table_height, table_width)

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

		// TODO: FREE char_table MEMORY

		if destination_found {
			//log.Printf("src & dst inside same hex, breaking dataset computations...")
			fmt.Fprintf(out, "0\n")
			continue DATASETS_LOOP
		}

		/* 		//log.Printf(
			"starting spreading from (%d;%d)",
			spread_start_y,
			spread_start_x,
		) */

		table[spread_start_y][spread_start_x] = 'G'
		borders_table[spread_start_y][spread_start_x] = 0
		replace_what = 'G'
		replace_with = '*'
		already_sailed = false
		element_borders_reached = 0

		//print_table(&table, &borders_table, table_height, table_width)

		//reset_visited_table(&visited_table, table_height, table_width)
		spread_area_from_point_helper(
			&table,
			&borders_table,
			&visited_table,
			table_height,
			table_width,
			spread_start_y,
			spread_start_x,
			element_borders_reached)

		for !destination_found {

			reset_visited_table(&visited_table, table_height, table_width)
			spread_start_y,
				spread_start_x,
				element_borders_reached =
				find_coords_touching_revealed_areas(
					&table,
					&borders_table,
					&visited_table,
					table_height,
					table_width,
					spread_start_y,
					spread_start_x,
				)

			reset_visited_table(&visited_table, table_height, table_width)
			spread_area_from_point_helper(
				&table,
				&borders_table,
				&visited_table,
				table_height,
				table_width,
				spread_start_y,
				spread_start_x,
				element_borders_reached)

			//log.Printf("iteration passed")
			//print_table(&table, &borders_table, table_height, table_width)

		}

		//print_table(&table, &borders_table, table_height, table_width)

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

				//log.Printf("CHANGED : NOW replacing %c -> %c", replace_what, replace_with)

				if !destination_found {
					element_borders_reached++

					//log.Printf(
						"at current table state INCREASING 'element_borders_reached' to %d\n#########################\n",
						element_borders_reached)

					print_table(&table, table_height, table_width)
					//log.Printf("\n#########################\n")

					if element_borders_reached > 200 {
						//log.Printf("ERROR !!! ERROR !!! ERROR !!! ERROR !!! ERROR !!! ERROR !!!")
						//log.Printf("too much borders crossed, possibly error! check!")
						break DATASETS_LOOP
					}
				}
			} */

		fmt.Fprintf(out, "%d\n", element_borders_reached)
	}
}

/* func switch_replacing_what_with() {
	switch replace_what {
	case 'G':
		replace_what = '~'
	case '~':
		replace_what = 'G'
	default:
		replace_what = '_'
	}
} */

func reset_visited_table(vt *[][]bool, table_height, table_width int) {

	for i := range table_height {
		for j := range table_width {
			(*vt)[i][j] = false
		}
	}

}
