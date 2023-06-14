package handler

import (
	"database/sql"
	"net/http"
	"problem1/model"
	"problem1/repository"

	"github.com/labstack/echo/v4"
)

func GetFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		ur := repository.NewUserRepository(db)

		friends, err := ur.GetFriendsByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		/*if len(friends) == 0 {
			panic("No friends")
		}*/

		blockedUsers, err := ur.GetBlockedUsersByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		/*if len(blockedUsers) == 0 {
			panic("No blocked users")
		}*/

		friendNames := make([]string, 0)
		for _, friend := range friends {
			if !contains(blockedUsers, friend) {
				friendNames = append(friendNames, friend.Name)
			}
		}

		/*if len(friendNames) == 0 {
			panic("No friend names")
		}*/

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func contains(slice []model.User, value model.User) bool {
	for _, item := range slice {
		if item.ID == value.ID {
			return true
		}
	}
	return false
}
