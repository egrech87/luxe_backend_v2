package middleware

import (
	"Quote/db/daoPostgres/daos/userDao"
	"Quote/db/entity/entities/userEntity"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		unauthorize(c)
		return
	}

	authId, err := strconv.ParseUint(username, 10, 64)
	if err != nil {
		unauthorize(c)
		return
	}

	friendlyName := password

	userI, found, err := userDao.GetInstance().ReadByID(context.Background(), authId, nil)
	if !found || err != nil {
		unauthorize(c)
		return
	}

	if (userI.(*userEntity.UserEntity)).GetFriendlyName() != friendlyName {
		unauthorize(c)
		return
	}

	// Continue to the next handler if authentication is successful
	c.Next()
}

func unauthorize(c *gin.Context) {
	c.JSON(401, gin.H{"message": "Unauthorized"})
	c.Abort()
}
