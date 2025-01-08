package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	InitializeDB()
	CloseDB()

	err := InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Define routes
	http.HandleFunc("/ws", WebSocketHandler)      // WebSocket connection
	http.HandleFunc("/register", RegisterHandler) // User registration
	http.HandleFunc("/login", LoginHandler)       // User login
	http.HandleFunc("/posts", PostsHandler)       // Posts (create and retrieve)

	http.HandleFunc("/navbar.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/navbar.html")
	})

	// Serve the `css` and `js` folders
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../frontend/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../frontend/js"))))

	// Serve HTML pages
	http.HandleFunc("/homepage.html", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving: ../frontend/homepage.html")
		http.ServeFile(w, r, "../frontend/homepage.html")
	})

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Request received for path: %s", r.URL.Path)
	//log.Println("Attempting to serve: ../frontend/index.html")
	//log.Println("Serving:", "../frontend/index.html")
	//http.ServeFile(w, r, "../frontend/index.html")
	//})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Println("Serving index.html")
			http.ServeFile(w, r, "../frontend/index.html")
		} else {
			log.Printf("Unexpected path: %s\n", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/registration.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/registration.html")
	})

	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/login.html")
	})

	http.HandleFunc("/posts.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/posts.html")
	})
	http.HandleFunc("/chat.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/chat.html")
	})

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
