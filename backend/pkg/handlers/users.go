package handlers

import (
    "context"
    "net/http"
    "time"
    "strconv"
    "math"
    // "fmt"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "arritech-challenge/pkg/models"
)

// UserHandler struct holds a connection to the database
type UserHandler struct {
    collection *mongo.Collection
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *mongo.Client) *UserHandler {
    return &UserHandler{
        collection: db.Database("arritech").Collection("users"),
    }
}

// CreateUser - POST /users
func (handler *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var existingUser models.User
    err := handler.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
    if err != nil && err != mongo.ErrNoDocuments {
        // If we have an error that is not the ErrNoDocuments, it means we had an issue checking the email
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing email"})
        return
    }

    if existingUser.Email != "" {
        // If we found a user with that email, return an error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
        return
    }

    count, err := handler.collection.CountDocuments(ctx, bson.D{{}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
        return
    }
    user.UserID = int(count) + 1

    result, err := handler.collection.InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"result": result, "userID": user.UserID})
}

// GetUser - GET /users/:id
func (handler *UserHandler) GetUser(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    if err := handler.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

// UpdateUser - PUT /users/:id
func (handler *UserHandler) UpdateUser(c *gin.Context) {
    id, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }
    var updateUser models.User
    if err := c.BindJSON(&updateUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check for email duplication with other users
    emailDuplicationFilter := bson.M{"email": updateUser.Email, "_id": bson.M{"$ne": id}}
    count, err := handler.collection.CountDocuments(ctx, emailDuplicationFilter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for email duplication"})
        return
    }
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists with another user"})
        return
    }

    // Check for userID duplication with other users, if userID is being changed
    if updateUser.UserID != 0 {
        userIDDuplicationFilter := bson.M{"user_id": updateUser.UserID, "_id": bson.M{"$ne": id}}
        count, err = handler.collection.CountDocuments(ctx, userIDDuplicationFilter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for userID duplication"})
            return
        }
        if count > 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "UserID already exists with another user"})
            return
        }
    }

    updateResult, err := handler.collection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.D{{"$set", bson.D{
            {"user_id", updateUser.UserID},
            {"first_name", updateUser.FirstName},
            {"last_name", updateUser.LastName},
            {"email", updateUser.Email},
            {"age", updateUser.Age},
        }}},
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    if updateResult.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}



func (handler *UserHandler) DeleteUser(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    result, err := handler.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.JSON(http.StatusOK, result)
}

// GetUsers - GET /users
func (handler *UserHandler) GetUsers(c *gin.Context) {
    var page, pageSize int64 = 1, 10
    search := c.Query("search")

    if p, err := strconv.ParseInt(c.Query("page"), 10, 64); err == nil && p > 0 {
        page = p
    }
    if ps, err := strconv.ParseInt(c.Query("pageSize"), 10, 64); err == nil && ps > 0 {
        pageSize = ps
    }

    searchFilter := bson.M{}
    if search != "" {
        // Use $or to search across multiple fields
        regex := bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}} // Case-insensitive search
        searchFilter = bson.M{"$or": []interface{}{
            bson.M{"first_name": regex},
            bson.M{"last_name": regex},
            bson.M{"email": regex},
            bson.M{"age": regex},
        }}
    }

    var users []models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Calculating the number of documents to skip
    skip := (page - 1) * pageSize

    // Finding documents with pagination
    findOptions := options.Find().SetSkip(skip).SetLimit(pageSize)
    cursor, err := handler.collection.Find(ctx, searchFilter, findOptions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user data"})
            return
        }
        users = append(users, user)
    }

    // Optional: Count total documents for pagination metadata (e.g., total pages)
    total, err := handler.collection.CountDocuments(ctx, searchFilter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":      users,
        "total":     total,
        "page":      page,
        "pageSize":  pageSize,
        "totalPages": math.Ceil(float64(total) / float64(pageSize)),
    })
}
