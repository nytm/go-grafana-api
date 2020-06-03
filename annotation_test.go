package gapi

import (
	"net/url"
	"testing"

	"github.com/gobs/pretty"
)

const (
	annotationsJSON = `[{
		"id": 1124,
		"alertId": 0,
		"dashboardId": 468,
		"panelId": 2,
		"userId": 1,
		"userName": "",
		"newState": "",
		"prevState": "",
		"time": 1507266395000,
		"text": "test",
		"metric": "",
		"regionId": 1123,
		"type": "event",
		"tags": [
			"tag1",
			"tag2"
		],
		"data": {}
	}]`

	newAnnotationJSON = `{
		"message":"Annotation added",
		"id": 1,
		"endId": 2
	}`

	newGraphiteAnnotationJSON = `{
		"message":"Annotation added",
		"id": 1
	}`

	deleteAnnotationJSON = `{"message":"Annotation deleted"}`
)

func TestAnnotations(t *testing.T) {
	server, client := gapiTestTools(200, annotationsJSON)
	defer server.Close()

	params := url.Values{}
	params.Add("from", "1506676478816")
	params.Add("to", "1507281278816")
	params.Add("limit", "100")

	as, err := client.Annotations(params)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(as))

	if as[0].ID != 1124 {
		t.Error("annotations response should contain annotations with an ID")
	}
}

func TestNewAnnotation(t *testing.T) {
	server, client := gapiTestTools(200, newAnnotationJSON)
	defer server.Close()

	a := Annotation{
		DashboardID: 123,
		PanelID:     456,
		Time:        1507037197339,
		IsRegion:    true,
		TimeEnd:     1507180805056,
		Tags:        []string{"tag1", "tag2"},
		Text:        "text description",
	}
	res, err := client.NewAnnotation(&a)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(res))

	if res != 1 {
		t.Error("new annotation response should contain the ID of the new annotation")
	}
}

func TestNewGraphiteAnnotation(t *testing.T) {
	server, client := gapiTestTools(200, newGraphiteAnnotationJSON)
	defer server.Close()

	a := GraphiteAnnotation{
		What: "what",
		When: 1507180805056,
		Tags: []string{"tag1", "tag2"},
		Data: "data",
	}
	res, err := client.NewGraphiteAnnotation(&a)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(res))

	if res != 1 {
		t.Error("new annotation response should contain the ID of the new annotation")
	}
}

func TestDeleteAnnotation(t *testing.T) {
	server, client := gapiTestTools(200, deleteAnnotationJSON)
	defer server.Close()

	res, err := client.DeleteAnnotation(1)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(res))

	if res != "Annotation deleted" {
		t.Error("delete annotation response should contain the correct response message")
	}
}

func TestDeleteAnnotationByRegionID(t *testing.T) {
	server, client := gapiTestTools(200, deleteAnnotationJSON)
	defer server.Close()

	res, err := client.DeleteAnnotationByRegionID(1)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(res))

	if res != "Annotation deleted" {
		t.Error("delete annotation by region ID response should contain the correct response message")
	}
}