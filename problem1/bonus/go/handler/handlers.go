package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"problem1/domain/entity"
	"problem1/domain/valueObject"
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
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		c.JSON(http.StatusOK, users)
		return nil
	}
}

func (h *Handler) GetUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}
		if id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "idは1以上の正の整数を入力してください")
		}

		// ユーザー情報取得用の設定
		uc := usecases.NewUserService(h.db)

		userID := valueObject.NewUserID(id)
		user, err := uc.GetUserByID(userID)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		c.JSON(http.StatusOK, user)
		return nil
	}
}

func (h *Handler) GetFriendListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}
		if id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "idは1以上の正の整数を入力してください")
		}

		// ユーザー情報取得用の設定
		us := usecases.NewUserService(h.db)

		userID := valueObject.NewUserID(id)
		friends, err := us.GetFriendList(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		friendNames := friends.GetUserNames()

		c.JSON(http.StatusOK, friendNames)
		return nil
	}
}

func (h *Handler) GetFriendOfFriendListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}
		hopStr := c.QueryParam("hop")
		if hopStr == "" {
			hopStr = "1"
		}
		hop, err := strconv.Atoi(hopStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}
		if hop < 1 || hop > 10 {
			return echo.NewHTTPError(http.StatusBadRequest, "hopは1より大きく，10未満の整数を入力してください")
		}
		params, err := get_limit_page(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "TODO")
		}

		us := usecases.NewUserService(h.db)

		userID := valueObject.NewUserID(id)
		friendList, err := us.GetFriendList(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		friend_of_friends, err := us.GetFriendListFromUsers(friendList, hop)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		friend_of_friends = pagenate(params, friend_of_friends)

		filteredNames := friend_of_friends.GetUserNames()
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, "サーバーで問題が発生しました")
		}

		c.JSON(http.StatusOK, filteredNames)
		return nil
	}
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
		return PagenationParams{}, errors.New("limitとpageはそれぞれ正の整数で無くてはならない")
	}

	params := PagenationParams{
		Limit:  limit,
		Offset: page,
	}

	return params, nil
}

func pagenate(params PagenationParams, UserCollection *entity.UserCollection) *entity.UserCollection {
	if params.Offset > len(UserCollection.Users) {
		return UserCollection
	}

	if params.Offset+params.Limit > len(UserCollection.Users) {
		offset := UserCollection.Users[params.Offset:]
		return entity.NewUserCollection(offset)
	}

	divided := UserCollection.Users[params.Offset : params.Offset+params.Limit]
	return entity.NewUserCollection(divided)
}
