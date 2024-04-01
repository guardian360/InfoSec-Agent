package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/scan"
)

type DataBase struct {
}

func NewDataBase() *DataBase {
	return &DataBase{}
}

func (d *DataBase) GetAllSeverities(checks []checks.Check, resultIDs []int) ([]int, error) {
	return scan.GetAllSeverities(checks, resultIDs)
}
