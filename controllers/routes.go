package controllers

func (s *Server) SetUpRoutes() {
	s.Router.HandleFunc("/signup", s.Signup).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/user", s.User).Methods("GET")
	s.Router.HandleFunc("/logout", s.Logout).Methods("POST")
}