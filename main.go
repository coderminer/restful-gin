package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"

	"github.com/gin-gonic/gin"
)

const (
	db         = "ToDo"
	collection = "ToDoList"
	host       = "127.0.0.1:27017"
)

var globalS *mgo.Session

type (
	todoModel struct {
		ID        string    `bson:"_id" json:"id"`
		Title     string    `bson:"title" json:"title"`
		Completed int       `bson:"completed" json:"completed"`
		CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	}

	transformedTodo struct {
		ID        string    `bson:"_id" json:"id"`
		Title     string    `bson:"title" json:"title"`
		Completed bool      `bson:"completed" json:"completed"`
		CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	}
)

func init() {
	ms, err := mgo.Dial(host)
	if err != nil {
		log.Fatal(err)
	}
	globalS = ms
}

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{
		ID:        bson.NewObjectId().Hex(),
		Title:     c.PostForm("title"),
		Completed: completed,
		CreatedAt: time.Now(),
	}

	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	defer ms.Close()
	mc.Insert(todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated,
		"message": "Todo item created successfully!", "resourceId": todo.ID})
}

func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	defer ms.Close()
	mc.Find(nil).All(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!!"})
		return
	}
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed, CreatedAt: item.CreatedAt})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	id := c.Param("id")
	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	mc.FindId(id).One(&todo)
	if todo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed, CreatedAt: todo.CreatedAt}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

func updateTodo(c *gin.Context) {
	var todo todoModel
	id := c.Param("id")
	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	defer ms.Close()
	mc.FindId(id).One(&todo)
	if todo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	todo.Title = c.PostForm("title")
	todo.Completed, _ = strconv.Atoi(c.PostForm("completed"))

	err := mc.UpdateId(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "update error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo update successfully!"})
}

func deleteTodo(c *gin.Context) {
	var todo todoModel
	id := c.Param("id")
	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	defer ms.Close()
	mc.FindId(id).One(&todo)
	if todo.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	mc.RemoveId(id)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run()
}
