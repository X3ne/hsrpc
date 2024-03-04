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
	fmt.Println("3. Prepare Bosses")
	fmt.Println("4. Execute All Functions")
	fmt.Print("Enter your choice (1, 2, 3, or 4): ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		PrepareCharacters(config)
	case 2:
		PrepareLocations(config)
	case 3:
		PrepareBosses(config)
	case 4:
		PrepareCharacters(config)
		PrepareLocations(config)
		PrepareBosses(config)
	default:
		fmt.Println("Invalid choice. Please enter 1, 2, 3, or 4.")
	}
}
