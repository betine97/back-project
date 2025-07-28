package view

import (
	"github.com/betine97/back-project.git/src/controller/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
)

func ConvertDomainToResponse(resp *entity.User) dtos.NewUser {
	return dtos.NewUser{
		ID:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		City:      resp.City,
	}
}
