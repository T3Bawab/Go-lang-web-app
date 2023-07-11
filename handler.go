package handlers

import (
	"MyWeb/models"
	"MyWeb/pkg/configs"
	"MyWeb/pkg/dbdriver"
	"MyWeb/pkg/forms"
	"MyWeb/pkg/renders"
	"MyWeb/pkg/repository"
	"MyWeb/pkg/repository/dbrepo"
	"log"
	"net/http"
)

type Repository struct {
	App      *configs.AppConfig
	DataBase repository.DataBaseRepo
}

var Repo *Repository // when we call Repo we can control the Repository esey

func NewRepository(appConfig *configs.AppConfig, db *dbdriver.DataBase) *Repository {
	return &Repository{
		App:      appConfig,
		DataBase: dbrepo.NewPostgresRepo(db.SQL, appConfig),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {

	var articlesList models.ArticleList

	articlesList, err := m.DataBase.Get3AnArticle()
	if err != nil {
		log.Println(err)
		return
	}

	m.App.Session.Put(r.Context(), "userid", "awab")

	data := make(map[string]interface{})

	data["articleList"] = articlesList

	renders.RenderTemplate(w, r, "home.page.tmpl", &models.DataPage{
		Data: data,
	})

}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {

	strMapData := make(map[string]string)

	renders.RenderTemplate(w, r, "about.page.tmpl", &models.DataPage{StrMap: strMapData})

}

func (m *Repository) MakePostHandler(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	var emptyArticle models.Article

	data := make(map[string]interface{})

	data["article"] = emptyArticle

	renders.RenderTemplate(w, r, "make-post.page.tmpl",
		&models.DataPage{
			Data: data,
			Form: forms.NewForm(nil),
		})
}

// Handler for posting articles useing post
func (m *Repository) PostMakePostHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println("( ! ) Error with Parse Form : ", err)
		return
	}

	uID := (m.App.Session.Get(r.Context(), "user_id")).(int)

	article := models.Post{
		Title:   r.Form.Get("blog_title"),
		Content: r.Form.Get("blog_article"),
		UserId:  int(uID),
	}

	form := forms.NewForm(r.PostForm)

	form.HasRequired("blog_title", "blog_article")

	form.MinLength("blog_title", 5, r)
	form.MinLength("blog_article", 5, r)

	if !form.Valid() { // if the Valid Not True
		data := make(map[string]interface{})

		data["article"] = article

		renders.RenderTemplate(w, r, "make-post.page.tmpl", &models.DataPage{Form: form, Data: data})
		return
	}

	// write to the db

	err = m.DataBase.InsertPost(article)

	if err != nil {
		log.Fatal(err)
	}

	m.App.Session.Put(r.Context(), "article", article)

	http.Redirect(w, r, "/article-received", http.StatusSeeOther)

}

func (m *Repository) ArticleReceived(w http.ResponseWriter, r *http.Request) {
	article, ok := m.App.Session.Get(r.Context(), "article").(models.Article)

	if !ok {

		m.App.Session.Put(r.Context(), "error", "Can't get data from session")

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})

	data["article"] = article

	renders.RenderTemplate(w, r, "article-received.page.tmpl", &models.DataPage{
		Data: data,
	})
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	strMapData := make(map[string]string)

	renders.RenderTemplate(w, r, "login.page.tmpl", &models.DataPage{StrMap: strMapData})
}

func (m *Repository) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.NewForm(r.PostForm)

	form.HasRequired("email", "password")

	form.IsEmail("email")

	if !form.Valid() {

		renders.RenderTemplate(w, r, "login.page.tmpl", &models.DataPage{Form: form})
		return
	}

	id, _, err := m.DataBase.AuthenticateUser(email, password)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid email or password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Valid Login")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository) PageHandler(w http.ResponseWriter, r *http.Request) {
	strMapData := make(map[string]string)
	renders.RenderTemplate(w, r, "page.page.tmpl", &models.DataPage{StrMap: strMapData})
}

func (m *Repository) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())

	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
