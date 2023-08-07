package store

import (
	"database/sql"
	"fmt"

	"github.com/MelloHey/htmxtodo/types"
	"github.com/MelloHey/htmxtodo/utils"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetTodos() ([]*types.Todo, error)
	CreateTodo(*types.Todo) error
	DeleteTodos(int)
	//GetItemsByTodoID(int) ([]*types.Item, error)
	//CreateItem(*types.Item) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	fmt.Println("YOYO")
	db_user := utils.ViperEnvVariable("DB_USER")
	db_name := utils.ViperEnvVariable("DB_NAME")
	db_password := utils.ViperEnvVariable("DB_PASSWORD")
	fmt.Printf("DB_PASSWORD: %s\n", db_name)
	connStr := "user=" + db_user + " dbname=" + db_name + " password=" + db_password + " sslmode=disable"
	fmt.Printf("DB_PASSWORD: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.createTodoTable(); err != nil {
		return err
	}
	if err := s.createItemTable(); err != nil {
		return err
	}

	if err := s.createTestTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) createTodoTable() error {
	query := `create table if not exists todo (
		id serial primary key,
		name varchar(100),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createTestTable() error {
	fmt.Println("TST")
	query := `create table if not exists test (
		id serial primary key,
		name varchar(100),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createTestTwoTable() error {
	fmt.Println("TST TWO")
	query := `create table if not exists testtwo (
		id serial primary key,
		name varchar(100),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateTodo(todo *types.Todo) error {
	query := `insert into todo 
	(name, created_at)
	values ($1, $2)`

	_, err := s.db.Query(
		query,
		todo.Name,
		todo.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetTodos() ([]*types.Todo, error) {
	rows, err := s.db.Query("select * from todo")
	if err != nil {
		return nil, err
	}

	todos := []*types.Todo{}
	for rows.Next() {
		account, err := scanIntoTodo(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, account)
	}

	return todos, nil
}

func (s *PostgresStore) DeleteTodos(id int) {

	_, err := s.db.Query("delete from todo where id = $1", id)

	if err == nil {
		fmt.Printf("ERROR DELETING RECORD with id %d\n", id)
	}

}

/* func (s *PostgresStore) GetItemsByTodoID(id int) ([]*types.Item, error) {
	rows, err := s.db.Query("select * from item where todo_id = $1", id)
	if err != nil {
		return nil, err
	}

	items := []*types.Item{}
	for rows.Next() {
		item, err := scanIntoItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
} */

/* func (s *PostgresStore) CreateItem(item *types.Item) error {
	query := `insert into item
	(todo_id, name, created_at, done)
	values ($1, $2, $3, $4)`

	_, err := s.db.Query(
		query,
		item.TodoID,
		item.Name,
		item.CreatedAt,
		item.Done,
	)

	if err != nil {
		return err
	}

	return nil
} */

func scanIntoTodo(rows *sql.Rows) (*types.Todo, error) {
	todo := new(types.Todo)
	err := rows.Scan(
		&todo.ID,
		&todo.Name,
		&todo.CreatedAt,
	)

	return todo, err
}

/* func scanIntoItem(rows *sql.Rows) (*types.Item, error) {
	item := new(types.Item)
	err := rows.Scan(
		&item.ID,
		&item.TodoID,
		&item.Name,
		&item.CreatedAt,
		&item.Done,
	)

	return item, err
} */

func (s *PostgresStore) createItemTable() error {
	query := `create table if not exists item (
		id serial primary key,
		todo_id INT,
		name VARCHAR(255) NOT NULL,
		created_at timestamp,
		done boolean,
		CONSTRAINT fk_todo
		   FOREIGN KEY(todo_id) 
		   REFERENCES todo(id)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createUserTable() error {
	query := `create table if not exists user (
		id serial primary key,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at timestamp,

	)`

	_, err := s.db.Exec(query)
	return err
}
