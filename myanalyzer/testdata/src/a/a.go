package a

import "fmt"

func testFindPointerOfLoopVar() {
	/*
		{
			for foo := 0 ; foo < 3 ; foo++ {
				bar := &foo
				f(foo) // OK
				g(&foo) // NG

				f(*bar) // OK
				g(bar) // NG

				bar = nil
				g(bar) // OK
			}
		}
	*/
	{
		foo := 0
		for foo := foo; foo < 3; foo++ { // want "foo found"
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			foo := foo
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			foo := &foo
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			fmt.Println(&foo) // want "unary expr found"
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			a := 1
			fmt.Println(&a)
		}
		for foo := &foo; *foo < 3; *foo++ { // want "foo found"
			fmt.Println(*foo) // TODO: これも通る
		}
		for foo := foo; foo < 3; foo++ { // want "foo found"
			foo := foo
			fmt.Println(&foo)
		}
	}
	{
		foo := 0
		for ; foo < 3; foo++ {
			fmt.Println(foo)
		}
		foo = 0
		for ; foo < 3; foo++ {
			fmt.Println(foo)
		}
		foo = 0
		for ; foo < 3; foo++ {
			foo := foo
			fmt.Println(foo)
		}
		foo = 0
		for ; foo < 3; foo++ {
			fmt.Println(&foo)
		}
		foo = 0
		for ; foo < 3; foo++ {
			foo := foo
			fmt.Println(&foo)
		}
	}
}

func test() {
	foo := 0
	for foo := foo; foo < 3; foo++ { // want "foo found"
		foo := &foo
		fmt.Println(foo)
	}
	for foo := foo; foo < 3; foo++ { // want "foo found"
		fmt.Println(&foo) // want "unary expr found"
	}
}

func appendTest1() {
	var m []*int
	foo := 0
	for foo := foo; foo < 3; foo++ { // want "foo found"
		m = append(m, &foo) // want "unary expr found"
	}
	fmt.Println(m)
}

func appendTest2() {
	var m []*int
	foo := 0
	for foo := foo; foo < 3; foo++ { // want "foo found"
		foo := foo
		m = append(m, &foo)
	}
	fmt.Println(m)
}
