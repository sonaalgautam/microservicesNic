package main

import (
	"net/http"
	"os"
	"log"
	"time"
	"microServicesNick/handlers"
	"context"
	"os/signal"
	"github.com/gorilla/mux"
)

func main(){
	l:= log.New(os.Stdout,"product-api", log.LstdFlags)

	//hh := handlers.NewHello(l)
	//gh := handlers.NewGoodbye(l) 
	ph := handlers.NewProducts(l)

 	//sm := http.NewServeMux()
 	sm:= mux.NewRouter()
 	getRouter := sm.Methods("GET").Subrouter()
 	getRouter.HandleFunc("/", ph.GetProducts)

 	putRouter := sm.Methods(http.MethodPut).Subrouter()
 	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
 	putRouter.Use(ph.MiddlewareProductValidation)

 	postRouter := sm.Methods(http.MethodPost).Subrouter()
 	postRouter.HandleFunc("/", ph.AddProducts)
 	postRouter.Use(ph.MiddlewareProductValidation)

 	
 	//sm.Handle("/",hh)
	//sm.Handle("/goodbye", gh)
	//sm.Handle("/products", ph)
	//http.ListenAndServe(":9090", sm)
	
	//sm.Handle("/",ph)
	s := &http.Server{
		Addr:			":9090",
		Handler:		sm,
		IdleTimeout:	120 * time.Second,
		ReadTimeout:	1* time.Second,
		WriteTimeout:	1*time.Second,

	}

go func (){

	err := s.ListenAndServe()
	if(err != nil){
		l.Fatal(err)
	}

}()
	sigChan := make (chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}