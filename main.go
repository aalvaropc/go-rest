package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []Task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "tarea uno",
		Content: "terminar de setear tindev",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Insert a valid task", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newTask)
	if err != nil {
		http.Error(w, "Invalid task format", http.StatusBadRequest)
		return
	}

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome que")
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if taskID == task.ID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if taskID == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "The task with ID %v has been removed successfully", taskID)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func updateTask(w http.ResponseWriter, r *http.Request) {

	log.Println("updateTask called")

	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Please enter valid data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &updatedTask)
	if err != nil {
		http.Error(w, "Invalid task format", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if taskID == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
	})

	router.StrictSlash(true)

	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/", indexRoute)
	log.Fatal(http.ListenAndServe("localhost:3000", router))
}
