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

		blockedUsers, err := ur.GetBlockedUsersByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		friendNames := fileterBlockedFriends(friends, blockedUsers)

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func fileterBlockedFriends(friends []model.User, blocked []model.User) []string {
	friendNames := make([]string, 0)

	for _, friend := range friends {
		if !contains(blocked, friend) {
			friendNames = append(friendNames, friend.Name)
		}
	}

	return friendNames
}

func contains(slice []model.User, value model.User) bool {
	for _, item := range slice {
		if item.ID == value.ID {
			return true
		}
	}
	return false
}
