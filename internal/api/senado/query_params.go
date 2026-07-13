package senado

import "github.com/gin-gonic/gin"

func QueryParams(c *gin.Context) map[string]string {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 && v[0] != "" {
			params[k] = v[0]
		}
	}
	return params
}
