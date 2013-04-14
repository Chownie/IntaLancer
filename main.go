package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	Instructions = map[string]func(*Interpreter, string){
		"INC":		INC,
		"RUN":		RUN,
		"RET":		RET,
		"AOUT":		AOUT,
		"OUT":      OUT,
		"IN":       IN,
		"STORE":    STORE,
		"LOAD":     LOAD,
		"ADD":      ADD,
		"SUBTRACT": SUBTRACT,
		"MULTIPLY": MULTIPLY,
		"DIVIDE":   DIVIDE,
		"JUMP":     JUMP,
		"JINEG":    JINEG,
		"JIZERO":   JIZERO,
		"HALT":     HALT,
	}
)

type Interpreter struct {
	Accum  int
	Data   map[string]int
	Code   []string
	LineP  int
	Labels map[string]int
	//Function specific jazz
	Wait *sync.WaitGroup
	Functions map[string]*Interpreter
	Parent *Interpreter
}

func main() {
	var wg sync.WaitGroup
	Interp := Interpreter{0, map[string]int{}, []string{""}, 0, map[string]int{}, &wg, map[string]*Interpreter{}, &Interpreter{}}
	Interp.LoadFile(os.Args[1])
	Interp.PopulateLabels()
	Interp.Run()
}

func (I *Interpreter) LoadFile(filename string) {
    byteInput, fileErr := ioutil.ReadFile(filename)
    if fileErr != nil {
        fmt.Println(fileErr)
        os.Exit(1)
    }

	input := strings.TrimRight(string(byteInput), "\r\n")
	splitCode := strings.Split(input, "\n")
	code := []string{}
	for i := 0; i < len(splitCode); i++ {
		begin := strings.Trim(splitCode[i], "\r\n")
		//append(code, strings.Trim(splitCode[i], "\r\n"))
		pos := strings.Index(begin, ";")
		if pos == 0 {
			continue
		} else if pos > 0 {
			begin := strings.Split(begin, ";")[0]
			code = append(code, begin)
		} else {
			code = append(code, begin)
		}
	}

	I.Code = code
}

func (I *Interpreter) PopulateLabels() {
	// Populating the label list
	for i := 0; i < len(I.Code); i++ {
		Line := strings.Split(I.Code[i], " ")
		if _, ok := Instructions[Line[0]]; !ok {
			I.Labels[Line[0]] = i
		}
	}
}

func (I *Interpreter) Run() {
	for I.LineP = 0; I.LineP < len(I.Code); I.LineP++ {
		Line := strings.Split(I.Code[I.LineP], " ")
		Command(I, Line)
	}
}

/*func AIN(I *Interpreter, _ string) {
    var s string
    fmt.Print(Green(">"), Reset())
    fmt.Scan(&s)

    input := strings.TrimRight(s, "\r\n")
    var intErr error
    I.Accum, intErr = strconv.Atoi(input)
    if intErr != nil {
        fmt.Println(intErr)
        os.Exit(1)
    }
}*/

func AOUT(I *Interpreter, _ string) {
	fmt.Println(Blue("<"), string(rune(I.Accum)), Reset())
}

func INC(I *Interpreter, location string) {
    var wg sync.WaitGroup
    Interp := Interpreter{0, map[string]int{}, []string{""}, 0, map[string]int{}, &wg, map[string]*Interpreter{}, &Interpreter{}}
	Interp.LoadFile(location)
	Interp.PopulateLabels()
	Interp.Parent = I
	I.Functions[location] = &Interp
}

func RUN(I *Interpreter, name string) {
	I.Functions[name].Accum = I.Accum
	I.Wait.Add(1)
	go I.Functions[name].Run()
	I.Wait.Wait()
}

func RET(I *Interpreter, _ string) {
	I.Parent.Accum = I.Accum
	I.Parent.Wait.Done()
}

func IN(I *Interpreter, _ string) {
	var s string
	fmt.Print(Green("> "), Reset())
	fmt.Scan(&s)

	input := strings.TrimRight(s, "\r\n")
	var intErr error
	I.Accum, intErr = strconv.Atoi(input)
	if intErr != nil {
		fmt.Println(intErr)
		os.Exit(1)
	}
}

func Command(I *Interpreter, Line []string) {
	if _, ok := Instructions[Line[0]]; !ok {
		if len(Line) == 2 {
			Instructions[Line[1]](I, "")
		} else {
			Instructions[Line[1]](I, Line[2])
		}
	} else {
		if len(Line) == 1 {
			Instructions[Line[0]](I, "")
		} else if len(Line) == 2 {
			Instructions[Line[0]](I, Line[1])
		}
	}
}

func OUT(I *Interpreter, _ string) {
	fmt.Println(Blue("<")+Reset(), I.Accum)
}

func STORE(I *Interpreter, address string) {
	I.Data[address] = I.Accum
}

func LOAD(I *Interpreter, address string) {
	integer, intErr := strconv.Atoi(address)
	if intErr != nil {
		I.Accum = I.Data[address]
	} else {
		I.Accum = integer
	}
}

func ADD(I *Interpreter, address string) {
	integer, intErr := strconv.Atoi(address)
	if intErr != nil {
		I.Accum += I.Data[address]
	} else {
		I.Accum += integer
	}
}

func SUBTRACT(I *Interpreter, address string) {
	integer, intErr := strconv.Atoi(address)
	if intErr != nil {
		I.Accum -= I.Data[address]
	} else {
		I.Accum -= integer
	}
}

func MULTIPLY(I *Interpreter, address string) {
	integer, intErr := strconv.Atoi(address)
	if intErr != nil {
		I.Accum = I.Accum * I.Data[address]
	} else {
		I.Accum = I.Accum * integer
	}
}

func DIVIDE(I *Interpreter, address string) {
	integer, intErr := strconv.Atoi(address)
	if intErr != nil {
		I.Accum = I.Accum / I.Data[address]
	} else {
		I.Accum = I.Accum / integer
	}
}

func JUMP(I *Interpreter, label string) {
	I.LineP = I.Labels[label] - 1
}

func JINEG(I *Interpreter, label string) {
	if I.Accum < 0 {
		I.LineP = I.Labels[label] - 1
	}
}

func JIZERO(I *Interpreter, label string) {
	if I.Accum == 0 {
		I.LineP = I.Labels[label] - 1
	}
}

func HALT(I *Interpreter, _ string) {
	os.Exit(0)
}
