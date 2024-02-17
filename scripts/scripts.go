package main

import (
	"fmt"
)

func main() {
	config := InitConfig()

	var choice int
	fmt.Println("Choose an option:")
	fmt.Println("1. Prepare Characters")
	fmt.Println("2. Prepare Locations")
	fmt.Println("3. Execute All Functions")
	fmt.Print("Enter your choice (1, 2, or 3): ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		PrepareCharacters(config)
	case 2:
		PrepareLocations(config)
	case 3:
		PrepareCharacters(config)
		PrepareLocations(config)
	default:
		fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
	}
}
