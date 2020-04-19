package main

import (
	"context"
	"fmt"

	"github.com/zzzhr1990/go-mongo-util/helper"
	"github.com/zzzhr1990/go-mongo-util/model"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	test := &model.TestModel{
		Identity:  "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		TestParam: "BBBB",
	}
	fmt.Println("Hello")
	_, err := helper.UpdateMany(context.Background(), nil, bson.D{}, test, map[string]bool{"test_param_2": false, "empty_string": true})
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}
	fmt.Println("============================================================")

}
