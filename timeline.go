package main

import (
	"encoding/json"
	"os"
	"path"

	"plato/db"
	"plato/debug"
	"plato/entity"
)

type Asset struct {
	Media   string `json:"media"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
}

type Timeline struct {
	Headline  string `json:"headline"`
	Type      string `json:"type"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Text      string `json:"text"`
	Asset     Asset  `json:"asset"`
}

func generateProjectTimelineJSON(user entity.User) {
	var ts []Timeline

	ps := getProjectsByAuthorID(user.ID())
	for _, p := range ps {
		var t Timeline

		t.Headline = p.Title()
		t.StartDate = p.StartDate().Format("2006,1,2")
		t.EndDate = p.EndDate().Format("2006,1,2")
		t.Text = p.ShortContent(140)
		t.Asset.Media = p.ImageURL()

		ts = append(ts, t)
	}

	v := struct {
		Headline string     `json:"headline"`
		Type     string     `json:"type"`
		Text     string     `json:"text"`
		Asset    Asset      `json:"asset"`
		Date     []Timeline `json:"date"`
	}{
		"Project Timeline",
		"default",
		"<p>Lorem ipsum dolor sit amet</p>",
		Asset{},
		ts,
	}

	vv := struct {
		Timeline interface{} `json:"timeline"`
	}{
		v,
	}

	data, err := json.Marshal(vv)
	if err != nil {
		debug.Warn(err)
		return
	}

	filepath := projectTimelinePath(user)
	if err = os.MkdirAll(path.Dir(filepath), 0700); err != nil {
		debug.Warn(err)
		return
	}

	f, err := os.OpenFile(filepath, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0600)
	if err !=  nil {
		debug.Warn(err)
		return
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		debug.Warn(err)
		return
	}
}

func projectTimelinePath(user entity.User) string {
	return db.DataDir + "/timeline/" + user.Email() + "/timeline.json"
}
