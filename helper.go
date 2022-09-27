package whatsapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Media struct {
	MessagingProduct string `json:"messaging_product"`
	URL              string `json:"url"`
	MIMEType         string `json:"mime_type"`
	SHA256           string `json:"sha256"`
	FileSize         string `json:"file_size"`
	ID               string `json:"id"`
}

func (wa *Whatsapp) getMedia(mediaID string) (media Media, err error) {

	endpoint := fmt.Sprintf("https://graph.facebook.com/%s/%s", wa.APIVersion, mediaID)

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return media, err
	}

	req.Header.Set("Authorization", "Bearer "+wa.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return media, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&media)

	if err != nil {
		return media, err
	}
	return media, err
}

func (wa *Whatsapp) UploadMedia(filepath string) (id string, err error) {
	endpoint := fmt.Sprintf("https://graph.facebook.com/%s/%s/media", wa.APIVersion, wa.PhoneNumberID)

	// read filepath and get base64
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return id, err
	}

	// add messaging_product="whatsapp" to form data body
	resp, err := http.PostForm(endpoint, url.Values{
		"messaging_product": {"whatsapp"},
		"file":              {string(file)},
	})

	if err != nil {
		return id, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return id, err
	}
	fmt.Printf("%s\n", string(body))
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return id, err
	}
	log.Println("res", res)
	// id = res["id"]
	return id, err
}

// Http post request to send the message
func (wa *Whatsapp) sendMessage(request any) (res map[string]interface{}, err error) {

	marshaledJSON, err := json.Marshal(request)
	if err != nil {
		return res, err
	}
	reqString := string(marshaledJSON)

	log.Println("body", reqString)

	body := strings.NewReader(reqString)

	endpoint := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", wa.APIVersion, wa.PhoneNumberID)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return res, err
	}
	req.Header.Set("Authorization", "Bearer "+wa.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return res, err
	}

	return res, err
}