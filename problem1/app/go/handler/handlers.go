package handler

import (
	"database/sql"
	"net/http"
	"problem1/model"
	"problem1/repository"

	"github.com/labstack/echo"
)

func GetFriendListHandler(db *sql.DB) (echo.HandlerFunc, error) {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		ur := repository.NewUserRepository(db)

		friends, err := ur.GetFriendsByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		blockedUsers, err := ur.GetBlockedUsersByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		filteredFriends := make([]model.User, 0)
		for _, friend := range friends {
			if !contains(blockedUsers, friend) {
				filteredFriends = append(filteredFriends, friend)
			}
		}

		c.JSON(http.StatusOK, filteredFriends)
		return nil
	}, nil
}

func contains(slice []model.User, value model.User) bool {
	for _, item := range slice {
		if item.ID == value.ID {
			return true
		}
	}
	return false
}
