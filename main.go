package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username   string
	Password   string
	email      string
	sexe       string
	name       string
	first_name string
	birth_date string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/muller-iafrate-forum")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/login", loginFormHandler)
	http.HandleFunc("/home", homeFormHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/discussion", createDiscussionHandler)
	http.HandleFunc("/createurs", creatorHandler)
	http.ListenAndServe(":9000", nil)

}

func handler(mux *http.ServeMux) {
	mux.HandleFunc("/css/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "main.css")
	})

	mux.HandleFunc("/css/connexion.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "connexion.css")
	})

	mux.HandleFunc("/connexion.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "connexion.js")
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	sexe := r.FormValue("sexe")
	name := r.FormValue("name")
	first_name := r.FormValue("first_name")
	birth_date := r.FormValue("birth_date")

	// Vérifier si l'utilisateur existe déjà dans la base de données
	if userExists(username) {
		http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusBadRequest)
		return
	}

	// Insérer l'utilisateur dans la base de données
	err := insertUser(username, email, password, sexe, name, first_name, birth_date)
	if err != nil {
		log.Println("Erreur lors de l'enregistrement:", err)
		http.Error(w, "Erreur lors de l'enregistrement", http.StatusInternalServerError)
		return
	}

	// L'utilisateur est enregistré avec succès
	// Vous pouvez effectuer d'autres actions ici, par exemple, définir une session ou rediriger vers une page d'accueil

	http.Redirect(w, r, "/login", http.StatusFound)
}

func insertUser(username, mail, password, sexe, name, first_name, birth_date string) error {
	_, err := db.Exec("INSERT INTO users (username, mail, password, sexe, name, first_name, birth_date) VALUES (?, ?, ?, ?, ?, ?, ?)", username, mail, password, sexe, name, first_name, birth_date)
	return err
}

func loginRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.Path == "/login" {
			loginHandler(w, r)
		} else if r.URL.Path == "/register" {
			registerHandler(w, r)
		}
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/muller-iafrate-forum")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Exécuter une requête pour vérifier les informations de connexion
	row := db.QueryRow("SELECT id_users FROM users WHERE username=? AND password=?", username, password)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Échec de la connexion:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// L'utilisateur est connecté avec succès
	// Vous pouvez effectuer d'autres actions ici, par exemple, définir une session ou rediriger vers une page d'accueil

	http.Redirect(w, r, "/home", http.StatusFound)
}
func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode de la requête
	if r.Method != http.MethodPost {
		// Affichez le formulaire de connexion
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Le formulaire a été soumis, appelez la fonction de gestion de la soumission du formulaire
	loginHandler(w, r)
}

func homeFormHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategories()
	if err != nil {
		http.Error(w, "Impossible de charger les catégories", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("acceuil.html"))
	tmpl.Execute(w, categories)
}

func creatorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl1 := template.Must(template.ParseFiles("créateurs.html"))
	tmpl1.Execute(w, nil)
}

func userExists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'utilisateur:", err)
		return false
	}

	return count > 0
}

func authenticateUser(username, password string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=? AND password=?", username, password).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de l'authentification:", err)
		return false
	}

	return count > 0
}

/* CATEGORIE */

type Category struct {
	Id    string
	Genre string
}

func getCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id_cat, gender FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Genre); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func createDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode de la requête
	if r.Method != http.MethodPost {
		// Si ce n'est pas une requête POST, renvoyez une erreur 405 (Méthode non autorisée)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Récupérez les valeurs du formulaire
	nameDiscussion := r.FormValue("name_discussion")
	dateStart := r.FormValue("date_start")
	idUsers, err := strconv.Atoi(r.FormValue("id_users"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insérez les données dans la base de données
	_, err = db.Exec("INSERT INTO discussion (name_discussion, date_start, id_users) VALUES (?, ?, ?)", nameDiscussion, dateStart, idUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl1 := template.Must(template.ParseFiles("discussion.html"))
	tmpl1.Execute(w, nil)
	// Si tout va bien, redirigez vers la page discussion.html
	http.Redirect(w, r, "/discussion", http.StatusSeeOther)
}
