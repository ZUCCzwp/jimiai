package adminRouterService

import (
	"errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"jiyu/model/adminRouterModel"
	"jiyu/repo/adminRouterRepo"
	"log"
)

func DelRouter(id int) error {
	return adminRouterRepo.Delete(id)
}

func SaveRouter(data adminRouterModel.RouterRequest) error {
	if data.ParentId > 0 && !adminRouterRepo.ExistById(data.ParentId) {
		return errors.New("父路由不存在")
	}

	router := adminRouterModel.Router{
		Model:      gorm.Model{ID: uint(data.Id)},
		RouterBase: data.RouterBase,
	}

	return adminRouterRepo.Save(router)
}

func ListRouter(role int) ([]adminRouterModel.RouterResponse, error) {
	routers, err := adminRouterRepo.List()
	if err != nil {
		log.Println("adminRouterService.ListRouter error: ", err)
		return nil, err
	}

	parents := lo.Filter(routers, func(item adminRouterModel.Router, index int) bool {
		return item.ParentId == 0
	})

	var result []adminRouterModel.RouterResponse
	for _, parent := range parents {
		if role != 1 && !lo.Contains(parent.Meta.Roles, role) {
			continue
		}

		response := parent.ConvertToResponse()

		response.Children = lo.FilterMap(routers, func(item adminRouterModel.Router, index int) (adminRouterModel.RouterResponse, bool) {
			return item.ConvertToResponse(), item.ParentId == int(parent.ID)
		})

		result = append(result, response)
	}

	return result, nil
}
