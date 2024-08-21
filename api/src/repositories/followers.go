package repositories

import (
	"database/sql"
)

type Follower struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *Follower {
	return &Follower{db}
}

func (repository Follower) Exists(followedID uint64, followerID uint64) (bool, error) {
	rows, err := repository.db.Query("select * from followers where followed_id = ? and follower_id = ?", followedID, followerID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (repository Follower) Create(followedID uint64, followerID uint64) error {
	statement, err := repository.db.Prepare("insert into followers (followed_id, follower_id) values (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(followedID, followerID)
	return err
}

func (repository Follower) Delete(followedID uint64, followerID uint64) error {
	statement, err := repository.db.Prepare("delete from followers where followed_id = ? and follower_id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(followedID, followerID)
	return err
}
