package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func OnlyComments() {
	// Also signel comments are ignored
	// foo := false
	// if true {
	//
	//	foo = true
	//
	// }
}

func NoBlock(r string) (x uintptr)

func CommentOnFnLine() { // This is just a comment
	fmt.Println("hello")
}

func ReturnOnSameLine() error {
	s := func() error { return nil }
	_ = s
}

func LabeledBlocks() {
	goto end
end:
}

func If() {
	v, err := strconv.Atoi("1")
	if err != nil {
		fmt.Println(v)
	}

	a := true
	if a {
		fmt.Println("I'm OK")
	}

	a, b := true, false
	if !a {
		fmt.Println("I'm OK")
	}

	a, b = true, false
	if err != nil { // want "if statements should only be cuddled with assignments used in the if statement itself"
		fmt.Println("I'm not OK")
	}

	if true {
		fmt.Println("I'm OK")
	}
	if false { // want "if statements should only be cuddled with assignments"
		fmt.Println("I'm OK")
	}

	c := true
	fooBar := "string"
	a, b := true, false
	if a && b || !c && len(fooBar) > 2 { // want "only one cuddle assignment allowed before if statement"
		return false
	}
	a, b, c := someFunc()          // want "assignments should only be cuddled with other assignments"
	if z && len(y) > 0 || 3 == 4 { // want "only one cuddle assignment allowed before if statement"
		return true
	}

	a, b := true, false
	if c && b || !c && len(fooBar) > 2 {
		return false
	}

	r := map[string]interface{}{}
	if err := json.NewDecoder(someReader).Decode(&r); err != nil {
		return "this should be OK"
	}
}

func Return() {
	if true {
		if false {
			//
		}
		return // want "return statements should not be cuddled if block has more than two lines"
	}

	if false {
		fmt.Println("this is less than two lines")
		return
	}
	return // want "return statements should not be cuddled if block has more than two lines"
}

func AssignmentAndDeclaration() {
	if true {
		fmt.Println("I'm not OK")
	}
	foo := true // want "assignments should only be cuddled with other assignments"

	bar := foo
	baz := true
	biz := false

	var a bool
	fmt.Println(a) // want "expressions should not be cuddled with declarations or returns"

	b := true
	fmt.Println(b)
}

func Range() {
	anotherList := make([]string, 5)

	myList := make([]int, 10)
	for i := range anotherList { // want "ranges should only be cuddled with assignments used in the iteration"
		fmt.Println(i)
	}

	myList = make([]string, 5)

	anotherList = make([]int, 10)
	for i := range anotherList {
		fmt.Println(i)
	}

	someList, anotherList := GetTwoListsOfStrings()
	for i := range append(someList, anotherList...) {
		fmt.Println(i)
	}

	aThirdList := GetList()
	for i := range append(someList, aThirdList...) {
		fmt.Println(i)
	}
}

func FirstInBlock() {
	idx := i
	if i > 0 {
		idx = i - 1
	}

	vals := map[int]struct{}{}
	for i := range make([]int, 5) {
		vals[i] = struct{}{}
	}

	x := []int{}

	vals := map[int]struct{}{}
	for i := range make([]int, 5) { // want "ranges should only be cuddled with assignments used in the iteration"
		x = append(x, i)
	}
}

func OnlyCuddleOneAssignment() {
	foo := true
	bar := false

	biz := true || false
	if biz {
		return true
	}

	foo := true
	bar := false
	biz := true || false

	if biz {
		return true
	}

	foo := true
	bar := false
	biz := true || false
	if biz { // want "only one cuddle assignment allowed before if statement"
		return false
	}
}

