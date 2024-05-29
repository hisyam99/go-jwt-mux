package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NamaLengkap string             `bson:"nama_lengkap" json:"nama_lengkap"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password" json:"password"`
}
