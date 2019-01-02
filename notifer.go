package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type postParams struct {
	ProjectFolder string `json:"project_folder"`
	File          string `json:"file"`
	Changed       int64  `json:"changed"`
	Ext           string `json:"ext"`
}

type NetNotifier struct {
	CallUrl string
	CanPost bool
}

func newNetNotifier(callUrl string) *NetNotifier {
	callPost := true
	if strings.TrimSpace(callUrl) == "" {
		callPost = false
	}
	return &NetNotifier{
		CallUrl: callUrl,
		CanPost: callPost,
	}
}

func (n *NetNotifier) Put(cf *changeFile) {
	if !n.CanPost {
		log.Println("notifier call url ignore. ", n.CallUrl)
		return
	}
	n.dispatch(&postParams{
		ProjectFolder: projectFolder,
		File:          cf.Name,
		Changed:       cf.Changed,
		Ext:           cf.Ext,
	})
}

func (n *NetNotifier) dispatch(params *postParams) {
	b, err := json.Marshal(params)
	if err != nil {
		log.Println("error: json.Marshal n.params. ", err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", n.CallUrl, bytes.NewBuffer(b))
	if err != nil {
		log.Println("error: http.NewRequest. ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", "FileBoy Net Notifier v1.2")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		// todo retry???
	}
	log.Println("notifier done .")
}
