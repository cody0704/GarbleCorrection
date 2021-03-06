package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/hydra13142/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type resContent struct {
	Status  string `json:"Status"`
	Content string `json:"Content"`
}

type content struct {
	Content string `json:"Content"`
}

type file struct {
	FileName string `json:"FileName"`
	Data     string `json:"Data"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/simplifiedGarbleds", simplifiedGarbleds).Methods("Post")
	router.HandleFunc("/simplifiedGarbled", simplifiedGarbled).Methods("Post")

	n := negroni.New()
	n.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}))
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(":8080", n))
}

func simplifiedGarbled(w http.ResponseWriter, r *http.Request) {
	var newContent content

	var response resContent
	response.Status = "false"
	response.Content = ""

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		json.NewEncoder(w).Encode(response)
	}

	err = json.Unmarshal(reqBody, &newContent)
	if err != nil {

		json.NewEncoder(w).Encode(response)
	} else {

		decoded, err := base64.StdEncoding.DecodeString(newContent.Content)

		if err == nil {
			var bytes = []byte(decoded)

			if !strings.Contains(chardet.Mostlike(bytes), "utf") &&
				!strings.Contains(chardet.Mostlike(bytes), "utf16") {

				response.Status = "true"

				// GBK 2 UTF-8
				s, _ := decodeGBK(bytes)
				bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
				data := string(bomUtf8) + string(s)

				response.Content = data
			}

			json.NewEncoder(w).Encode(response)
		}
	}
}

func simplifiedGarbleds(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

//convert GBK to UTF-8
func decodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//convert UTF-8 to GBK
func encodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//convert BIG5 to UTF-8
func decodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//convert UTF-8 to BIG5
func encodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
