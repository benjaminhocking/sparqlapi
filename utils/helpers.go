package utils

import (
	"os"
)

func ReadFile(path string) string{
	dat, err := os.ReadFile(path)
	if err != nil{
		return "500: Internal Server Error"
	}
	return string(dat)
}

func WriteFile(path string, cont string){
	dat := []byte(cont)
	os.WriteFile(path, dat, 0644)
}