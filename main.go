package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/trewanek/repository-with-tx/infrastructure/persistence/rdb"
	"github.com/trewanek/repository-with-tx/interface/repository"
	"github.com/trewanek/repository-with-tx/model"
	"github.com/trewanek/repository-with-tx/presenter"
)

var userRepo repository.IUserRepository

func main() {
	ctx := context.Background()
	dbConn, err := rdb.NewDBConn()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	userRepo = rdb.NewUserMySqlRepository(dbConn)

	err = dbConn.Transact(ctx, func() error {
		users, err := userRepo.FindAll(ctx)
		if err != nil {
			return err
		}

		found, err := userRepo.Find(ctx, "2")
		if err != nil {
			return err
		}

		newUser := new(model.User)
		newUser.UserName = "hogefoobar"
		newUser.Email = "hogefoobar@hoge.com"
		newUser.Telephone = "YYY-YYYY-YYYY"

		err = userRepo.Create(ctx, newUser)
		if err != nil {
			return err
		}

		found.UserName = "new name"
		err = userRepo.Update(ctx, found)
		if err != nil {
			return err
		}

		err = userRepo.Delete(ctx, "1")
		if err != nil {
			return err
		}

		bs, err := json.Marshal(users)
		if err != nil {
			return err
		}

		pre := presenter.NewStdoutPresenter()
		_, _ = pre.Write(bs)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
