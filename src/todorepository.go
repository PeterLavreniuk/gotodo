package gotodo

import (
	"fmt"
)

// ToDoItem represents a note with title, description and unique id
type ToDoItem struct {
	// The title of the note
	Title string `json:"title"`
	// The description of the note
	Description string `json:"description"`
	// The unique id of the note
	Id int64 `json:"id"`
}

type ToDoRepository struct {
	BaseRepository
}

func CreateTodoRepository(config *Config) (*ToDoRepository, error) {
	baseRepository, err := CreateBaseRepository(config)
	if err != nil {
		return nil, err
	}

	repository := ToDoRepository{
		BaseRepository: *baseRepository,
	}

	return &repository, nil
}

func (r *ToDoRepository) Create(item *ToDoItem) (*ToDoItem, error) {
	query := fmt.Sprintf("insert into gotododb.todoitem values (null, '%s', '%s')", item.Title, item.Description)
	res, err := r.DB.Exec(query)

	defer r.DB.Close()

	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()

	todoItem := ToDoItem{
		Title:       item.Title,
		Description: item.Description,
		Id:          id}

	return &todoItem, err
}

func (r *ToDoRepository) Update(item *ToDoItem) (*ToDoItem, error) {
	query := fmt.Sprintf("update gotododb.todoitem set title = '%s', description = '%s' where id = %d", item.Title, item.Description, item.Id)
	_, err := r.DB.Exec(query)

	defer r.DB.Close()

	return item, err
}

func (r *ToDoRepository) Delete(id int) (bool, error) {
	query := fmt.Sprintf("delete from gotododb.todoitem where id = %d", id)
	res, err := r.DB.Exec(query)

	defer r.DB.Close()

	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ToDoRepository) All() ([]*ToDoItem, error) {
	temp := make([]*ToDoItem, 0)
	query := "select id, title, description from todoitem"

	rows, err := r.DB.Query(query)

	defer r.DB.Close()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		todoitem := new(ToDoItem)
		err := rows.Scan(&todoitem.Id, &todoitem.Title, &todoitem.Description)
		if err != nil {
			return nil, err
		}

		temp = append(temp, todoitem)
	}

	return temp, nil
}

func (r *ToDoRepository) Get(id int) (*ToDoItem, error) {
	row := r.DB.QueryRow("select id, title, description from todoitem where id = ?", id)

	defer r.DB.Close()

	todoitem := new(ToDoItem)
	err := row.Scan(&todoitem.Id, &todoitem.Title, &todoitem.Description)

	if err != nil {
		return nil, err
	}

	return todoitem, nil
}
