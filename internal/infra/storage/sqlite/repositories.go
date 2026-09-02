package sqlite

import (
	"database/sql"

	"github.com/geoikonomou/forum/internal/domain/activity"
	"github.com/geoikonomou/forum/internal/domain/category"
	"github.com/geoikonomou/forum/internal/domain/comment"
	"github.com/geoikonomou/forum/internal/domain/notification"
	"github.com/geoikonomou/forum/internal/domain/oauth"
	"github.com/geoikonomou/forum/internal/domain/topic"
	"github.com/geoikonomou/forum/internal/domain/user"
	"github.com/geoikonomou/forum/internal/domain/vote"
	activities "github.com/geoikonomou/forum/internal/infra/storage/sqlite/activity"
	"github.com/geoikonomou/forum/internal/infra/storage/sqlite/categories"
	"github.com/geoikonomou/forum/internal/infra/storage/sqlite/comments"
	oauthrepo "github.com/geoikonomou/forum/internal/infra/storage/sqlite/oauth"
	"github.com/geoikonomou/forum/internal/infra/storage/sqlite/topics"
	"github.com/geoikonomou/forum/internal/infra/storage/sqlite/users"
	"github.com/geoikonomou/forum/internal/infra/storage/sqlite/votes"
)

type Repositories struct {
	UserRepo         user.Repository
	CategoryRepo     category.Repository
	TopicRepo        topic.Repository
	CommentRepo      comment.Repository
	VoteRepo         vote.Repository
	NotificationRepo notification.Repository
	OauthRepo        oauth.Repository
	ActivityRepo     activity.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     users.NewRepo(db),
		CategoryRepo: categories.NewRepo(db),
		TopicRepo:    topics.NewRepo(db),
		CommentRepo:  comments.NewRepo(db),
		VoteRepo:     votes.NewRepo(db),
		OauthRepo:    oauthrepo.NewOAuthRepository(db),
		ActivityRepo: activities.NewRepo(db),
	}
}
