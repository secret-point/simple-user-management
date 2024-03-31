
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID    int                `bson:"user_id" gorm:"type:int" json:"userID"`
    FirstName string             `bson:"first_name" gorm:"type:varchar(100)" json:"firstName"`
    LastName  string             `bson:"last_name" gorm:"type:varchar(100)" json:"lastName"`
    Email     string             `bson:"email" gorm:"type:varchar(100);unique_index" json:"email"`
    Age       int                `bson:"age" gorm:"type:int" json:"age"`
}
