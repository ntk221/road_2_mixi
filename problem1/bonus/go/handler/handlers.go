package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"problem1/domain"
	"problem1/infra"
	"problem1/types"
	"problem1/usecases"
	"strconv"

	"github.com/labstack/echo"
)

func GetUserListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ユーザー情報取得用の設定
		ur := infra.NewUserRepository()

		users, err := ur.GetUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, users)
		return nil
	}
}

func GetUserHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, errors.New("id is empty"))
			return nil
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		// ユーザー情報取得用の設定
		ur := infra.NewUserRepository()
		txAdmin := usecases.NewTxAdmin(db)
		us := usecases.NewUserService(txAdmin, ur)

		user, err := us.GetUserByID((domain.UserID(idInt)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, user)
		return nil
	}
}

func GetFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, errors.New("id is empty"))
			return nil
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		// ユーザー情報取得用の設定
		ur := infra.NewUserRepository()
		txAdmin := usecases.NewTxAdmin(db)
		us := usecases.NewUserService(txAdmin, ur)

		friends, err := us.GetFriendList(domain.UserID(idInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		friendNames := make([]string, 0)
		for _, friend := range friends {
			friendNames = append(friendNames, friend.Name)
		}

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func GetFriendOfFriendListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, errors.New("id is empty"))
			return nil
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		ur := infra.NewUserRepository()
		txAdmin := usecases.NewTxAdmin(db)
		us := usecases.NewUserService(txAdmin, ur)

		friendList, err := us.GetFriendList(domain.UserID(idInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		fof, err := us.GetFriendListFromUsers(friendList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		filteredNames, err := get_names(fof)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, filteredNames)
		return nil
	}
}

func GetFriendOfFriendListPagingHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// parameterを取得
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, errors.New("id is empty"))
			return nil
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}
		params, err := get_limit_page(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		ur := infra.NewUserRepository()
		txAdmin := usecases.NewTxAdmin(db)
		us := usecases.NewUserService(txAdmin, ur)

		friendList, err := us.GetFriendList(domain.UserID(idInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		fof, err := us.GetFriendListFromUsers(friendList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		fof = pagenate(params, fof)

		filteredNames, err := get_names(fof)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		c.JSON(http.StatusOK, filteredNames)
		return nil
	}
}

func get_names(fof []domain.User) ([]string, error) {
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
	return filteredNames, nil
}

/*func get_id(c echo.Context) (int, error) {
	id_s := c.QueryParam("id")
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return 0, err
	}
	return id, nil
}*/

func get_limit_page(c echo.Context) (types.PagenationParams, error) {
	limit_s := c.QueryParam("limit")
	if limit_s == "" {
		limit_s = "10"
	}

	limit, err := strconv.Atoi(limit_s)
	if err != nil {
		return types.PagenationParams{}, err
	}

	page_s := c.QueryParam("page")
	if page_s == "" {
		page_s = "0"
	}

	page, err := strconv.Atoi(page_s)
	if err != nil {
		return types.PagenationParams{}, err
	}

	// limit and page should be positive
	if limit < 0 || page < 0 {
		return types.PagenationParams{}, err
	}

	params := types.PagenationParams{
		Limit:  limit,
		Offset: page,
	}

	return params, nil
}

func pagenate(params types.PagenationParams, users []domain.User) []domain.User {
	if params.Offset > len(users) {
		return make([]domain.User, 0)
	}

	if params.Offset+params.Limit > len(users) {
		return users[params.Offset:]
	}

	return users[params.Offset : params.Offset+params.Limit]
}
