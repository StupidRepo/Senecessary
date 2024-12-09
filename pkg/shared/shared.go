package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/StupidRepo/Senecessary/pkg/models"
	"github.com/google/uuid"
)

var Client http.Client
var User models.User

type URLs string

const (
	Courses_SectionsQuery URLs = "https://course-cdn-v2.app.senecalearning.com/api/courses/%s/sections?limit=3000"
	User_MeQuery          URLs = "https://user-info.app.senecalearning.com/api/user-info/me"
	Sessions_Submit       URLs = "https://stats.app.senecalearning.com/api/stats/sessions"
)

func Login() *models.User {
	res, result, err := DoReq[models.User]("GET", string(User_MeQuery), nil)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			panic("Invalid token.")
		}
		panic("Error: " + res.Status)
	}

	User = result
	return &result
}

func RefreshAssessments() {
	_, result, err := DoReq[models.AssignmentResponse]("GET", "https://assignments.app.senecalearning.com/api/students/me/assignments?limit=1000", nil)
	if err != nil {
		panic(err)
	}

	User.Assignments = result.Items
}

func GetSectionsInCourse(CourseId string) (*[]models.Section, error) {
	_, result, err := DoReq[models.CoursesSectionsResponse]("GET", fmt.Sprintf(string(Courses_SectionsQuery), CourseId), nil)
	if err != nil {
		return nil, err
	}

	sections := result.Sections
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Number < sections[j].Number
	})

	return &sections, nil
}

func DoReq[T any](method string, url string, body interface{}) (*http.Response, T, error) {
	var result T
	var bodyString []byte = nil

	if body != nil {
		str, err := json.Marshal(body)
		if err != nil {
			return nil, result, err
		}

		bodyString = str
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyString))
	if err != nil {
		return nil, result, err
	}

	req.Header.Add("origin", "https://app.senecalearning.com")
	req.Header.Add("referer", "https://app.senecalearning.com/")

	req.Header.Add("user-agent", "Senecessary/1.0/Made in Golang by Bradlee Barnes")
	req.Header.Add("user-region", "GB")

	req.Header.Add("access-key", os.Getenv("TOKEN"))
	req.Header.Add("correlationId", GenerateCorrelationId())

	req.Header.Add("content-type", "application/json")

	res, err := Client.Do(req)
	if err != nil {
		return nil, result, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, result, err
	}

	return res, result, nil
}

func GenerateCorrelationId() string {
	timestamp := time.Now().UnixMilli()
	gen := fmt.Sprintf("%d::%s", timestamp, uuid.New())

	//fmt.Println("Generated Correlation ID: ", gen)
	return gen
}
