// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// Base /
// Version:1.0.0
//
// Consumes:
// 		-application/json
//
//	Produces:
// 		-application/json
// swagger:meta


package handlers

import(
	"log"
	"net/http"
	"microServicesNick/data"
	"github.com/gorilla/mux"
//	"regexp"
 	"strconv"
 	"context"
 	"fmt"
	
)
type Products struct {

	l *log.Logger
}

func NewProducts(l*log.Logger) *Products{
	return &Products{l}
}
/*
func  (p *Products) ServeHTTP (rw http.ResponseWriter, r *http.Request){

if r.Method == http.MethodGet {
	p.getProducts(rw, r)
	return
}

if(r.Method == http.MethodPost){
	p.addProducts(rw,r)
	return 

}
if r.Method == http.MethodPut{

	p.l.Println("PUT", r.URL.Path)
	reg:= regexp.MustCompile(`/([0-9]+)`)
	g:= reg.FindAllStringSubmatch(r.URL.Path,-1)

	if(len(g)!=1){
		p.l.Println("Invalid URI more than one id")
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	if(len(g[0]) !=2){
		p.l.Println("Invalid URI more than one capture group")
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	idString := g[0][1]
	id,err := strconv.Atoi(idString)
	if err!=nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
	}

	p.updateProduct(id,rw,r)
	p.l.Println("got id", id)
	return
}
rw.WriteHeader(http.StatusMethodNotAllowed)
	
 

 ###############delted below#########
	lp := data.GetProducts()
	err:= lp.ToJSON(rw)
	//err := data.ToJSON(lp)
 	//d, err:= json.Marshal(lp)
 	if(err!=nil){
 		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
 	}
 	#######################deleted above############
 	//rw.Write(d)
}
*/
func (p*Products) GetProducts (rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle Get Products")
	lp := data.GetProducts()
	err:= lp.ToJSON(rw)
	if(err!=nil){
 		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
 	}

}
func (p*Products) AddProducts (rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle ADD Products")
/*
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if(err != nil){
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	*/
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
	p.l.Printf("Prod:  %#v", prod)

}

func (p Products) UpdateProduct ( rw http.ResponseWriter, r *http.Request){
	 vars := mux.Vars(r)
	 id, err := strconv.Atoi(vars["id"])
	 if err != nil {
	 	http.Error(rw,"unable to convert id",http.StatusBadRequest)
	 }

	p.l.Println("Handle PUT Product",id)
	fmt.Println("going to fetch product from context")
//##################### here it wasnot accepting data.Product. look back again

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	fmt.Println("before going to update product ")
	err = data.UpdateProduct(id, prod)
	fmt.Println("updated")
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil{
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}
	type KeyProduct struct {

	}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler{

fmt.Println("inside middelware")
	return http.HandlerFunc( func (rw http.ResponseWriter, r *http.Request){
fmt.Println("inside middelware1")
	prod := &data.Product{}
	fmt.Println("inside middelware2")
	//err := prod.FromJSON(r.Body)
	err := data.FromJSON(prod, r.Body)
	fmt.Println("inside middelware3")
	p.l.Println(*prod)
	fmt.Println("inside middelware4")
		if err != nil {
			fmt.Println("inside middelware5")
		http.Error (rw, "Unable to unmarshal json", http.StatusBadRequest)
		 }

		 // validate the product
		 err = prod.Validate()
		 if err != nil {
		 	p.l.Println("[ERROR] validating product", err)
		 	http.Error(
		 		rw, 
		 		fmt.Sprintf("Error validating product: %s", err), 
		 		http.StatusBadRequest,
		 	)
		 	return
		 }
		 fmt.Println("inside middelware6")

		 ctx := context.WithValue(r.Context(),KeyProduct{}, prod)
		 fmt.Println("inside middelware7")
		 r = r.WithContext(ctx)
		 fmt.Println("inside middelware8")
		 next.ServeHTTP(rw, r)
		 fmt.Println("inside middelware9")

fmt.Println("context set done..")
	})
}


