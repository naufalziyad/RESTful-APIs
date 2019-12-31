package action

import "github.com/naufalziyad/RESTful-APIs/apis/middlewares"

func (s *Server) initializeRoutes() {

	//home
	s.Router.HandleFunc("/", middlewares.MiddlewareJSON(s.Home)).Methods("GET")

	//login
	s.Router.HandleFunc("/login", middlewares.MiddlewareJSON(s.Login)).Methods("POST")

	//user
	s.Router.HandleFunc("/users", middlewares.MiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.MiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.MiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.MiddlewareJSON(middlewares.MiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.MiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//ADS
	s.Router.HandleFunc("/ads", middlewares.MiddlewareJSON(s.CreateAds)).Methods("POST")
	s.Router.HandleFunc("/ads", middlewares.MiddlewareJSON(s.GetAdsAll)).Methods("GET")
	s.Router.HandleFunc("/ads/{id}", middlewares.MiddlewareJSON(s.GetAds)).Methods("GET")
	s.Router.HandleFunc("/ads/{id}", middlewares.MiddlewareJSON(middlewares.MiddlewareAuthentication(s.UpdateAds))).Methods("PUT")
	s.Router.HandleFunc("/ads/{id}", middlewares.MiddlewareAuthentication(s.DeleteAds)).Methods("DELETE")

}
