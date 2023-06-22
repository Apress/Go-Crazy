package main

import (
	"encoding/json"
	"io/ioutil"
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
	data := Employee{
		FirstName:     "Nicolas",
		LastName:      "Modrzyk",
		Email:         "hellonico at gmail.com",
		Age:           43,
		MonthlySalary: []Salary{{Basic: 15000.00}, {Basic: 16000.00}, {Basic: 17000.00}},
	}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("my_salary.json", file, 0644)
}
