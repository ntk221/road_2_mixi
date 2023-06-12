package handler

import (
	"database/sql"
	"net/http"
	"problem1/repository"

	"github.com/labstack/echo"
)

func GetFriendListHandler(db *sql.DB) (echo.HandlerFunc, error) {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		ur := repository.NewUserRepository(db)

		friendIDs, err := ur.GetFriendIDs(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		blockedUsers, err := ur.GetBlockedUsers(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		filteredFriends := make([]string, 0)
		for _, friend := range friendIDs {
			if !contains(blockedUsers, friend) {
				filteredFriends = append(filteredFriends, friend)
			}
		}

		friendNames, err := ur.GetFriendNames(filteredFriends)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		c.JSON(http.StatusOK, friendNames)
		return nil
	}, nil
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
