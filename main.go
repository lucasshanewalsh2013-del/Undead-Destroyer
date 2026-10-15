package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lucaswalsh11/UndeadDestroyer/imageloading"
)



func main() {
	fmt.Println("Press Enter to load the character model...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	loading.OpenCharacterModel()
	
}
