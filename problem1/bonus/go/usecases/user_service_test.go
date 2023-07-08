package usecases

import (
	"errors"
	"log"
	"problem1/domain"
	"problem1/testutils"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetRealFriends(t *testing.T) {
	tx, err := testutils.OpenDBForTest(t).Begin()
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	testUsers := []domain.User{
		{
			UserID: 1,
			Name:   "test1",
			FriendList: []int{
				2, 3,
			},
			BlockList: []int{
				4,
			},
		},
	}

	moq := &UserRepositoryMock{}
	moq.GetByIDFunc = func(userID int, tx domain.Queryer) (domain.User, error) {
		// 特定のユーザーIDに対して異なる結果を返す
		switch userID {
		case 1:
			return testUsers[0], nil
		case 2:
			return domain.User{UserID: 2, Name: "test2", BlockList: []int{1}}, nil
		case 3:
			return domain.User{UserID: 3, Name: "test3"}, nil
		default:
			return domain.User{}, errors.New("user not found")
		}
	}

	sut := NewUserService(tx, moq)

	ret, err := sut.GetFriendList(testUsers[0].UserID)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%v\n", ret)

}
