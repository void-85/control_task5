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
)

const (
	test_dir = "tests"
)

func main() {

	fmt.Printf("\n\033[33m[ STARTED ]\n")

	test_input_files, err := filepath.Glob(filepath.Join(test_dir, "*"))
	if err != nil {
		fmt.Printf("\033[31mfailed to list input files: %v", err)
	}

	for _, in_file := range test_input_files {

		file_info, err := os.Stat(in_file)
		if err != nil || !file_info.Mode().IsRegular() || strings.Contains(file_info.Name(), ".a") {
			continue
		}

		fmt.Printf("\033[37mtesting file \"%s\"\n", in_file)
		out_file := in_file + ".a"

		input, err := os.ReadFile(in_file)
		if err != nil {
			fmt.Printf("\033[31mfailed to read input file %s: %v\n", in_file, err)
			continue
		}

		expectedOutput, err := os.ReadFile(out_file)
		if err != nil {
			fmt.Printf("\033[31mfailed to read output file %s: %v\n", out_file, err)
			continue
		}

		origStdin := os.Stdin
		origStdout := os.Stdout
		rIn, wIn, _ := os.Pipe()
		wIn.Write(input)
		wIn.Close()
		os.Stdin = rIn
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut
		start := time.Now()
		test_func()
		wOut.Close()
		duration := time.Since(start)
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		os.Stdin = origStdin
		os.Stdout = origStdout

		actualOutput := strings.TrimSpace(buf.String())
		expected := strings.TrimSpace(string(expectedOutput))

		if actualOutput != expected {
			fmt.Printf("\n\033[31mFAILED %s (worked %s)\nExpected:\n%s\n--------------------------------------------------------------\nGot:\n%s\n", in_file, duration, expected, actualOutput)
		} else {
			fmt.Printf("\033[32mPASSED %s 	(worked %s)\n", in_file, duration)
		}

	}

	fmt.Printf("\033[33m[ FINISHED ]\n\n")
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

		table := make([][]byte, n)
		for i := range n {
			table[i] = make([]byte, m)
		}

		// adding data to table
		for y := range n {

			chars, _ := inp.ReadString('\n') ///FUCK spaces trailing!!!!
			//log.Printf("line is %d chars length\n", len(chars))

			for x := range m {

				//log.Printf("reading (%d;%d)\n", y, x)

				table[y][x] = chars[x]
			}
		}

		var y1, y2, x1, x2 int
		fmt.Fscanln(inp, &y1, &x1)
		fmt.Fscanln(inp, &y2, &x2)

		s := ""
		for y := range n {

			for x := range m {
				if y == y1-1 && x == x1-1 && y == y2-1 && x == x2-1 {
					s += "B" + " "
				} else if y == y1-1 && x == x1-1 {
					s += "F" + " "
				} else if y == y2-1 && x == x2-1 {
					s += "T" + " "
				} else {
					s += string(table[y][x]) + " "
				}
			}
			s += "\n"
		}

		log.Printf("-------------------------------------\nN == %d   M == %d\nloaded table:\n%s\nCHECK PATH: (%d;%d) --> (%d;%d)", n, m, s, y1, x1, y2, x2)

		first_hex_skipped, hex_height, hex_width := get_hexagon_height_width(&table, n, m)

		num_of_y_hexes := (n - 1) / (hex_height + hex_height)
		num_of_x_hexes := (m - hex_height) / (hex_height + hex_width)

		log.Printf("reading %d x %d hex field (first hex is skipped : %t)", num_of_y_hexes, num_of_x_hexes, first_hex_skipped)

	}

	fmt.Fprintf(out, "%s\n", "asd")
}

func get_hexagon_height_width(table *[][]byte, n, m int) (first_hex_skipped bool, height, width int) {

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
							if (*table)[y][x+h+i] == '_' && (*table)[y+h+h][x+h+i] == '_' {
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
							if (*table)[y+h-i][x+i] == '/' && (*table)[y+h+i+1][x+i] == '\\' {
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
							if (*table)[y+h+h-i][x+h+i+w] == '/' && (*table)[y+1+i][x+h+i+w] == '\\' {
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
		log.Printf("searching upto h==%d w==%d failed", n, m)
	} else {
		log.Printf("h==%d w==%d PASSED", height, width)
	}

	return first_hex_skipped, height, width

}
