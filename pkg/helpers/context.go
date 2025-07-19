package helpers

import (
	"errors"
	"time"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	"github.com/gin-gonic/gin"
)

func getContextString(name string, c *gin.Context, deps *deps.Dependencies) (string, error) {
	context, exists := c.Get(name)
	if !exists {
		deps.Logger.Debug().Err(errors.New("context doenst exist")).Msgf(
			"error getting %s from context string",
			name,
		)
		return "", errors.New("error found, try again later")
	}

	contextStr, ok := context.(string)
	if !ok {
		deps.Logger.Debug().Err(errors.New("context dont is string")).Msgf(
			"error transforming %s from context in string",
			name,
		)
		return "", errors.New("error found, try again later")
	}
	return contextStr, nil
}

func getContextTime(name string, c *gin.Context, deps *deps.Dependencies) (time.Time, error) {
	context, exists := c.Get(name)
	if !exists {
		deps.Logger.Debug().Err(errors.New("context doenst exist")).Msgf(
			"error getting %s from context time",
			name,
		)
		return time.Now(), errors.New("error found, try again later")
	}

	contextStr, ok := context.(time.Time)
	if !ok {
		deps.Logger.Debug().Err(errors.New("context dont is time.Time")).Msgf(
			"error getting %s from context time",
			name,
		)
		return time.Now(), errors.New("error found, try again later")
	}
	return contextStr, nil
}

func GetAccessTokenContext(c *gin.Context, deps *deps.Dependencies) (string, error) {
	return getContextString("accessToken", c, deps)
}

func GetRefreshTokenContext(c *gin.Context, deps *deps.Dependencies) (string, error) {
	return getContextString("refreshToken", c, deps)
}

func GetExpireTokenContext(c *gin.Context, deps *deps.Dependencies) (time.Time, error) {
	return getContextTime("expireToken", c, deps)
}
