package ginvalidator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	valid "github.com/asaskevich/govalidator"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func ExampleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(ExampleMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	p := NewParam("person", "person is not valid.")

	router.GET("/hello/:person",
		// IsAlpha("person is not alpha").
		p.Chain().
			Not().
			IsArray("person is an array", &ArrayLengthCheckerOpts{}).
			Validate(),
		// p.Chain().
		// 	IsASCII("person is not ascii").
		// 	Bail().
		// 	Not().
		// 	IsAlphanumeric("").
		// 	Validate(),
		func(c *gin.Context) {
			person := c.Query("person")
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": person})
		})
	return router
}

func TestExampleMiddleware(t *testing.T) {
	router := setupRouter()

	// Test the /test route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	fmt.Println("isValid", valid.IsASCII("你好，世界！"))

	// fmt.Printf("json field(s) %v", splitFieldOnPeriod("name"))

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Equal(t, `{"message":"success"}`, w.Body.String())

	data := []byte(`{
	  "person": {
	    "name": {
	      "first": "Leonid",
	      "last": "Bugaev",
	      "fullName": "Leonid Bugaev",
	    },
	    "github": {
	      "handle": "buger",
	      "followers": 109
	    },
	    "avatars": [
	      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
	    ]
	  },
	  "company": {
	    "name": "Acme"
	  }
	}`)

	_ = data

	ty := "person"
	key, typ, _, _ := jsonparser.Get([]byte(fmt.Sprintf(`{"key":"%s"}`, ty)), "key")
	// key, typ, _, _ := jsonparser.Get(data, "person", "avatars", "[0]", "url")

	fmt.Printf("key is %s while datatype is %v", key, typ)
}

func TestParamMiddleware(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello/[]", nil)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// fmt.Println("Response:", w.Body.String())

	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Equal(t, `{"message":"david"}`, w.Body.String())
}
