package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"problem1/domain"
	"problem1/infra"
	"problem1/usecases"
	"strconv"

	"github.com/labstack/echo"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetUserListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// ユーザー情報取得用の設定
		ur := infra.NewUserRepository()

		users, err := ur.GetUsers(h.db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get users"))
		}

		c.JSON(http.StatusOK, users)
		return nil
	}
}

func (h *Handler) GetUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		// ユーザー情報取得用の設定
		uc := usecases.NewUserService(h.db)

		user, err := uc.GetUserByID((domain.UserID(id)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get user"))
		}

		c.JSON(http.StatusOK, user)
		return nil
	}
}

func (h *Handler) GetFriendListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("id must be integer"))
			return nil
		}

		// ユーザー情報取得用の設定
		us := usecases.NewUserService(h.db)

		friends, err := us.GetFriendList(domain.UserID(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get friend list"))
		}

		friendCollection := domain.NewUserCollection(friends)
		friendCollection = friendCollection.GetUniqueUsers()
		friendNames := friendCollection.GetUserNames()

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func (h *Handler) GetFriendOfFriendListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("id must be integer"))
			return nil
		}
		hop, err := strconv.Atoi(c.QueryParam("hop"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("hop must be integer"))
			return nil
		}
		if hop < 1 || hop > 10 {
			c.JSON(http.StatusBadRequest, errors.New("hop must be between 1 and 10"))
			return nil
		}
		params, err := get_limit_page(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		us := usecases.NewUserService(h.db)

		friendList, err := us.GetFriendList(domain.UserID(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get friend list"))
		}

		friend_of_friends, err := us.GetFriendListFromUsers(friendList, hop)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get friend of friend list"))
		}

		friend_of_friends = pagenate(params, friend_of_friends)

		filteredNames, err := get_names(friend_of_friends)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("failed to get friend of friend list"))
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

func get_limit_page(c echo.Context) (PagenationParams, error) {
	limit_s := c.QueryParam("limit")
	if limit_s == "" {
		limit_s = "10"
	}

	limit, err := strconv.Atoi(limit_s)
	if err != nil {
		return PagenationParams{}, err
	}

	page_s := c.QueryParam("page")
	if page_s == "" {
		page_s = "0"
	}

	page, err := strconv.Atoi(page_s)
	if err != nil {
		return PagenationParams{}, err
	}

	// limit and page should be positive
	if limit < 0 || page < 0 {
		return PagenationParams{}, err
	}

	params := PagenationParams{
		Limit:  limit,
		Offset: page,
	}

	return params, nil
}

func pagenate(params PagenationParams, users []domain.User) []domain.User {
	if params.Offset > len(users) {
		return make([]domain.User, 0)
	}

	if params.Offset+params.Limit > len(users) {
		return users[params.Offset:]
	}

	return users[params.Offset : params.Offset+params.Limit]
}
