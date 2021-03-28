package server

import (
	"encoding/json"
	"go_basics/packages/remoteblog/models"
	"html/template"
	"net/http"
	"path"

	"github.com/go-chi/chi"
)

// IndexHadle returns index page with max 20 posts.
// @Summary serve index templated page
// @Description returns index page with max 20 posts
// @Success 200 {string} string
// @Failure 500 {string} models.ErrorModel
// @Produce html
// @Router / [get]
func (s *Server) IndexHadle(w http.ResponseWriter, r *http.Request) {
	posts, err := models.AllPosts(s.ctx, s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	s.Page.Data = posts
	templ := template.Must(
		template.New("page").ParseFiles(
			path.Join(s.config.Root, s.config.TemplatesDir, "index.tpl"),
			path.Join(s.config.Root, s.config.TemplatesDir, "common.tpl"),
		))
	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// ViewHandle returns page with one post by id.
// @Summary serve view post templated page
// @Description returns page with one post by id
// @Param id path string true "id"
// @Success 200 {string} string
// @Failure 500 {string} models.ErrorModel
// @Produce html
// @Router /view/{id} [get]
func (s *Server) ViewHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := models.GetPost(s.ctx, s.db, id)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	s.Page.Data = post
	templ := template.Must(
		template.New("page").ParseFiles(
			path.Join(s.config.Root, s.config.TemplatesDir, "view.tpl"),
			path.Join(s.config.Root, s.config.TemplatesDir, "common.tpl"),
		))
	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// NewHandle returns page with new post form.
// @Summary serve new post templated page
// @Description returns page with new post form
// @Success 200 {string} string
// @Failure 500 {string} models.ErrorModel
// @Produce html
// @Router /new [get]
func (s *Server) NewHandle(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(
		template.New("page").ParseFiles(
			path.Join(s.config.Root, s.config.TemplatesDir, "new.tpl"),
			path.Join(s.config.Root, s.config.TemplatesDir, "common.tpl"),
		))
	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// EditHandle returns page with edit post form.
// @Summary serve edit post templated page
// @Description returns page with edit post form
// @Param id path string true "id"
// @Success 200 {string} string
// @Failure 500 {string} models.ErrorModel
// @Produce html
// @Router /edit/{id} [get]
func (s *Server) EditHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := models.GetPost(s.ctx, s.db, id)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	s.Page.Data = post
	templ := template.Must(
		template.New("page").ParseFiles(
			path.Join(s.config.Root, s.config.TemplatesDir, "edit.tpl"),
			path.Join(s.config.Root, s.config.TemplatesDir, "common.tpl"),
		))
	if err := templ.Execute(w, s.Page); err != nil {
		s.SendInternalErr(w, err)
		return
	}
}

// PostHandle create new post.
// @Summary create new post
// @Description create new post and return json
// @Success 200 {string} models.Post
// @Failure 500 {string} models.ErrorModel
// @Produce json
// @Param post body models.Post true "title, author, content"
// @Router /api/v1/posts [post]
func (s *Server) PostHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	err = post.Insert(s.ctx, s.db)
	if err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// DeleteHandle deletes a post.
// @Summary deletes a post by id
// @Description deletes a post by id
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Param id path string true "id"
// @Param post body models.Post true "id"
// @Router /api/v1/posts/{id} [delete]
func (s *Server) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	if _, err := post.Delete(s.ctx, s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PutHandle renew post.
// @Summary renew post by id
// @Description renew post by id
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Param id path string true "id"
// @Param post body models.Post true "title, author, content"
// @Router /api/v1/posts/{id} [put]
func (s *Server) PutHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	post := &models.Post{}
	if err = json.NewDecoder(r.Body).Decode(post); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	if _, err = post.Update(s.ctx, s.db); err != nil {
		s.SendInternalErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// SwaggerHandle serve path to json for swagger.
// @Summary returns json docs
// @Description Returns swagger's json docs
// @Produce  json
// @Success 200 {string} string "ok"
// @Router /api/v1/docs/swagger.json [get]
func (s *Server) SwaggerHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.config.SwagJSON)
}
