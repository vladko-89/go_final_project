package tasks

import (
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	Database *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{Database: db}
}

func (repo *TaskRepository) Create(task *Task) (id int64, err error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, ?, ?, ?);
	`
	result, err := repo.Database.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения задачи: %v", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID задачи: %v", err)
	}

	return id, nil
}

func (repo *TaskRepository) GetList() (*sql.Rows, error) {

	query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT 50;
		`
	rows, err := repo.Database.Query(query)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка задач: %v", err)
	}

	return rows, nil

}

func (repo *TaskRepository) Get(id string) *sql.Row {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	return repo.Database.QueryRow(query, id)
}

func (repo *TaskRepository) Delete(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	_, err := repo.Database.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepository) Update(task *TaskRow) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?;
	`

	result, err := repo.Database.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача с id=%s не найдена", task.ID)
	}
	return nil
}

func (repo *TaskRepository) UpdateDate(id, date string) error {

	_, err := repo.Database.Exec(`
		UPDATE scheduler
		SET date = ?
		WHERE id = ?;
	`, date, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}

	return nil

}
