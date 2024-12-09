package web

import (
	"fmt"
	"html/template"
	"io"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/StupidRepo/Senecessary/pkg/models"
	"github.com/StupidRepo/Senecessary/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func StartMux() {
	fmt.Println("Starting server on :2020")

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assignments := shared.User.Assignments

		sort.Slice(assignments, func(i, j int) bool {
			return assignments[i].DueDate.After(assignments[j].DueDate)
		})
		if len(assignments) > 4 {
			assignments = assignments[:4]
		}

		for i := range assignments {
			assignment := &assignments[i]

			var sections []models.Section
			allSections, err := shared.GetSectionsInCourse(assignment.Spec.CourseId)
			if err != nil {
				panic(err)
			}

			for _, section := range *allSections {
				for _, id := range assignment.Spec.SectionIds {
					if section.Id == id {
						sections = append(sections, section)
					}
				}
			}

			assignment.Sections = sections
		}

		RenderTemplate(w, "index.html", map[string]interface{}{
			"Assignments": assignments,
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
	assignmentId := string(body)

	var assignment *models.Assignment
	for i := range shared.User.Assignments {
		if shared.User.Assignments[i].Id == assignmentId {
			assignment = &shared.User.Assignments[i]
			break
		}
	}
	if assignment == nil {
		panic("Assignment not found with ID " + assignmentId)
	}

	fmt.Println("Solving assignment", assignment.Id)

	sessionId := uuid.New().String()
	sessionReq := models.SessionRequest{
		ClientVersion: "2.13.8",
		Platform:      "seneca",

		Modules: []models.AnswerModule{
			{
				ModuleId:  "1",
				CourseId:  assignment.Spec.CourseId,
				SectionId: assignment.Spec.SectionIds[0],
				SessionId: sessionId,
				Completed: true,
				Score:     1,
				Submitted: false,
			},
			{
				ModuleId:  "2",
				CourseId:  assignment.Spec.CourseId,
				SectionId: assignment.Spec.SectionIds[0],
				SessionId: sessionId,
				Completed: true,
				Score:     1,
				Submitted: false,
			},
			{
				ModuleId:  "3",
				CourseId:  assignment.Spec.CourseId,
				SectionId: assignment.Spec.SectionIds[0],
				SessionId: sessionId,
				Completed: true,
				Score:     1,
				Submitted: false,
			},
			{
				ModuleId:  "4",
				CourseId:  assignment.Spec.CourseId,
				SectionId: assignment.Spec.SectionIds[0],
				SessionId: sessionId,
				Completed: true,
				Score:     1,
				Submitted: false,
			},
		},
		Session: models.Session{
			SessionId: sessionId,
			CourseId:  assignment.Spec.CourseId,

			Completed: true,

			ModulesCorrect: 20,
			ModulesStudied: 20,
			ModulesTested:  20,

			SessionScore: 1,

			SectionIds: []string{
				assignment.Spec.SectionIds[3],
			},
			ContentIds: []string{},

			SessionType: "adaptive",

			TimeStarted:  time.Now(),
			TimeFinished: time.Now().Add(time.Duration(rand.Float64() * (float64(time.Minute) * 8))),
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
