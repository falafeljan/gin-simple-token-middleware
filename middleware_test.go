package tokenmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getFreePort() (string, error) {
	portNumber, err := freeport.GetFreePort()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(int64(portNumber), 10), nil
}

func performRequest(url string, queryToken *string, headerToken *string) (int, error) {
	queryString := ""

	if queryToken != nil {
		queryString = "access_token=" + *queryToken
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"?"+queryString, nil)
	if err != nil {
		return 0, err
	}

	req.Header = make(map[string][]string)
	if headerToken != nil {
		req.Header["Authorization"] = []string{"Token " + *headerToken}
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	parts := strings.Split(res.Status, " ")
	statusCode, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	return statusCode, nil
}

func startServer(port string, token string) (*http.Server, error) {
	router := gin.Default()
	router.GET("/", NewHandler(token), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"yee": "boi",
		})
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()

	return server, nil
}

func TestNewHandler(t *testing.T) {
	token := "foobar123"
	emptyValue := ""
	fooValue := "foo"

	gin.SetMode(gin.ReleaseMode)

	port, err := getFreePort()
	if err != nil {
		t.Fatal(err)
	}

	fixtures := [][2]*string{
		[2]*string{nil, nil},
		[2]*string{&emptyValue, nil},
		[2]*string{nil, &emptyValue},
		[2]*string{&emptyValue, &emptyValue},
		[2]*string{&fooValue, nil},
		[2]*string{nil, &fooValue},
		[2]*string{&fooValue, &fooValue},
		[2]*string{&token, nil},
		[2]*string{nil, &token},
		[2]*string{&token, &token},
	}

	expected := []int{
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusUnauthorized,
		http.StatusOK,
		http.StatusOK,
		http.StatusOK,
	}

	server, err := startServer(port, token)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	for i, fixture := range fixtures {
		received, err := performRequest("http://127.0.0.1:"+port+"/", fixture[0], fixture[1])
		if err != nil {
			t.Fatal(err)
		}

		if received != expected[i] {
			t.Errorf("#%d: Result did not match for request %+v:\nExpected: %d\nReceived: %d\n",
				i,
				fixture,
				expected[i],
				received,
			)

			server.Close()
		}
	}

	err = server.Close()
	if err != nil {
		t.Fatal(err)
	}
}
