package handler

import (
	"database/sql"
	"net/http"
	"problem1/repository"
	service "problem1/services"

	"github.com/labstack/echo/v4"
)

func GetFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		ur := repository.NewUserRepository(db)
		us := service.NewUserService(ur)

		filteredFriends, err := us.GetFriendList(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		friendNames := make([]string, 0)
		for _, friend := range filteredFriends {
			friendNames = append(friendNames, friend.Name)
		}

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func GetFriendOfFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		ur := repository.NewUserRepository(db)
		us := service.NewUserService(ur)

		friendList, err := us.GetFriendList(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		fof, err := us.GetFriendListFromUsers(friendList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		fofNames := make([]string, 0)
		for _, v := range fof {
			fofNames = append(fofNames, v.Name)
		}

		uniqueNames := make(map[string]bool)
		filteredNames := make([]string, 0)

		for _, name := range fofNames {
			if !uniqueNames[name] {
				uniqueNames[name] = true
				filteredNames = append(filteredNames, name)
			}
		}

		c.JSON(http.StatusOK, filteredNames)
		return nil
	}
}
