package main

import (
	"context"
	"day9/connection"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func handleRequests() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	// router path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact", Contact).Methods("GET")
	route.HandleFunc("/formProject", formProject).Methods("GET")
	route.HandleFunc("/detailProject/{id}", DetailProject).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("POST")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/formEditProject/{id}", formEditProject).Methods("GET")
	route.HandleFunc("/editProject/{id}", editProject).Methods("POST")


	fmt.Println("Go Running on Port 5000")
	http.ListenAndServe(":5000", route)
}

type Project struct {
	Name 						string
	StartDate 			time.Time
	EndDate 				time.Time
	StartDateFormat string
	EndDateFormat 	string
	DescFormat			string
	Duration				string
	Desc 						string
	Id							int
	Tech						[]string
}

var data = []Project{}


func Home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	var result []Project
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects")
	for data.Next() {
		var each = Project{}
		err := data.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Desc, &each.Tech)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// each.DescFormat = each.Desc[:200]

		each.Duration = ""

		day :=  24 //in hours
		month :=  24 * 30 // in hours
		year :=  24 * 365 // in hours
		differHour := each.EndDate.Sub(each.StartDate).Hours()
		var differHours int = int(differHour)
		days := differHours / day
		months := differHours / month
		years := differHours / year
		if differHours < month {
			each.Duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < year {
			each.Duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			each.Duration = strconv.Itoa(int(years)) + " Years"
		}

		result = append(result, each)
	}

	response := map[string]interface{}{
		"Projects": result,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		w.Write([]byte("Message: "))
		w.Write([]byte(err.Error()))
	}
}


func Contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func formProject(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/addProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func DetailProject(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")
	
	var tmpl, err = template.ParseFiles("views/detailProject.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Desc, &ProjectDetail.Tech)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	ProjectDetail.StartDateFormat = ProjectDetail.StartDate.Format("02 Jan 2006")
	ProjectDetail.EndDateFormat = ProjectDetail.EndDate.Format("02 Jan 2006")
	ProjectDetail.Duration = ""

	day :=  24 //in hours
	month :=  24 * 30 // in hours
	year :=  24 * 365 // in hours
	differHour := ProjectDetail.EndDate.Sub(ProjectDetail.StartDate).Hours()
	var differHours int = int(differHour)
	days := differHours / day
	months := differHours / month
	years := differHours / year
	if differHours < month {
		ProjectDetail.Duration = strconv.Itoa(int(days)) + " Days"
	} else if differHours < year {
		ProjectDetail.Duration = strconv.Itoa(int(months)) + " Months"
	} else if differHours > year {
		ProjectDetail.Duration = strconv.Itoa(int(years)) + " Years"
	}

	response := map[string]interface{}{
		"Details" : ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}


func addProject(w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var inputTitle			string
	var inputStartDate 	string
	var inputEndDate 		string
	var inputDesc 			string
	var inputTech				[]string

	for i, values := range r.Form {
		for _ , value := range values {
			if i == "inputTitle" {
				inputTitle = value
			}
			if i == "inputStartDate" {
				inputStartDate = value
			}
			if i == "inputEndDate" {
				inputEndDate = value
			}
			if i == "inputDesc" {
				inputDesc = value
			}
			if i == "inputTech" {
				inputTech = append(inputTech, value)
			}
		}
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies) VALUES ($1, $2, $3, $4, $5) ", inputTitle, inputStartDate, inputEndDate, inputDesc, inputTech)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
	}

	http.Redirect(w,r, "/", http.StatusMovedPermanently)
}


func editProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var inputTitle string
	var inputStartDate string
	var inputEndDate string
	var inputDesc string
	var inputTech []string

	for i, values := range r.Form {
		for _, value := range values {
			if i == "inputTitle"{
				inputTitle = value
			}
			if i == "inputStartDate"{
				inputStartDate = value
			}
			if i == "inputEndDate"{
				inputEndDate = value
			}
			if i == "inputDesc" {
				inputDesc = value
			}
			if i == "inputTech" {
				inputTech = append(inputTech, value)
			}
		}

	}
		_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5 WHERE id=$6", inputTitle, inputStartDate, inputEndDate, inputDesc, inputTech, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/editMyProject.html")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectEdit := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", id).Scan(&ProjectEdit.Id, &ProjectEdit.Name, &ProjectEdit.StartDate, &ProjectEdit.EndDate, &ProjectEdit.Desc, &ProjectEdit.Tech)

	ProjectEdit.StartDateFormat = ProjectEdit.StartDate.Format("2006-01-02")
	ProjectEdit.EndDateFormat = ProjectEdit.EndDate.Format("2006-01-02")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
	}

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

		if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	handleRequests() 

	
}