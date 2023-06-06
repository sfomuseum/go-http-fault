package main

import (
	"fmt"
	"encoding/json"
	"os"
	"reflect"
	"log"
	
	"github.com/sfomuseum/go-http-fault/v2"
)

type Bar struct {
	Hello string
	Status int
	Error error
}
	
type Foo struct {
	Hello string
}

// https://groups.google.com/g/golang-nuts/c/dV7Yw78wWzU

func Merge(a interface{}, b interface{})( d interface{}) {
	aType := reflect.TypeOf(a)
	if aType.Kind() != reflect.Struct {
		panic("a is not a struct")
	}

	bType := reflect.TypeOf(b)

	log.Println(bType.Kind())
	
	if bType.Kind() != reflect.Struct {
		panic("b is not a struct")
	}

	var fields []reflect.StructField
	for i:=0 ; i< aType.NumField(); i++{
	  fields =append(fields, aType.Field(i))
	}
	for i:=0 ; i< bType.NumField(); i++{
	  fields =append(fields, bType.Field(i))
	}

	dType := reflect.StructOf(fields)
	dVal := reflect.Indirect(reflect.New(dType))


	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	

	for i := 0; i < aType.NumField(); i++ {
		dVal.FieldByName(aType.Field(i).Name).Set(aVal.Field(i))
	}
	for i := 0; i < bType.NumField(); i++ {
		dVal.FieldByName(bType.Field(i).Name).Set(bVal.Field(i))
	}
	
	d = dVal.Interface()
		
	return 
}

func main() {
	
	vars := fault.FaultHandlerVars{
		Status: 999,
		Error:  fmt.Errorf("SAD"),
	}
	
	foo := &Foo{
		Hello: "world",
	}

	v := Merge(vars, *foo)

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(v)

	log.Println(v.Status)
}
