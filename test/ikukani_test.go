package ikukani

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/softpunks/ikukani"
)

const (
	token      = "5a6a5234-a392-4a87-8f3f-33342afe8a42"
	apiVersion = "20170710"
	mockUser   = `
{
  "object": "user",
  "url": "https://api.wanikani.com/v2/user",
  "data_updated_at": "2018-04-06T14:26:53.022245Z",
  "data": {
    "id": "5a6a5234-a392-4a87-8f3f-33342afe8a42",
    "username": "example_user",
    "level": 5,
    "profile_url": "https://www.wanikani.com/users/example_user",
    "started_at": "2012-05-11T00:52:18.958466Z",
    "current_vacation_started_at": null,
    "subscription": {
      "active": true,
      "type": "recurring",
      "max_level_granted": 60,
      "period_ends_at": "2018-12-11T13:32:19.485748Z"
    },
    "preferences": {
      "default_voice_actor_id": 1,
      "lessons_autoplay_audio": false,
      "lessons_batch_size": 10,
      "lessons_presentation_order": "ascending_level_then_subject",
      "reviews_autoplay_audio": false,
      "reviews_display_srs_indicator": true
    }
  }
}`
	mockSummary = `
{
  "object": "report",
  "url": "https://api.wanikani.com/v2/summary",
  "data_updated_at": "2018-04-11T21:00:00.000000Z",
  "data": {
    "next_reviews_at": "2018-04-11T21:00:00.000000Z"
  }
}
`
)

func TestGetUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			t.Error("Missing api token")
		}
		if r.Header.Get("Wanikani-Revision") != apiVersion {
			t.Error("Missing api version")
		}
		w.Write([]byte(mockUser))
	}))
	defer ts.Close()

	client := ikukani.NewClient(token, apiVersion)
	client.Conn = ts.Client()
	client.BaseUrl = ts.URL

	user, err := client.GetUser()

	if err != nil {
		t.Error("expected nil got", err)
	}

	if user.Level != 5 {
		t.Error("expected 5 got ", user.Level)
	}
	if user.Id != token {
		t.Error("expected", token, "got ", user.Id)
	}
	if user.Username != "example_user" {
		t.Error("expected example_user got ", user.Username)
	}

}

func TestGetSummary(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			t.Error("Missing api token")
		}
		if r.Header.Get("Wanikani-Revision") != apiVersion {
			t.Error("Missing api version")
		}
		w.Write([]byte(mockSummary))
	}))
	defer ts.Close()

	client := ikukani.NewClient(token, apiVersion)
	client.Conn = ts.Client()
	client.BaseUrl = ts.URL

	summary, err := client.GetSummary()

	if err != nil {
		t.Error("expected nil got", err)
	}

	if summary.NextReviewsAt != "2018-04-11T21:00:00.000000Z" {
		t.Error("expected 2018-04-11T21:00:00.000000Z got ", summary.NextReviewsAt)
	}

}