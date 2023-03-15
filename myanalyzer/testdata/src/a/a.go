package a

import "fmt"

func testFindPointerOfLoopVar() {
	{
		foo := 0
		for foo := foo; foo < 3; foo++ {
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ {
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ {
			foo := foo
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ {
			foo := &foo
			fmt.Println(foo)
		}
		for foo := foo; foo < 3; foo++ {
			fmt.Println(&foo) // want "foo is pointer"
		}
		for foo := foo; foo < 3; foo++ {
			a := 1
			fmt.Println(&a)
		}
		for foo := &foo; *foo < 3; *foo++ {
			fmt.Println(*foo)
		}
		for foo := foo; foo < 3; foo++ {
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

func appendTest1() {
	var m []*int
	foo := 0
	for foo := foo; foo < 3; foo++ {
		m = append(m, &foo) // want "foo is pointer"
	}
	fmt.Println(m)
}

func appendTest2() {
	var m []*int
	foo := 0
	for foo := foo; foo < 3; foo++ {
		foo := foo
		m = append(m, &foo)
	}
	fmt.Println(m)
}

func anonymousFuncTest() {
	foo := 0
	for foo := foo; foo < 3; foo++ {
		func() {
			fmt.Println(foo)
		}()
	}
	for foo := foo; foo < 3; foo++ {
		func() {
			fmt.Println(&foo) // TODO: want "foo is pointer" のエラーが必要
		}()
	}
	for foo := foo; foo < 3; foo++ {
		foo := foo
		func() {
			fmt.Println(&foo)
		}()
	}
}
