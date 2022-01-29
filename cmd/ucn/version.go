package main

import (
	"fmt"
	"os"

	"github.com/mgumz/ucn/internal/pkg"
)

func printVersion(_ string) error {
	fmt.Println("ucn", pkg.Version, pkg.GitHash, pkg.BuildDate)
	os.Exit(0)
	return nil
}
