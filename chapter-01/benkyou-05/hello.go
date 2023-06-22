package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Salary struct {
	Basic float64
}

type Employee struct {
	FirstName, LastName, Email string
	Age                        int
	MonthlySalary              []Salary
}

func main() {
	jsonFile, _ := os.Open("my_salary.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var employee Employee
	_ = json.Unmarshal(byteValue, &employee)
	//fmt.Printf("%+v", employee)
	json, _ := json.MarshalIndent(employee, "", " ")

	fmt.Println(string(json))
}
