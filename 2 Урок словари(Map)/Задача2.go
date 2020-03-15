package main

import "fmt"
import s "strings"

func isIsomorphic(c string, t string) (answer bool) {
	var num int  = 0 
	map2:= map[string] int{
	"q":1, "w":1, "e":1, "r":1, "t":1, "y":1, "u":1, "i":1, "o":1, "p":1,
	"a":1, "s":1, "d":1, "f":1, "g":1, "h":1, "j":1, "k":1, "l":1, "z":1,
	"x":1, "c":1, "v":1, "b":1, "n":1, "m":1,
	}
	
	arr:= []string{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
	"a", "s", "d", "f", "g", "h", "j", "k", "l", "z",
	"x", "c", "v", "b", "n", "m"}
	
	el1 := len(c)
	el2 := len(t)
	
	if el1 == el2 {
	
	for i:= 0; i<len(arr); i++{
	
	number := s.Count(t, arr[i])
	
	if number >= 2{
	num = map2[arr[i]]
	}
	
	}
	for i:= 0; i<len(arr); i++{
	
	number := s.Count(c, arr[i])
	
	if number >= 2{
	num = num + map2[arr[i]]
	}
	
	}
	}else{
	answer = false
	}
	
	if num >= 2{
	answer = true
	}

	return answer

}

func title (){
	fmt.Println("Программа на 2 задание, которая сравнивает 2 строки")
	fmt.Println("Она проверяет 2 строки изоморфные они или нет")
	fmt.Println("false строки не изоморфные, true изоморфные\n")
	
}

func end (){
	fmt.Println("\nОтвет получен")
	fmt.Println("Программа заканчивает свою работу")
}


func main() {
	var (c string = "foo"
		t string = "baa")
	title()	
	fmt.Println("Ответ программы\n")
	fmt.Println(isIsomorphic(c, t))
	end()
}
