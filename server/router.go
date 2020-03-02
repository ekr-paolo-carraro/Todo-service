package server

import (
	"net/http"
	"strconv"

	"github.com/ekr-paolo-carraro/todoTest/Todo-service/model"
	"github.com/gin-gonic/gin"
)

type HandlerRouter struct {
	persistDelegate model.TodoDelegate
}

func NewRouter(delegate model.TodoDelegate) (*gin.Engine, error) {

	hr := HandlerRouter{delegate}

	err := delegate.InitData()
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	router.Static("/client", "client")

	v1 := router.Group("api/v1/todos")
	{
		v1.GET("/", hr.getAllHandler)
		v1.GET("/:index", hr.getTodoByIDHandler)
		v1.PUT("/", hr.insertTodo)
		v1.POST("/", hr.updateTodo)
	}

	return router, nil
}

func (hr HandlerRouter) getAllHandler(c *gin.Context) {
	todos, err := hr.persistDelegate.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, todos)
}

func (hr HandlerRouter) getTodoByIDHandler(c *gin.Context) {
	originalIndex := c.Param("index")
	index, err := strconv.Atoi(c.Param("index"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if index == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + originalIndex + " not allowed"})
	} else {
		todo, err := hr.persistDelegate.GetTodo(index)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, todo)
		}
	}
}

func (hr HandlerRouter) insertTodo(c *gin.Context) {
	var todoItem model.TodoItem
	c.BindJSON(&todoItem)

	result, err := hr.persistDelegate.InsertTodo(todoItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"sid": result})
	}
}

func (hr HandlerRouter) updateTodo(c *gin.Context) {
	var todoItem model.TodoItem
	c.BindJSON(&todoItem)
	result, err := hr.persistDelegate.UpdateTodo(todoItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"sid": result})
	}
}
