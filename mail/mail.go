package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

const api = "https://burner.kiwi/api/v1/"

var client = &http.Client{}

type Inbox struct {
	Email string
	Id string
	Token string
}

type Message struct {
	Sender string `json:"sender"`
	Subject string `json:"subject"`
	Time int `json:"received_at"`
	BodyHtml string `json:"body_html"`
	BodyPlain string `json:"body_plain"`
}

func GetMessagesByInbox(inbox *Inbox) []*Message {
	request, err := http.NewRequest("GET", api + "inbox/" + inbox.Id + "/messages", bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error whilst creating a new request: ", err)
		return nil
	}
	request.Header.Set("X-Burner-Key", inbox.Token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error whilst doing a request: ", err)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	resp.Body.Close()

	var messages []*Message
	for _, result := range gjson.Get(bodyString, "result").Array() {
		var message Message
		json.Unmarshal([]byte(result.String()), &message)
		messages = append(messages, &message)
	}
	return messages
}

func GenerateInbox() *Inbox {
	resp, err := http.Get(api + "inbox")
	if err != nil{
		fmt.Println(err)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	resp.Body.Close()

	return &Inbox{
		Email: gjson.Get(bodyString, "result.email.address").String(),
		Id: gjson.Get(bodyString, "result.email.id").String(),
		Token: gjson.Get(bodyString, "result.token").String(),
	}
}
