package app

import (
	activityQueries "github.com/geoikonomou/forum/internal/app/activities/queries"
	categoryCommands "github.com/geoikonomou/forum/internal/app/categories/commands"
	categoryQueries "github.com/geoikonomou/forum/internal/app/categories/queries"
	commentCommands "github.com/geoikonomou/forum/internal/app/comments/commands"
	commentQueries "github.com/geoikonomou/forum/internal/app/comments/queries"
	oauthservice "github.com/geoikonomou/forum/internal/app/oauth"
	topicCommands "github.com/geoikonomou/forum/internal/app/topics/commands"
	topicQueries "github.com/geoikonomou/forum/internal/app/topics/queries"
	userCommands "github.com/geoikonomou/forum/internal/app/user/commands"
	userQueries "github.com/geoikonomou/forum/internal/app/user/queries"
	votecommands "github.com/geoikonomou/forum/internal/app/votes/commands"
	voteQueries "github.com/geoikonomou/forum/internal/app/votes/queries"
	"github.com/geoikonomou/forum/internal/domain/activity"
	"github.com/geoikonomou/forum/internal/domain/category"
	"github.com/geoikonomou/forum/internal/domain/comment"
	"github.com/geoikonomou/forum/internal/domain/oauth"
	"github.com/geoikonomou/forum/internal/domain/topic"
	"github.com/geoikonomou/forum/internal/domain/user"
	"github.com/geoikonomou/forum/internal/domain/vote"
	"github.com/geoikonomou/forum/internal/pkg/bcrypt"
	"github.com/geoikonomou/forum/internal/pkg/uuid"
)

type Queries struct {
	UserLoginGithub    oauthservice.OAuthService
	GetTopic           topicQueries.GetTopicRequestHandler
	GetAllTopics       topicQueries.GetAllTopicsRequestHandler
	GetComment         commentQueries.GetCommentRequestHandler
	GetCommentsByTopic commentQueries.GetCommentsByTopicRequestHandler
	UserLoginEmail     userQueries.UserLoginEmailRequestHandler
	UserLoginUsername  userQueries.UserLoginUsernameRequestHandler
	GetCategoryByID    categoryQueries.GetCategoryByIDHandler
	GetAllCategories   categoryQueries.GetAllCategoriesRequestHandler
	GetCounts          voteQueries.GetCountsRequestHandler
	GetUserActivity    activityQueries.GetUserActivityHandler
}

type Commands struct {
	UserRegister   userCommands.UserRegisterRequestHandler
	CreateTopic    topicCommands.CreateTopicRequestHandler
	UpdateTopic    topicCommands.UpdateTopicRequestHandler
	DeleteTopic    topicCommands.DeleteTopicRequestHandler
	CreateComment  commentCommands.CreateCommentRequestHandler
	UpdateComment  commentCommands.UpdateCommentRequestHandler
	DeleteComment  commentCommands.DeleteCommentRequestHandler
	CreateCategory categoryCommands.CreateCategoryRequestHandler
	UpdateCategory categoryCommands.UpdateCategoryRequestHandler
	DeleteCategory categoryCommands.DeleteCategoryRequestHandler
	CastVote       votecommands.CastVoteRequestHandler
	DeleteVote     votecommands.DeleteVoteRequestHandler
}

type UserServices struct {
	Queries  Queries
	Commands Commands
}

type Services struct {
	UserServices UserServices
}

func NewServices(userRepo user.Repository, categoryRepo category.Repository, topicRepo topic.Repository, commentRepo comment.Repository, voteRepo vote.Repository, oauthRepo oauth.Repository, activityRepo activity.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				*oauthservice.NewOAuthService(oauthRepo, uuidProvider),
				topicQueries.NewGetTopicHandler(topicRepo, commentRepo),
				topicQueries.NewGetAllTopicsHandler(topicRepo, categoryRepo),
				commentQueries.NewGetCommentHandler(commentRepo),
				commentQueries.NewGetCommentsByTopicRequestHandler(commentRepo),
				userQueries.NewUserLoginEmailHandler(userRepo, encryption),
				userQueries.NewUserLoginUsernameHandler(userRepo, encryption),
				categoryQueries.NewGetCategoryByIDHandler(categoryRepo),
				categoryQueries.NewGetAllCategoriesHandler(categoryRepo),
				voteQueries.NewGetCountsRequestHandler(voteRepo),
				activityQueries.NewGetUserActivityHandler(activityRepo),
			},
			Commands: Commands{
				userCommands.NewUserRegisterHandler(userRepo, uuidProvider, encryption),
				topicCommands.NewCreateTopicHandler(topicRepo),
				topicCommands.NewUpdateTopicHandler(topicRepo),
				topicCommands.NewDeleteTopicHandler(topicRepo),
				commentCommands.NewCreateCommentRequestHandler(commentRepo),
				commentCommands.NewUpdateCommentRequestHandler(commentRepo),
				commentCommands.NewDeleteCommentHandler(commentRepo),
				categoryCommands.NewCreateCategoryHandler(categoryRepo),
				categoryCommands.NewUpdateCategoryHandler(categoryRepo),
				categoryCommands.NewDeleteCategoryHandler(categoryRepo),
				votecommands.NewCastVoteHandler(voteRepo),
				votecommands.NewDeleteVoteHandler(voteRepo),
			},
		},
	}
}
