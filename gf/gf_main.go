package gf

import (
	"fmt"
)

//-------------------------------------------------------------------------------
type GFruntime struct {
	
}

//-------------------------------------------------------------------------------
func Init() *GFruntime {
	fmt.Printf("GF ===========>\n")
	runtime := &GFruntime{}
	return runtime
}