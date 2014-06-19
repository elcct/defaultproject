package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/golang/glog"
	"time"
)

type User struct {
  ID       bson.ObjectId `bson:"_id,omitempty"`
  Email string `bson:"e"`
  Username string `bson:"u"`
  Password []byte  `bson:"p"`
  Timestamp time.Time `bson:"t"`
}

func (user *User) HashPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		glog.Fatalf("Couldn't hash password: %v", err)
		panic(err)
	}
	user.Password = hash
}

func GetUserByEmail(database *mgo.Database, email string) (user *User) {
	err := database.C("users").Find(bson.M{"e": email}).One(&user)

	if err != nil {
		glog.Warningf("Can't get user by email: %v", err)
	}	
	return
}

func InsertUser(database *mgo.Database, user *User) error {	
	user.ID = bson.NewObjectId()
	return database.C("users").Insert(user)	
}
