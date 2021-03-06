### 使用的Go语言的web框架gin-gonic构建RESTful API服务

使用`Go`语言构建一个`todo`应用的API，使用的是简单快速的框架`gin-goni`框架，后端使用`mongodb`,需要安装的三方库

```
go get gopkg.in/gin-gonic/gin.v1
go get github.com/globalsign/mgo
```

创建下面的API接口  

* POST      todos/
* GET       todos/ 
* GET       todos/{id}
* PUT       todos/{id}
* DELETE    todos/{id}

#### 初始化路由

```
package main

import (
	"github.com/gin-gonic/gin"
)

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

```

关于 `createTodo`等方法会在后面一一介绍，先创建一下`ToDo`相关的model  

```
const (
	db         = "ToDo"
	collection = "ToDoList"
	host       = "127.0.0.1:27017"
)

var globalS *mgo.Session

type (
	todoModel struct {
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
	globalS, err := mgo.Dial(host)
	if err != nil {
		log.Fatal(err)
	}
	
}
```

#### createTodo 路由方法

```
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
```

从`gin`的上下文获取post的数据，并保存数据到数据库，如果成功，返回对应的id

#### fetchAllTodo 路由方法

```
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
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed, CreatedAt: time.Now()})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}
```

#### fetchSingleTodo 路由方法

```
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	id := c.Param("id")
	ms := globalS.Copy()
	mc := ms.DB(db).C(collection)
	mc.FindId(id)
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
```

#### updateTodo 路由方法

```
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

	mc.UpdateId(id, todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo update successfully!"})
}
```

#### deleteTodo 路由方法

```
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
```

#### 运行测试

运行  

`go run main.go`  

使用 `postman`测试提交数据  

* createTodo
  ![createTodo](http://wx4.sinaimg.cn/large/6d6f88d1ly1fvo44szr3sj218g0bj0td.jpg)
* fetchSingleTodo
  ![fetchSingleTodo](http://wx4.sinaimg.cn/large/6d6f88d1ly1fvo44vpatij218i0aiwf4.jpg)
* deleteTodo
  ![deleteTodo](http://wx2.sinaimg.cn/large/6d6f88d1ly1fvo44y83m2j218i08fgm2.jpg)

[更多精彩内容](http://coderminer.com)  
[译自](https://medium.com/@thedevsaddam/build-restful-api-service-in-golang-using-gin-gonic-framework-85b1a6e176f3)    
