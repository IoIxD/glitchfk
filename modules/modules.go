package modules

import (
	"fmt"
	"image"
	"math/rand"
	"sync"
	"time"
)

const NORMAL_WIDTH float64 = 1024
const NORMAL_HEIGHT float64 = 768

const SUNDAY_WIDTH float64 = 1920
const SUNDAY_HEIGHT float64 = 1080

var WIDTH float64
var HEIGHT float64

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

	if(time.Now().Weekday() == 0) {
		WIDTH = SUNDAY_WIDTH
		HEIGHT = SUNDAY_HEIGHT
	} else {
		WIDTH = NORMAL_WIDTH
		HEIGHT = NORMAL_HEIGHT
	}
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
	if(ok) {
		return function
	} else {
		return nil
	}
}

func (f *FunctionPoolStruct) Random() ImageFunction {
	f.Lock()
	randomIndex := rand.Intn(len(f.keys))
	key := f.keys[randomIndex]
	fmt.Printf("generating %v image...\n", key)
	f.Unlock()
	return f.functions[key]
}
