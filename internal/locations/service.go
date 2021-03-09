package locations

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type service struct {
	repo     irepo
	log      ilogger
	response iresponse
}

type iresponse interface {
	Resp(http.ResponseWriter, int, interface{})
	Err(http.ResponseWriter, int, error)
}

type ilogger interface {
	Println(...interface{})
	Errorln(...interface{})
}

type ServiceConfig struct {
	Db       idb
	Repo     irepo
	Logger   ilogger
	Response iresponse
}

func DefaultService(cfg ServiceConfig) http.Handler {
	svc := service{
		repo:     NewRepo(cfg.Db),
		log:      cfg.Logger,
		response: cfg.Response,
	}

	r := chi.NewRouter()
	r.Post("/", svc.Insert)
	r.Get("/", svc.Select)
	r.Get("/city", svc.Select)
	return r
}

func (s service) Insert(w http.ResponseWriter, r *http.Request) {
	loc := Location{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loc)
	if err != nil {
		s.log.Errorln(err)
		s.response.Err(w, http.StatusBadRequest, err)
		return
	}

	err = loc.ValidateInsert()
	if err != nil {
		s.log.Errorln(err)
		s.response.Err(w, http.StatusBadRequest, err)
		return
	}

	err = s.repo.Save(loc)
	if err != nil {
		s.log.Errorln(err)
		s.response.Err(w, http.StatusBadRequest, err)
		return
	}

	s.response.Resp(w, http.StatusCreated, nil)

}

func (s service) Select(w http.ResponseWriter, r *http.Request) {
	locs := []Location{}
	err := s.repo.Select(r.URL.Query(), &locs)
	if err != nil {
		s.log.Errorln(err)
		s.response.Err(w, http.StatusBadRequest, err)
		return
	}

	s.response.Resp(w, http.StatusOK, locs)
}
