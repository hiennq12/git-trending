package cache

import (
	"time"
)

type RepoDescription struct {
	ID          int64     `db:"id"`
	RepoName    string    `db:"repo_name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
