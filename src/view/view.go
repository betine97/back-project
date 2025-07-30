package view

import (
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
)

func ConvertDomainToResponse(resp *entity.User) dtos.NewUser {
	return dtos.NewUser{
		ID:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
	}
}
