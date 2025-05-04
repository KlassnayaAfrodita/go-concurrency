package main

import "testing"

//func NewObject(params []string) (Srayer, eror){
//	if len(params) == 1 {
//		return Base{name:params[0]}, nil
//	} else if len(params) == 2 {
//		return Child{Base{params[0]}, params[1]}, nil
//	} else {
//		return nil, errors.New("some error")
//	}
//
//}

func TestNewObject(t *testing.T) {
	firstTestSlice := []string{"name"}
	secondTestSlice := []string{"name", "lastname"}
	thirdTestSlice := []string{}

	firstResult, err := NewObject(firstTestSlice)
	if err != nil {
		t.Error()
	}
	secondResult, err := NewObject(secondTestSlice)
	if err != nil {
		t.Error()
	}

	thirdResult, err := NewObject(thirdTestSlice)
	if !(err != nil && thirdResult == nil) {
		t.Error()
	}

	//firstType := firstResult.(type)
	switch firstResult.(type) {
	case Child:
		t.Error()
	}

	//secondResult := secondResult.(type)
	switch secondResult.(type) {
	case Base:
		t.Error()
	}
}
