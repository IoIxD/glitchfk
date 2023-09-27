package modules

import (
	"fmt"
	"image"
	"math/rand"
	"sync"
)

type ImageFunction func(seed int64, width, height float64) (image.Image, error)

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
	function, ok := f.functions[key]
	if ok {
		return function
	} else {
		return nil
	}
}

func (f *FunctionPoolStruct) Random() (ImageFunction, string) {
	f.Lock()
	randomIndex := rand.Intn(len(f.keys))
	key := f.keys[randomIndex]
	fmt.Printf("generating %v image...\n", key)
	f.Unlock()
	return f.functions[key], key
}
