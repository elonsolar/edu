package util_test

import (
	"edu/internal/util"
	"fmt"
	"testing"
)

type Menu struct {
	Id       int32
	Name     string
	ParentId int32
	Children []*Menu
}

func TestTree(t *testing.T) {

	menuList := []*Menu{{
		Id:       101,
		Name:     "101",
		ParentId: 10,
	},
		{
			Id:       10,
			Name:     "10",
			ParentId: 1,
		},
		{
			Id:       1,
			Name:     "1",
			ParentId: 0,
		},
	}

	t.Run("test buildTree", func(t *testing.T) {

		treeDatList := util.NewTreeBuilder[Menu, int32]().Build(menuList, 0)

		fmt.Println(treeDatList)
	})
}
