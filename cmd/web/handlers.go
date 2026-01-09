package main 

import ( 
    "fmt"
    "errors"
    "strconv"
    "net/http"
    //"html/template"

    "modulo.porreiro/internal/models"
)

// posso aceder ao app logger dentro da função
func (app *application) home ( w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server","Go")

    snippets , err := app.snippets.Latest()

    if err != nil {
        app.serverError(w,r,err)
        return
    }
    for _, snippet := range snippets {
        fmt.Fprintf(w,"%+v\n\nAKAKAKAK\n",snippet)
    }

    //files := [] string  {
    //    "./ui/html/base.tmpl",
    //    "./ui/html/pages/home.tmpl",
    //    "./ui/html/partials/nav.tmpl",
    //}

    //ts, err := template.ParseFiles( files... )
    //if err != nil {
    //    app.serverError( w, r, err )
    //    return
    //}

    //err = ts.ExecuteTemplate(w,"base",nil)
    //// o ts.Execute escreve o template no body, posso dar parametros de cenas
    //if err != nil {
    //    app.serverError( w, r, err )
    //    return
    //}

}


func (app *application) create ( w http.ResponseWriter, _ *http.Request) {
    //ss,_ := app.snippets.Latest( )
    //app.logger.Info(ss)

    w.Write ( [] byte ("create something"))
}

func (app *application) createPost ( w http.ResponseWriter, r *http.Request) {
    t :=  "AKAKKA"
    c := "AKJHDGFAKJBLAUKBGAJGFAJKGFAKJSDF&%$%&/()(/#&%&/(/&54345678"
    e := 7

    id,err := app.snippets.Insert(t,c,e)
    if err != nil {
        app.serverError(w,r,err)
        return
    }

    newPath := fmt.Sprintf("/view/%d",id)
    http.Redirect(w,r, newPath , http.StatusSeeOther)
}

func (app *application) view ( w http.ResponseWriter, r *http.Request) {
    id,err := strconv.Atoi( r.PathValue("id") ) 
    if err != nil || id < 1 {
        http.NotFound(w,r)
        return 
    }   
    
    snippet, err := app.snippets.Get( id )
    if err != nil {
        if errors.Is( err, models.ErrNoRecord ) {
            http.NotFound(w,r)
        } else {
            app.serverError(w, r, err)
        }
        return 
    }

    message := fmt.Sprintf("id = %d\n\n %+v", id, snippet)
    w.Write ( [] byte (message) )
}

