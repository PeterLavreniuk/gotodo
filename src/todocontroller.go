package gotodo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToDoController struct {
	Config Config
}

// @Summary Deletes a note by its id
// @ID delete-todoitem
// @Produce json
// @Param id path int true "note id"
// @Success 200 {object} ToDoItem
// @Failure 404 {object} NotFoundResponse
// @Failure 500 {object} ErrorResponse
// @Router /{id} [delete]
func (c *ToDoController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	repository, err := CreateTodoRepository(&c.Config)
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := repository.Delete(parsedId)
	if err != nil {
		ctx.Error(err)
		return
	}

	if !res {
		ctx.IndentedJSON(http.StatusNotFound, NotFoundResponse{Message: "todo not found", Id: parsedId})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, EntityDeleted{Message: "todo item deleted", Id: parsedId})
}

// @Summary GetAll returns all existing notes
// @ID get-all-todoitems
// @Produce json
// @Success 200 {array} ToDoItem
// @Failure 500 {object} ErrorResponse
// @Router / [get]
func (c *ToDoController) GetAll(ctx *gin.Context) {
	repository, err := CreateTodoRepository(&c.Config)
	if err != nil {
		ctx.Error(err)
		return
	}

	items, err := repository.All()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, items)
}

// @Summary Get returns a note by its id
// @ID get-todoitem
// @Produce json
// @Param id path int true "note id"
// @Success 200 {object} ToDoItem
// @Failure 500 {object} ErrorResponse
// @Router /{id} [get]
func (c *ToDoController) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	repository, err := CreateTodoRepository(&c.Config)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := repository.Get(parsedId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, item)
}

// @Summary Create a new note and returns it
// @ID create-todoitem
// @Produce json
// @Param data body ToDoItem true "note"
// @Success 200 {object} ToDoItem
// @Failure 500 {object} ErrorResponse
// @Router / [post]
func (c *ToDoController) Create(ctx *gin.Context) {
	var todoItemRequest ToDoItemRequest

	err := ctx.BindJSON(&todoItemRequest)
	if err != nil {
		return
	}

	repository, err := CreateTodoRepository(&c.Config)
	if err != nil {
		ctx.Error(err)
		return
	}

	todoItem := ToDoItem{Title: todoItemRequest.Title, Description: todoItemRequest.Description}
	createdToDoItem, err := repository.Create(&todoItem)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdToDoItem)
}

// @Summary Updates a note by its id
// @ID update-todoitem
// @Produce json
// @Param data body ToDoItemRequest true "note"
// @Param id path int true "note id"
// @Success 200 {object} ToDoItem
// @Failure 500 {object} ErrorResponse
// @Router /{id} [put]
func (c *ToDoController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedId, _ := strconv.Atoi(id)

	var todoItemRequest ToDoItemRequest

	err := ctx.BindJSON(&todoItemRequest)
	if err != nil {
		return
	}

	repository, err := CreateTodoRepository(&c.Config)
	if err != nil {
		ctx.Error(err)
		return
	}

	todoItem, err := repository.Get(parsedId)
	if err != nil {
		ctx.Error(err)
		return
	}

	todoItem.Title = todoItemRequest.Title
	todoItem.Description = todoItemRequest.Description

	updatedToDoItem, err := repository.Update(todoItem)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedToDoItem)
}

// ToDoItem represents a blog post with a title, description, body, author, and publication status
type ToDoItemRequest struct {
	// The author of the blog
	Title string `json:"title"`
	// The author of the blog
	Description string `json:"description"`
}

// ErrorResponse represents an error response from the api
type ErrorResponse struct {
	// The error message
	Message string `json:"message"`
}

// NotFoundResponse represents an error response from the api when entity is not found by id
type NotFoundResponse struct {
	// The error message
	Message string `json:"message"`
	// Id of the item not found
	Id int `json:"id"`
}

// EntityDeleted represents an response when item has been deleted by id
type EntityDeleted struct {
	// The additional message
	Message string `json:"message"`
	// Id of the item not found
	Id int `json:"id"`
}
