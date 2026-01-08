package main 

import ( 
    "fmt"
    "strconv"
    "net/http"
    "html/template"
)

// posso aceder ao app logger dentro da função
func (app *application) home ( w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server","Go")

    files := [] string  {
        "./ui/html/base.tmpl",
        "./ui/html/pages/home.tmpl",
        "./ui/html/partials/nav.tmpl",
    }

    ts, err := template.ParseFiles( files... )
    if err != nil {
        app.serverError( w, r, err )
        return
    }

    err = ts.ExecuteTemplate(w,"base",nil)
    // o ts.Execute escreve o template no body, posso dar parametros de cenas
    if err != nil {
        app.serverError( w, r, err )
        return
    }

}


func (app *application) create ( w http.ResponseWriter, _ *http.Request) {
    w.Write ( [] byte ("create something"))
}

func (app *application) createPost ( w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader (http.StatusCreated) // cria status code sozinho
    w.Write([] byte ("save a new snippet"))
}

func (app *application) view ( w http.ResponseWriter, r *http.Request) {
    id,err := strconv.Atoi( r.PathValue("id") ) 
    if err != nil || id < 1 {
        http.NotFound(w,r)
        return 
    }   
    message := fmt.Sprintf("display something =  %d", id)
    w.Write ( [] byte (message) )
}

