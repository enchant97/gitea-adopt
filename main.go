package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type BasicAuth struct {
	Username string
	Password string
}

type API struct {
	BaseURL string
	Auth    BasicAuth
}

func (api *API) getUnadopted() ([]string, error) {
	url, _ := url.JoinPath(api.BaseURL, "/api/v1/admin/unadopted")
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(api.Auth.Username, api.Auth.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (api *API) adopt(org_repo string) error {
	url, _ := url.JoinPath(api.BaseURL, "/api/v1/admin/unadopted/", org_repo)
	req, _ := http.NewRequest("POST", url, nil)
	req.SetBasicAuth(api.Auth.Username, api.Auth.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("adopting %s failed with status code %d", org_repo, resp.StatusCode)
	}
	return nil
}

func main() {
	var baseURL string
	var username string
	var password string

	fmt.Print("Gitea URL: ")
	fmt.Scanln(&baseURL)
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	api := API{
		BaseURL: baseURL,
		Auth: BasicAuth{
			Username: username,
			Password: password,
		},
	}

	for {
		repos, err := api.getUnadopted()
		if len(repos) == 0 && err == nil {
			fmt.Println("No more repos to adopt")
			break
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, org_repo := range repos {
			fmt.Println("Adopting:", org_repo)
			if err := api.adopt(org_repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
