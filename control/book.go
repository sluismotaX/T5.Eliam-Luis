package control

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// ---------------------------------------------------
// Book table contains the information for each book and language, editor and author

type Jperson struct {
	Id     uint32 `db:"id" bson:"id,omitempty"`
	Cuenta string `db:"cuenta" bson:"cuenta"`
	Uuid   string `db:"uuid" bson:"uuid,omitempty"`
	Nivel  uint32 `db:"nivel" bson:"nivel"`
	Email  string `db:"email" bson:"email"`
}

type User struct {
	Cuenta   string `json:"cuenta"`
	Password string `json:"password"`
}

type Server struct {
	Hostname string `json:"hostname"`
}

type BookZ struct {
	Title    string `db:"title" bson:"title"`
	Comment  string `db:"comment" bson:"comment"`
	Year     uint32 `db:"year" bson:"year"`
	Author   string `db:"author" bson:"author,omitempty"`
	Editor   string `db:"editor" bson:"editor,omitempty"`
	Language string `db:"language" bson:"language,omitempty"`
}

var (
	formato = "2006-01-02"
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Sentencia Case en codigo no es correcta.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result  no encontrado.")
	// ErrUnavailable is a database not available error
	ErrUnauthorized = errors.New("Usuario sin permiso para realizar esta operacion.")
	bookClient      http.Client
)

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Creando Cookie %s\n", err.Error())
	}
	bookClient = http.Client{
		Jar:     jar,
		Timeout: time.Second * 2,
	}
}

// getBody -> from server
func getBody(url string) (body []byte) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("getBody 1", url, err)
		log.Fatal(err)
	}
	res, getErr := bookClient.Do(req)
	if getErr != nil {
		fmt.Println("getBody 2", url, err)
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("getBody 3", url, err)
		log.Fatal(readErr)
	}
	return
}

// JLoginGET : get user id
func JLoginGET(server Server, user User) (id uint32) {
	cuenta := user.Cuenta
	passwd := user.Password
	encPass := base64.StdEncoding.EncodeToString([]byte(passwd))
	url := server.Hostname + "/jlogin/" + cuenta + "/" + encPass
	body := getBody(url)
	pers := Jperson{}
	jsonErr := json.Unmarshal(body, &pers)
	if jsonErr != nil {
		fmt.Println("JLoginGet", jsonErr)
		log.Fatal(jsonErr)
	}
	id = pers.Id
	//fmt.Println(body)
	return
}
