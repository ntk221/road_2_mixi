package main

import (
	"database/sql"
	"net/http"
	"problem1/configs"
	"problem1/dao"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, World!")
	})

	e.GET("/get_friend_list", func(c echo.Context) error {
		id := c.QueryParam("id")
		UserRepository := dao.NewUser(db)
		friends, err := UserRepository.GetFriends(id)
		if err != nil {
			return err
		} else if friends == nil {
			return c.JSON(http.StatusOK, "no friends")
		} else {
			var ret []string
			for _, friend := range friends {
				ret = append(ret, friend.Name)
			}
			return c.JSON(http.StatusOK, ret)
		}
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		id := c.QueryParam("id")
		UserRepository := dao.NewUser(db)
		user, err := UserRepository.GetByID(id)
		if err != nil {
			return err
		}
		friends, err := UserRepository.GetFriends(id)
		if err != nil {
			return err
		} else if friends == nil {
			return c.JSON(http.StatusOK, "no friends")
		} else {
			for _, friend := range friends {
				user.AddFriend(friend)
			}
			var ret []string
			for _, friend := range friends {
				fof, err := UserRepository.GetFriends(friend.Id)
				if err != nil {
					return err
				} else if fof == nil {
					continue
				} else if fof != nil {
					for _, f := range fof {
						if f.Id == user.Id {
							continue
						} else if user.IsFriend(f) {
							continue
						}
						ret = append(ret, f.Name)
					}
				}
			}
			return c.JSON(http.StatusOK, ret)
		}
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
