package types

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	Email             string `bson:"email" json:"email"`
	FullName          string `bson:"fullName" json:"fullName"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
	Admin             bool   `bson:"admin" json:"admin"`
}
