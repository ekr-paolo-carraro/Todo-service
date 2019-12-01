package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ekr-paolo-carraro/todoTest/todo-app/model"
	_ "github.com/lib/pq"
)

type PersistenceDelegate struct {
	Db *sql.DB
}

func NewTodoDelegate() (*model.TodoDelegate, error) {
	//connection := "postgresql://serveruser:bFZLUKE8RQ86CTh@35.242.213.160:5432"
	connection := "host=localhost port=5432 user=srvuser password=ekr dbname=tododb sslmode=disable"
	db, err := sql.Open("postgres", connection)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	var tododelegate model.TodoDelegate
	tododelegate = PersistenceDelegate{db}
	return &tododelegate, nil
}

func (todoDelegate PersistenceDelegate) InitData() error {
	script, err := ioutil.ReadFile("../sql/createTodoTable.sql")
	if err != nil {
		return fmt.Errorf("error on reading script sql: %s", err)
	}

	_, err = todoDelegate.Db.Exec(string(script))
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "duplicate key value violates unique constraint") == false {
			return fmt.Errorf("error on exec script sql: %s", err)
		}
	}

	result, err := todoDelegate.GetAllItems()
	if result == nil || len(result) == 0 {
		script, err = ioutil.ReadFile("../sql/populateTable.sql")
		if err != nil {
			return fmt.Errorf("error on reading script sql: %s", err)
		}

		_, err = todoDelegate.Db.Exec(string(script))
		if err != nil {
			return fmt.Errorf("error on exec script sql: %s", err)
		}
	}

	return nil
}

func (todoDelegate PersistenceDelegate) GetAllItems() ([]model.TodoItem, error) {
	rows, err := todoDelegate.Db.Query("SELECT * FROM public.todos;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseResultSet(rows)
}

func (todoDelegate PersistenceDelegate) GetTodo(index int) (*model.TodoItem, error) {
	rows, err := todoDelegate.Db.Query("SELECT * FROM public.todos WHERE id=$1;", index)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, err
	}

	result, err := parseResultSet(rows)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (todoDelegate PersistenceDelegate) InsertTodo(item model.TodoItem) (string, error) {
	script, err := ioutil.ReadFile("../sql/insert.sql")
	if err != nil {
		return "", fmt.Errorf("error on reading script sql: %s", err)
	}

	statement, err := todoDelegate.Db.Prepare(string(script))
	if err != nil {
		return "", fmt.Errorf("error on preare insert statement : %s", err)
	}
	result := statement.QueryRow(item.UserID, item.Title, item.Executed)

	var tempId int64
	if err := result.Scan(&tempId); err != nil {
		return "", fmt.Errorf("can't parse id after insert: %s", err)
	}

	return strconv.FormatInt(tempId, 10), nil
}

func (todoDelegate PersistenceDelegate) UpdateTodo(item model.TodoItem) (string, error) {
	script, err := ioutil.ReadFile("../sql/update.sql")
	if err != nil {
		return "", fmt.Errorf("error on reading script sql: %s", err)
	}
	statement, err := todoDelegate.Db.Prepare(string(script))
	if err != nil {
		return "", fmt.Errorf("error on prepare update statement sql: %s", err)
	}
	_, err = statement.Exec(item.ID, item.UserID, item.Title, item.Executed)
	if err != nil {
		return "", fmt.Errorf("error on update item: %s", err)
	}
	return item.ID, nil
}

func parseResultSet(todoResult *sql.Rows) ([]model.TodoItem, error) {
	if todoResult == nil {
		return nil, errors.New("no result")
	}
	ts := make([]model.TodoItem, 0)
	for todoResult.Next() {
		var tempTodoItem model.TodoItem
		if err := todoResult.Scan(&tempTodoItem.ID,
			&tempTodoItem.UserID,
			&tempTodoItem.Title,
			&tempTodoItem.Executed); err != nil {
			return nil, err
		}
		ts = append(ts, tempTodoItem)
	}

	if len(ts) == 0 {
		return nil, errors.New("no result")
	}

	return ts, nil
}
