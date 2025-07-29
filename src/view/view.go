package view

import (
	dtos_controllers "github.com/betine97/back-project.git/src/controller/dtos_controllers"
	entity "github.com/betine97/back-project.git/src/model/entitys"
)

func ConvertDomainToResponse(resp *entity.User) dtos_controllers.NewUser {
	return dtos_controllers.NewUser{
		ID:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		City:      resp.City,
	}
}
