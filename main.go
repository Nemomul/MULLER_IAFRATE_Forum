package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	r := mux.NewRouter()
	r.HandleFunc("/discussions", discussionsByCategoryHandler).Methods("GET")
	http.Handle("/acceuil", r)

	http.HandleFunc("/login", loginFormHandler)
	http.HandleFunc("/create-discussion", createDiscussionHandler)
	http.HandleFunc("/home", homeFormHandler)
	http.HandleFunc("/register", registerHandler)
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

	mux.HandleFunc("/css/acceuil.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "acceuil.css")
	})

	mux.HandleFunc("/discussions/{categoryID}", discussionsByCategoryHandler)

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
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
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

	// Si tout va bien, renvoyez une réponse 200 (OK)
	w.WriteHeader(http.StatusOK)
}

type Discussion struct {
	ID   int
	Name string
	// Ajoutez d'autres champs si nécessaire
}

func discussionsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la catégorie depuis la requête
	categoryID := mux.Vars(r)["categoryID"]

	// Effectuer une requête dans la base de données pour récupérer les discussions correspondantes à la catégorie
	rows, err := db.Query("SELECT d.id_discussion, d.name_discussion FROM discussion d JOIN category_discussion cd ON d.id_discussion = cd.id_discussion WHERE cd.id_cat = ?", categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parcourir les résultats de la requête et stocker les discussions dans une structure de données appropriée
	var discussions []Discussion
	for rows.Next() {
		var discussion Discussion
		if err := rows.Scan(&discussion.ID, &discussion.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		discussions = append(discussions, discussion)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir les discussions en JSON
	jsonData, err := json.Marshal(discussions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Définir l'en-tête Content-Type sur application/json
	w.Header().Set("Content-Type", "application/json")
	// Renvoyer les données JSON
	w.Write(jsonData)
}

func discuByCategoryHandler(w http.ResponseWriter, r *http.Request) {

	http.HandleFunc("/getDiscussions", func(w http.ResponseWriter, r *http.Request) {
		idCatStr := r.URL.Query().Get("id_cat")
		if idCatStr == "" {
			http.Error(w, "Missing id_cat", http.StatusBadRequest)
			return
		}
		idCat, err := strconv.Atoi(idCatStr)
		if err != nil {
			http.Error(w, "Invalid id_cat", http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT d.* FROM discussions d JOIN category_discussions cd ON d.id = cd.discussion_id WHERE cd.id_cat = ?", idCat)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var discussions []map[string]interface{}
		columns, err := rows.Columns()
		if err != nil {
			http.Error(w, "Failed to get columns", http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			values := make([]sql.RawBytes, len(columns))
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				scanArgs[i] = &values[i]
			}
			err = rows.Scan(scanArgs...)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			row := make(map[string]interface{})
			for i, value := range values {
				if value == nil {
					row[columns[i]] = nil
				} else {
					row[columns[i]] = string(value)
				}
			}
			discussions = append(discussions, row)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Failed to read rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discussions)
	})
}
