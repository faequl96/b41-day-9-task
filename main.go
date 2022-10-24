package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/register", registerForm).Methods("GET")
	route.HandleFunc("/login", loginForm).Methods("GET")

	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")

	route.HandleFunc("/form-add-project", addProjectForm).Methods("GET")
	route.HandleFunc("/send-data-add-project", sendDataAddProject).Methods("POST")

	route.HandleFunc("/form-edit-project/{id}", editProjectForm).Methods("GET")
	route.HandleFunc("/send-data-edit-project/{id}", sendDataEditProject).Methods("POST")

	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	fmt.Println("Server running on localhost:8000")
	http.ListenAndServe("localhost:8000", route)
}

type projectDataStruc struct {
	Id           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Description  string
	Technologies []string
	Image        string
}

var projectData = []projectDataStruc{
	// {
	// 	Id:           0,
	// 	ProjectName:  "Dummy Project 1",
	// 	StartDate:    "2022-09-12",
	// 	EndDate:      "2022-09-19",
	// 	Duration:     "1 Weeks",
	// 	Description:  "Description Dummy Project 1",
	// 	Technologies: []string{"NodeJs", "ReactJs"},
	// 	Image:        "gambar1.jpg",
	// },
	// {
	// 	Id:           1,
	// 	ProjectName:  "Dummy Project 2",
	// 	StartDate:    "2022-09-20",
	// 	EndDate:      "2022-09-25",
	// 	Duration:     "5 Days",
	// 	Description:  "Description Dummy Project 2",
	// 	Technologies: []string{"NodeJs", "ReactJs", "TypeScript"},
	// 	Image:        "gambar2.jpg",
	// },
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/home.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {
		var result []projectDataStruc
		data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM db_project")
		for data.Next() {
			var each = projectDataStruc{}
			err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result = append(result, each)
		}

		response := map[string]interface{}{
			"ProjectData": result,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func registerForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		for index, selectedProject := range projectData {
			if index == where_id {
				selectedProjectData = projectDataStruc{
					ProjectName:  selectedProject.ProjectName,
					StartDate:    selectedProject.StartDate,
					EndDate:      selectedProject.EndDate,
					Duration:     selectedProject.Duration,
					Description:  selectedProject.Description,
					Technologies: selectedProject.Technologies,
					Image:        selectedProject.Image,
				}
			}
		}

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func addProjectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func sendDataAddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		var technologies []string
		technologies = r.Form["technologies"]
		image := r.PostForm.Get("project-image")

		fmt.Println(startDate)

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		fmt.Println(startDateParse)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		addProjectData := projectDataStruc{
			ProjectName: projectName,
			// StartDate:    startDate,
			// EndDate:      endDate,
			Duration:     duration,
			Description:  description,
			Technologies: technologies,
			Image:        image,
		}

		projectData = append(projectData, addProjectData)

		fmt.Println(projectData)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func editProjectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		for index, selectedProject := range projectData {
			if id == index {
				selectedProjectData = projectDataStruc{
					Id:          id,
					ProjectName: selectedProject.ProjectName,
					StartDate:   selectedProject.StartDate,
					EndDate:     selectedProject.EndDate,
					Description: selectedProject.Description,
				}
				fmt.Println(selectedProjectData.Description)
			}
		}

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func sendDataEditProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		index, _ := strconv.Atoi(mux.Vars(r)["id"])

		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		var technologies []string
		technologies = r.Form["technologies"]
		image := r.PostForm.Get("project-image")

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		editSelectedProjectData := projectDataStruc{
			ProjectName: projectName,
			// StartDate:    startDate,
			// EndDate:      endDate,
			Duration:     duration,
			Description:  description,
			Technologies: technologies,
			Image:        image,
		}

		projectData[index] = editSelectedProjectData

		fmt.Println(projectData)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["id"])

	projectData = append(projectData[:index], projectData[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
