package handler

import (
	"database/sql"
	"net/http"
	"problem1/repository"
	service "problem1/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := get_id(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}
		ur := repository.NewUserRepository()
		us := service.NewUserService(db, ur)

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
		id, err := get_id(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}
		ur := repository.NewUserRepository()
		us := service.NewUserService(db, ur)

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

func get_id(c echo.Context) (int, error) {
	id_s := c.QueryParam("id")
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func get_limit_page(c echo.Context) (service.PagenationParams, error) {
	limit_s := c.QueryParam("limit")
	limit, err := strconv.Atoi(limit_s)
	if err != nil {
		return service.PagenationParams{}, err
	}

	page_s := c.QueryParam("page")
	page, err := strconv.Atoi(page_s)
	if err != nil {
		return service.PagenationParams{}, err
	}

	params := service.PagenationParams{
		Limit:  limit,
		Offset: page,
	}

	return params, nil
}

func GetFriendOfFriendListPagingHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := get_id(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}
		params, err := get_limit_page(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		ur := repository.NewUserRepository()
		us := service.NewUserService(db, ur)

		friendList, err := us.GetFriendList(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		// get friend of friend list pagenated
		fof, err := us.GetFriendListFromUsersWithPagenation(friendList, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
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
