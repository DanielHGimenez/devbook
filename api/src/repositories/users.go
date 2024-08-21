package repositories

import (
	"api/src/models"
	"database/sql"
	"strings"
)

type User struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *User {
	return &User{db}
}

func (repository *User) Create(userModel models.User) (uint64, error) {
	statement, err := repository.db.Prepare("insert into users (name, nick, email, password) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(userModel.Name, userModel.Nick, userModel.Email, userModel.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (repository *User) FindBy(name string, nick string, email string) ([]models.User, error) {
	var query strings.Builder
	query.WriteString("select id, name, nick, email, password, created_at from users where true = true")

	var params []any
	if name != "" {
		query.WriteString(" and name like ?")
		params = append(params, name)
	}
	if nick != "" {
		query.WriteString(" and nick like ?")
		params = append(params, nick)
	}
	if email != "" {
		query.WriteString(" and email like ?")
		params = append(params, email)
	}

	rows, err := repository.db.Query(query.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *User) FindOne(id uint64) (*models.User, error) {
	rows, err := repository.db.Query(
		"select id, name, nick, email, password, created_at from users where id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, nil
}

func (repository User) Update(id uint64, user *models.User) error {
	statement, err := repository.db.Prepare("update users set name = ?, nick = ?, email = ?, password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Nick, user.Email, user.Password, id)
	return err
}

func (repository *User) Delete(id uint64) error {
	statement, err := repository.db.Prepare("delete from users where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	return err
}

func (repository *User) FindAllFollowers(followedID uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
		select id, name, nick, email, created_at from users u
		inner join followers f
		on u.id = f.follower_id
		where f.user_id = ?
	`, followedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *User) FindAllFollowing(followerID uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
		select id, name, nick, email, created_at from users u
		inner join followers f
		on u.id = f.user_id
		where f.follower_id = ?
	`, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *User) FindPassword(userID uint64) (string, error) {
	rows, err := repository.db.Query("select password from users where id = ?", userID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var password string
	if rows.Next() {
		if err = rows.Scan(&password); err != nil {
			return "", err
		}
	}

	return password, nil
}

func (repository *User) SavePassword(userID uint64, password string) error {
	statement, err := repository.db.Prepare("update users set password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(password, userID)
	return err
}
