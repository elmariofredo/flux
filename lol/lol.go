package main

import (
	"github.com/Masterminds/semver"
	"fmt"
)

func main() {
	b, _ := semver.NewVersion("1.13.9")
	a, _ := semver.NewVersion("1.13.9-perl")
	c, _ := semver.NewVersion("1.13.9-alpine")
	fmt.Println(a.Compare(b)) // a < b
	fmt.Println(b.Compare(c)) // b > c
	fmt.Println(a.Compare(c)) // a > c
}
