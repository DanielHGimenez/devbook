package repositories

import (
	"api/src/models"
	"database/sql"
)

type Post struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *Post {
	return &Post{db}
}

func (repository *Post) Create(post models.Post) (uint64, error) {
	statement, err := repository.db.Prepare("insert into posts (title, content, author_id) values (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (repository *Post) FindRelatedPosts(userID uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		select
			distinct(p.id),
			p.title,
			p.content,
			p.likes,
			p.author_id,
			u.nick,
			p.created_at

		from posts p

		left join followers f
		on p.author_id = f.followed_id

		inner join users u
		on p.author_id = u.id

		where p.author_id = ?
		or f.follower_id = ?

		order by p.id
		desc`,
		userID,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.AuthorID,
			&post.AuthorNick,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *Post) FindOne(postID uint64) (*models.Post, error) {
	rows, err := repository.db.Query(`
		select
			p.id,
			p.title,
			p.content,
			p.likes,
			p.author_id,
			u.nick,
			p.created_at

		from posts p

		left join users u
		on p.author_id = u.id

		where p.id = ?
		limit 1`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post models.Post
	if rows.Next() {
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.AuthorID,
			&post.AuthorNick,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &post, nil
}

func (repository *Post) Update(postID uint64, post models.Post) error {
	statement, err := repository.db.Prepare("update posts set title = ?, content = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(post.Title, post.Content, postID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Post) Delete(postID uint64) error {
	statement, err := repository.db.Prepare("delete from posts where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(postID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Post) Exists(postID uint64) (bool, error) {
	rows, err := repository.db.Query("select true from posts where id = ? limit 1", postID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (repository *Post) IsAuthor(authorID uint64, postID uint64) (bool, error) {
	rows, err := repository.db.Query("select true from posts where id = ? and author_id = ? limit 1", postID, authorID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (repository *Post) WasLiked(userID uint64, postID uint64) (bool, error) {
	rows, err := repository.db.Query("select true from likes where user_id = ? and post_id = ?", userID, postID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (repository *Post) Like(userID uint64, postID uint64) error {
	_, err := repository.db.Exec("insert into likes (user_id, post_id) values (?, ?)", userID, postID)
	if err != nil {
		return err
	}

	transaction, err := repository.db.Begin()
	if err != nil {
		return err
	}

	_, err = transaction.Exec("select * from posts where id = ? for update", postID)
	if err != nil {
		return err
	}

	_, err = transaction.Exec("update posts set likes = likes + 1 where id = ?", postID)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *Post) UnLike(userID uint64, postID uint64) error {
	_, err := repository.db.Exec("delete from likes where user_id = ? and post_id = ?", userID, postID)
	if err != nil {
		return err
	}

	transaction, err := repository.db.Begin()
	if err != nil {
		return err
	}

	_, err = transaction.Exec("select * from posts where id = ? for update", postID)
	if err != nil {
		return err
	}

	_, err = transaction.Exec("update posts set likes = likes - 1 where id = ?", postID)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}
