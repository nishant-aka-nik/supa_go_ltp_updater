package main

import (
	"fmt"

	stockrate "github.com/shubhamgosain/stockrate"
)

func main() {
	fmt.Println(stockrate.GetPrice("mahanagar gas"))

	// fmt.Println(stockrate.GetCompanyList())

}
