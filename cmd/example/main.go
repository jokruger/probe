package main

import (
	"time"

	"github.com/jokruger/probe"
)

func Foo1() {
	defer probe.Start("Foo1").Stop()
	time.Sleep(10 * time.Millisecond)
}

func Foo2() {
	defer probe.Start("Foo2").Stop()
	time.Sleep(20 * time.Millisecond)
}

func main() {
	Foo1()
	Foo2()

	p1 := probe.Start("Block1")
	time.Sleep(30 * time.Millisecond)
	p1.Stop()

	p2 := probe.Start("Block2")
	time.Sleep(40 * time.Millisecond)
	p2.Stop()

	Foo2()

	probe.PrintReport()
}