func IdentifiersWithIndices() {
	runes := []rune{'+', '-'}
	if runes[0] == '+' || runes[0] == '-' {
		return string(runes[1:])
	}

	listTwo := []string{}

	listOne := []string{"one"}
	if listOne[0] == "two" {
		return "not allowed"
	}

	for i := range make([]int, 10) {
		key := GetKey()
		if val, ok := someMap[key]; ok {
			fmt.Println("ok!")
		}

		someOtherMap := GetMap()
		if val, ok := someOtherMap[key]; ok {
			fmt.Println("ok")
		}

		someIndex := 3
		if val := someSlice[someIndex]; val != nil {
			retunr
		}
	}
}

func Defer() {
	thingOne := getOne()
	thingTwo := getTwo()

	defer thingOne.Close()
	defer thingTwo.Close()

	thingOne := getOne()
	defer thingOne.Close()

	thingTwo := getTwo()
	defer thingTwo.Close()

	thingOne := getOne()
	defer thingOne.Close()
	thingTwo := getTwo()   // want "assignments should only be cuddled with other assignments"
	defer thingTwo.Close() // want "only one cuddle assignment allowed before defer statement"

	thingOne := getOne()
	thingTwo := getTwo()
	defer thingOne.Close() // want "only one cuddle assignment allowed before defer statement"
	defer thingTwo.Close()

	m := sync.Mutex{}

	m.Lock()
	defer m.Unlock()

	foo := true
	defer func(b bool) { // want "defer statements should only be cuddled with expressions on same variable"
		fmt.Printf("%v", b)
	}()
}

func For() {
	bool := true
	for { // want "for statement without condition should never be cuddled"
		fmt.Println("should not be allowed")

		if bool {
			break
		}
	}
}

func Go() {
	go func() {
		panic("is this real life?")
	}()

	fooFunc := func() {}
	go fooFunc()

	barFunc := func() {}
	go fooFunc() // want "go statements can only invoke functions assigned on line above"

	go func() {
		fmt.Println("hey")
	}()

	cuddled := true
	go func() { // want "go statements can only invoke functions assigned on line above"
		fmt.Println("hey")
	}()

	argToGo := 1
	go takesArg(argToGo)

	notArgToGo := 1
	go takesArg(argToGo) // want "go statements can only invoke functions assigned on line above"

	t1 := NewT()
	t2 := NewT()
	t3 := NewT()

	go t1()
	go t2()
	go t3()

	multiCuddle1 := NewT()
	multiCuddle2 := NewT()
	go multiCuddle2() // want "only one cuddle assignment allowed before go"

	t4 := NewT()
	t5 := NewT()
	go t5() // want "only one cuddle assignment allowed before go"
	go t4()
}

func Switch() {
	var b bool
	switch b {
	case true:
		return "a"
	case false:
		return "b"
	}

	t := time.Now()

	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	var b bool
	switch anotherBool { // want "switch statements should only be cuddled with variables switched"
	case true:
		return "a"
	case false:
		return "b"
	}

	t := time.Now()
	switch { // want "anonymous switch statements should never be cuddled"
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}
}

func TypeSwitch() {
	x := GetSome()
	switch v := x.(type) {
	case int:
		return "got int"
	default:
		return "got other"
	}

	var id string
	switch i := objectID.(type) {
	case int:
		id = strconv.Itoa(i)
	case uint32:
		id = strconv.Itoa(int(i))
	case string:
		id = i
	}

	var b bool
	switch AnotherVal.(type) { // want "type switch statements should only be cuddled with variables switched"
	case int:
		return "a"
	case string:
		return "b"
	}
}

func Append() {
	someList := []string{}

	bar := "baz"
	someList = append(someList, bar)

	bar := "baz"
	someList = append(someList, "notBar") // want "append only allowed to cuddle with appended value"

	bar := "baz"
	someList = append(someList, fmt.Sprintf("%s", bar))

	bar := "baz"
	whatever := appendFn(bar)

	s := []string{}

	s = append(s, "one")
	s = append(s, "two")
	s = append(s, "three")

	p.defs = append(p.defs, x)
	def.parseFrom(p)
	p.defs = append(p.defs, def)
	def.parseFrom(p)
	def.parseFrom(p)
	p.defs = append(p.defs, x)
}

