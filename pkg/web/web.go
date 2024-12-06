package web

import (
	"encoding/json"
	"fmt"
	"github.com/StupidRepo/Senecessary/pkg/models"
	"github.com/StupidRepo/Senecessary/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"time"
)

func StartMux() {
	fmt.Println("Starting server on :2020")

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, result, err := shared.DoReq[models.AssignmentResponse]("GET", "https://assignments.app.senecalearning.com/api/students/me/assignments?limit=100", nil)
		if err != nil {
			panic(err)
		}

		sort.Slice(result.Items, func(i, j int) bool {
			return result.Items[i].DueDate.After(result.Items[j].DueDate)
		})

		RenderTemplate(w, "index.html", map[string]interface{}{
			"Assignments": result.Items,
		})
	})
	r.HandleFunc("/solve", SolveAssignment)

	http.Handle("/", ErrMiddleware(r))
	err := http.ListenAndServe(":2020", nil)

	if err != nil {
		panic(err)
	}
}

func SolveAssignment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// parse body as assignment
	assignment := models.Assignment{}
	err = json.Unmarshal(body, &assignment)
	if err != nil {
		panic(err)
	}

	sessionId := uuid.New().String()
	sessionReq := models.SessionRequest{
		ClientVersion: "2.13.8",
		Platform:      "seneca",

		Modules: []models.Module{},
		Session: models.Session{
			SessionId: sessionId,
			CourseId:  assignment.Spec.CourseId,

			Completed: true,

			ModulesCorrect: 20,
			ModulesStudied: 20,
			ModulesTested:  20,

			SessionScore: 1,

			SectionIds: []string{
				assignment.Spec.SectionIds[0],
			},
			ContentIds: []string{},

			SessionType: "adaptive",

			TimeStarted:  time.Now(),
			TimeFinished: time.Now().Add(time.Minute * 7),
		},

		UserId: assignment.UserId,
	}

	res, _, err := shared.DoReq[any]("POST", "https://stats.app.senecalearning.com/api/stats/sessions", sessionReq)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic("Error: " + res.Status)
	}

	RenderTemplate(w, "assignment.html", nil)
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("Monday, 02 Jan 2006")
		},
	}

	tmpl := template.Must(template.New(name).Funcs(funcMap).ParseFiles(filepath.Join("pkg", "web", "tmpl", name)))
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func HandleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func ErrMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				HandleError(w, fmt.Errorf("panic: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
