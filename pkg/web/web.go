package web

import (
	"fmt"
	"html/template"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
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

		RenderTemplate(w, "index.html", map[string]interface{}{
			"Assignments": assignments,
		})
	})
	r.HandleFunc("/solve", SolveAssignment)
	r.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		shared.RefreshAssessments()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

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

	fmt.Println("Solving assignment", assignment.Id, "with", len(assignment.Sections), "section(s)!")
	if os.Getenv("FORCE_SOLVE") != "true" {
		fmt.Println("Uh... there's a problem")

		fmt.Println("")
		fmt.Println("__SO WHAT IS THE PROBLEM?__")
		fmt.Println("Basically, the code that runs below these messages (method: web.SolveAssignment()) sends a lot of API requests to Seneca, which don't currently set the completion of a section to 100%.")
		fmt.Println("The issue we are running into is that the API requests don't specify a proper user answer to the questions, so the completion of the section doesn't get set to 100%.")
		fmt.Println("If you have any Golang experience and are willing to look into the Seneca API, please consider PRing a fix to this issue.")
		fmt.Println("I honestly have no idea how to go about implementing fake answers for every different answer type (wordfill, toggles, etc.), so I'm stuck.")

		fmt.Println("")
		fmt.Println("Anyway, like I said, please PR a fix if you can. Thanks!")
		return
	}
	for _, section := range assignment.Sections {
		sessionId := uuid.New().String()
		contentModules, err := shared.GetModulesInSection(assignment.Spec.CourseId, section.Id)
		if err != nil {
			panic(err)
		}

		fmt.Println("Solving section", section.Id)

		var answerModules []models.AnswerModule

		for i, module := range *contentModules {
			fmt.Printf("Solving module %s\n", module.Id)
			answerModule := models.AnswerModule{
				ModuleId:    module.Id,
				ContentId:   module.ParentId,
				CourseId:    module.CourseId,
				SectionId:   section.Id,
				SessionId:   sessionId,
				ModuleType:  module.ModuleType,
				ModuleOrder: i,
				ModuleScore: models.AnswerModuleScore{
					Score: 1,
				},
				Content:       []struct{}{},
				TestingActive: true,
				Completed:     true,
				Submitted:     true,
				Score:         1,
				TimeStarted:   time.Now(),
				TimeFinished:  time.Now().Add(time.Duration(rand.Float64() * (float64(time.Minute) * 8))),
			}
			answerModules = append(answerModules, answerModule)
		}

		sessionReq := models.SessionRequest{
			ClientVersion: "2.13.8",
			Platform:      "seneca",

			Modules: answerModules,
			Session: models.Session{
				SessionId: sessionId,
				CourseId:  assignment.Spec.CourseId,

				Completed: true,

				ModulesCorrect: 20,
				ModulesStudied: 20,
				ModulesTested:  20,

				SessionScore: 1,

				SectionIds: []string{
					section.Id,
				},
				ContentIds: []string{},

				SessionType: "adaptive",

				TimeStarted:  time.Now(),
				TimeFinished: time.Now().Add(time.Duration(rand.Float64() * (float64(time.Minute) * 8))),
			},

			UserId: assignment.UserId,
		}

		res, _, err := shared.DoReq[any]("POST", string(shared.Sessions_Submit), sessionReq)
		if err != nil {
			fmt.Print(err)
			panic(err)
		}

		if res.StatusCode != 200 {
			fmt.Print(res.Status)
			panic("Error: " + res.Status)
		}
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
