package handlers

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/types-library"

	"os"
	"testing"
)

var (
	user  types.User
	token string
)

func TestMain(main *testing.M) {
	database.Connect(os.Getenv("DATABASE_CONNECTION"))
	database.Create()

	user = types.NewUser(nick, "", email)
	content = types.NewContent("", id, "png", nil, true, true)

	var err error
	if err = database.WriteUser(user.Map()); err != nil {
		panic(err)
	}

	if token, _, err = database.CreateToken(user.ID); err != nil {
		panic(err)
	}

	os.Exit(main.Run())
}
