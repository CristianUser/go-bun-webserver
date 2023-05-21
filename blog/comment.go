package blog

import (
	"time"

	"github.com/cristianuser/go-bun-webserver/org"
	"github.com/uptrace/bun"
)

type Comment struct {
	bun.BaseModel `bun:"comments,alias:c"`

	ID   uint64 `json:"id"`
	Body string `json:"body"`

	Author   *org.Profile `json:"author" bun:"rel:belongs-to"`
	AuthorID uint64       `json:"-"`

	ArticleID uint64 `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
