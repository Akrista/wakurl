package wakurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const cytrusURL = "https://launcher.cdn.ankama.com/cytrus.json"

type Helper struct {
	version string
	files   map[string]string
	mode    string
}

func NewHelper(betaMode bool) (*Helper, error) {
	h := new(Helper)
	cy, err := getCytrus()
	if err != nil {
		return nil, err
	}

	h.mode = "main"
	h.version = cy.Games.Wakfu.Platforms.Darwin.Main
	if betaMode {
		h.mode = "beta"
		h.version = cy.Games.Wakfu.Platforms.Darwin.Beta
	}
	if err := h.getFiles(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Helper) getFiles() error {
	h.files = make(map[string]string)
	resp, err := http.Get(fmt.Sprintf("https://launcher.cdn.ankama.com/wakfu/releases/%s/darwin/%s.json", h.mode, h.version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`"([^"]+)":{"hash":"([^"]+)","size":\d+}`)
	matches := re.FindAllSubmatch(json, -1)
	for _, i := range matches {
		h.files[string(i[1])] = string(i[2])
	}

	return nil
}

func (h *Helper) GetURL(path string) string {
	if val, ok := h.files[path]; ok {
		return fmt.Sprintf("https://launcher.cdn.ankama.com/wakfu/hashes/%s/%s", val[:2], val)
	}
	panic(errors.New("The file does not exist."))
}

type Cytrus struct {
	Version int    `json:"version"`
	Name    string `json:"name"`
	Games   struct {
		Wakfu struct {
			Name   string `json:"name"`
			Order  int    `json:"order"`
			GameID int    `json:"gameId"`
			Assets struct {
				Meta struct {
					Beta string `json:"beta"`
					Main string `json:"main"`
				} `json:"meta"`
			} `json:"assets"`
			Platforms struct {
				Darwin struct {
					Beta string `json:"beta"`
					Main string `json:"main"`
				} `json:"darwin"`
				Linux struct {
					Beta string `json:"beta"`
					Main string `json:"main"`
				} `json:"linux"`
				Windows struct {
					Beta string `json:"beta"`
					Main string `json:"main"`
				} `json:"windows"`
			} `json:"platforms"`
		} `json:"wakfu"`
	} `json:"games"`
}

func getCytrus() (*Cytrus, error) {
	resp, err := http.Get(cytrusURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cy := new(Cytrus)
	if err := json.NewDecoder(resp.Body).Decode(cy); err != nil {
		return nil, err
	}

	return cy, nil
}
