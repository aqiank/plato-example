package main

import (
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

func getProjectTimeline(user entity.User) interface{} {
	var ts []Timeline

	ps := getProjectsByAuthorID(user.ID())
	for _, p := range ps {
		var t Timeline

		t.Headline = p.Title()
		t.StartDate = p.StartDate().Format("2006,1,2")
		t.EndDate = p.EndDate().Format("2006,1,2")
		t.Text = p.ShortContent(64)
		t.Asset.Media = p.ImageURL()

		ts = append(ts, t)
	}

	return struct {
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
}
