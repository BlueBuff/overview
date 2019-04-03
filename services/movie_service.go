package services

import (
	"hdg.com/bulebuff/overview/datamodels"
	"hdg.com/bulebuff/overview/repositories"
)

type MovieService interface {
	GetAll() []datamodels.Movie
	GetByID(id int64) (datamodels.Movie, bool)
	DeleteByID(id int64) bool
	UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error)
}

type movieService struct {
	MovieService
	rep repositories.MovieRepository
}

func NewMovieService(rep repositories.MovieRepository) MovieService {
	return &movieService{rep: rep}
}

func (s *movieService) GetAll() []datamodels.Movie {
	return s.rep.SelectMany(func(datamodels.Movie) bool {
		return true
	}, -1)
}

func (s *movieService) GetByID(id int64) (datamodels.Movie, bool) {
	return s.rep.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})
}

func (s *movieService) UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error) {
	return s.rep.InsertOrUpdate(datamodels.Movie{
		ID:     id,
		Poster: poster,
		Genre:  genre,
	})
}

func (s *movieService) DeleteByID(id int64) bool {
	return s.rep.Delete(func(m datamodels.Movie) bool {
		return m.ID == id
	},1)
}


