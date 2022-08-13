package modules

import (
	"fmt"
	"image"
	"math/rand"
	"sync"
)

const WIDTH float64 = 1024
const HEIGHT float64 = 768

type ImageFunction func() (image.Image, error)

type FunctionPoolStruct struct {
	sync.RWMutex
	functions map[string]ImageFunction
	keys      []string
}

var FunctionPool FunctionPoolStruct

func init() {
	FunctionPool = FunctionPoolStruct{}
	FunctionPool.functions = make(map[string]ImageFunction)
	FunctionPool.keys = make([]string, 0)
}

func (f *FunctionPoolStruct) Add(key string, value ImageFunction) {
	// If the functions are a null map, set up another thread to wait
	// until they...aren't...wait what?
	if f.functions == nil {
		go func(key string, value ImageFunction) {
			for f.functions == nil {
			}
			f.add(key, value)
		}(key, value)
	} else {
		f.add(key, value)
	}
}

func (f *FunctionPoolStruct) add(key string, value ImageFunction) {
	f.Lock()
	f.functions[key] = value
	f.keys = append(f.keys, key)
	f.Unlock()
}

func (f *FunctionPoolStruct) Get(key string) ImageFunction {
	return f.functions[key]
}

func (f FunctionPoolStruct) Random() ImageFunction {
	randomIndex := rand.Intn(len(f.keys))
	key := f.keys[randomIndex]
	fmt.Printf("generating %v gradient...\n", key)
	return f.functions[key]
}
