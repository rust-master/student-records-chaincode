package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Student struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Degree string  `json:"degree"`
	GPA    float32 `json:"gpa"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) AddStudent(ctx contractapi.TransactionContextInterface, id string, name string, degree string, gpa float32) error {
	exists, err := s.StudentExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("student with ID %s already exists", id)
	}

	student := Student{
		ID:     id,
		Name:   name,
		Degree: degree,
		GPA:    gpa,
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, studentJSON)
}

func (s *SmartContract) UpdateGPA(ctx contractapi.TransactionContextInterface, id string, gpa float32) error {
	studentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if studentJSON == nil {
		return fmt.Errorf("student with ID %s does not exist", id)
	}

	var student Student
	err = json.Unmarshal(studentJSON, &student)
	if err != nil {
		return err
	}

	student.GPA = gpa

	updatedStudentJSON, err := json.Marshal(student)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedStudentJSON)
}

func (s *SmartContract) QueryStudent(ctx contractapi.TransactionContextInterface, id string) (*Student, error) {
	studentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if studentJSON == nil {
		return nil, fmt.Errorf("student with ID %s does not exist", id)
	}

	var student Student
	err = json.Unmarshal(studentJSON, &student)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (s *SmartContract) StudentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	studentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return studentJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating Student Records chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Student Records chaincode: %v", err)
	}
}
