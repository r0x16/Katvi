package main

import "github.com/r0x16/Katvi/src/shared/infraestructure"

func main() {
	application := &infraestructure.Main{}
	application.RunServices()
}
