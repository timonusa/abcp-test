package main

import (
	"fmt"
	"time"
)

//Struct for tasks
type TaskT struct {
	ID               int
	CreationTime     string // время создания
	ExecutingTime    string // время выполнения
	Result           bool
	ErrorMsg         []byte
}

//Channels to hold both successful and error tasks
var resultSeccess = make(chan map[int]TaskT)
var resultErrors = make(chan error)


func main() {

	//create channel for holding all tasks
	tasksChan := make(chan TaskT, 10)

	//go to create tasks
	go taskCreator(tasksChan)

	//go to work over the tasks
	go taskWorker(tasksChan)

	//show error about the task
	//better to write to error log
	go func(){
		println("Errors:")
		for {
			println(<- resultErrors)
		}
	}()

	//show successful
	//better to write to success log
	go func(){
		println("Done tasks:")
		for  {
			println(<- resultSeccess)
		}
	}()

	//to keep daemon mode
	for range time.Tick(time.Second * 1) {

	}

}



//function for creating tasks
func taskCreator(ch chan TaskT) {

	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond() % 2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		ch <- TaskT{CreationTime: ft, ID: int(time.Now().Nanosecond())} // передаем таск на выполнение
	}

}


//function for checking tasks
func taskWorker(ch chan TaskT)  {

	for {
		task := <-ch
		taskTime, _ := time.Parse(time.RFC3339, task.CreationTime)

		if taskTime.After(time.Now().Add(-20 * time.Second)) {
			task.Result = true
		} else {
			task.Result = false
			task.ErrorMsg = []byte("something went wrong")
		}
		task.ExecutingTime = time.Now().Format(time.RFC3339Nano)

		taskSorter(task)
	}

}

//function for sorting tasks
func taskSorter(task TaskT) {

	if task.Result{
		taskMap := map[int]TaskT{
			task.ID : task,
		}
		resultSeccess <- taskMap
	} else {
		resultErrors <-  fmt.Errorf("Task id %d time %s, error %s", task.ID, task.CreationTime, task.ErrorMsg)
	}

}
