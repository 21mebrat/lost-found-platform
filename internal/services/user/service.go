package user

import "github.com/21mebrat/lost-found-platform/internal/repository/user"

type Service struct {
	repo user.Repository
}

var _ UserService = (*Service)(nil)

func NewService(rep user.Repository) *Service {
	return &Service{
		repo: rep,
	}
}
