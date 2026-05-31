package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Place your code here.
	fmt.Println("Параметры программы:")
	fmt.Printf("Файл для чтения (-from): %s\n", from)
	fmt.Printf("Файл для записи (-to): %s\n", to)
	fmt.Printf("Лимит байт (-limit): %d\n", limit)
	fmt.Printf("Смещение (-offset): %d\n", offset)
	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
	}
}
