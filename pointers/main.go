package main

import "fmt"

type object struct {
	Id   string
	Name string
}

func main() {
	//Basic variable pointer stuff...
	obj1 := object{
		Id:   "1",
		Name: "Object 1",
	}

	obj2 := object{
		Id:   "2",
		Name: "Object 2",
	}

	Rename(obj1, "New Object 1")
	RenameByPointer(&obj2, "New Object 2")

	fmt.Println(obj1.Name)
	fmt.Println(obj2.Name)

	//Slice pointer stuff...

	slice1 := []object{obj1}
	slice2 := []*object{&obj2}

	RenameAllElements(slice1, "Slice Name - Object 1")
	RenameAllElementsByPointer(slice2, "Slice Name By Pointer - Object 2")

	for _, obj := range slice1 {
		fmt.Println(obj.Name)
	}
	for _, obj := range slice2 {
		fmt.Println(obj.Name)
	}

}

// Basic Variable Pointer Functions...
func Rename(obj object, name string) {
	obj.Name = name
}

func RenameByPointer(obj *object, name string) {
	obj.Name = name
}

// Slice pointer functions...
// These both change names correclty - find out why (slices must pass by pointer by default?)
func RenameAllElements(s []object, name string) {
	for i := range s {
		s[i].Name = name
	}
}

func RenameAllElementsByPointer(s []*object, name string) {
	for i := range s {
		s[i].Name = name
	}
}
