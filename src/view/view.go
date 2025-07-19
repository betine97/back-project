package view

import (
	"back-project/src/controller/dtos"
	entity "back-project/src/model/entitys.go"
)

func ConvertDomainToResponse(resp *entity.CreateUser) dtos.NewUser {
	return dtos.NewUser{
		First_Name: resp.First_Name,
		Email:      resp.Email,
	}
}
