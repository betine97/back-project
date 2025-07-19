package view

import (
	"github.com/betine97/back-project.git/src/controller/dtos"
	entity "github.com/betine97/back-project.git/src/model/entitys"
)

func ConvertDomainToResponse(resp *entity.CreateUser) dtos.NewUser {
	return dtos.NewUser{
		First_Name: resp.First_Name,
		Email:      resp.Email,
	}
}