func ExpressionAssignment() {
	foo := true
	someFunc(foo)

	foo := true
	someFunc(false) // want "only cuddled expressions if assigning variable or using from line above"
}

func Channel() {
	timeoutCh := time.After(timeout)

	for range make([]int, 10) {
		select {
		case <-timeoutCh:
			return true
		case <-time.After(10 * time.Millisecond):
			return false
		}
	}

	select {
	case <-time.After(1 * time.Second):
		return "1s"
	default:
		return "are we there yet?"
	}
}

func Branch() {
	for {
		if true {
			singleLine := true
			break
		}

		if true && false {
			multiLine := true
			maybeVar := "var"
			continue // want "branch statements should not be cuddled if block has more than two lines"
		}

		if false {
			multiLine := true
			maybeVar := "var"

			break
		}
	}
}

func Locks() {
	hashFileCache.Lock()
	out, ok := hashFileCache.m[file]
	hashFileCache.Unlock()

	mu := &sync.Mutex{}
	mu.X(y).Z.RLock()
	x, y := someMap[someKey]
	mu.RUnlock()
}

func SpliceAndSlice() {
	start := 0
	if v := aSlice[start:3]; v {
		fmt.Println("")
	}

	someKey := "str"
	if v, ok := someMap[obtain(someKey)+"str"]; ok {
		fmt.Println("Hey there")
	}

	end := 10
	if v := arr[3:notEnd]; !v { // want "if statements should only be cuddled with assignments used in the if statement itself"
		// Error
	}

	notKey := "str"
	if v, ok := someMap[someKey]; ok { // want "if statements should only be cuddled with assignments used in the if statement itself"
		// Error
	}
}

func Block() {
	var foo string

	x := func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}()

	x := func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}

	func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}()

	func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}

	func() {
		func() {
			return func() {
				var err error
				foo = "1" // want "assignments should only be cuddled with other assignments"
			}
		}()
	}

	var x error
	foo, err := func() { return "", nil } // want "assignments should only be cuddled with other assignments"

	defer func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}()

	go func() {
		var err error
		foo = "1" // want "assignments should only be cuddled with other assignments"
	}()
}

func CommentInLast() {
	switch nextStatement.(type) {
	case *someType, *anotherType:
	default:
		// Foo
		return
	}
}

func OnlyCheckTwoLinesAboveIfAssignment() {
	if err != nil {
		return
	}
	a, err := someFn() // want "assignments should only be cuddled with other assignments"
	if err != nil {    // want "only one cuddle assignment allowed before if statement"
		return result, err
	}
	b := someFn()                 // want "assignments should only be cuddled with other assignments"
	if err := fn(b); err != nil { // want "only one cuddle assignment allowed before if statement"
		return
	}

	if true {
		return
	}
	c, err := someFn() // want "assignments should only be cuddled with other assignments"
	if err != nil {    // want "only one cuddle assignment allowed before if statement"
		return result, err
	}
}

func IncrementDecrement() {
	a := 1
	a++
	a--
	b := 2

	if true {
		b--
	}
	b++ // want "assignments should only be cuddled with other assignments"

	go func() {}()
	b++            // want "assignments should only be cuddled with other assignments"
	go func() {}() // want "only one cuddle assignment allowed before go statement"
}

// Issue #09
func AnonymousFunc() {
	fmt.Println(func() string {
		_ = 1
		_ = 2
		_ = 3
		return "string" // want "return statements should not be cuddled if block has more than two lines"
	})

	fmt.Println(func() error {
		foo := "foo"
		fmt.Println("fmt.Println") // want "only cuddled expressions if assigning variable or using from line above"
		if foo == "bar" {          // want "if statements should only be cuddled with assignments"
			return fmt.Errorf("bar")
		}
		return nil // want "return statements should not be cuddled if block has more than two lines"
	}())
}
