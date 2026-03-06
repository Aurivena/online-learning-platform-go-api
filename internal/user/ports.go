package user

import "online-learning-platform-go-api/internal/user/dto"

type Postgres interface {
	Create(req dto.RegistrationRequest) error
	Get(id int) (dto.AccountResponse, error)
	Update(req dto.UpdateRequest, id int) error
}
