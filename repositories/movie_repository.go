package repositories

import (
	"hdg.com/bulebuff/overview/datamodels"
	"sync"
	"errors"
)

type Query func(movie datamodels.Movie) bool

type MODE int

const (
	READ_ONLY_MODE  MODE = iota
	READ_WRITE_MODE
)

type MovieRepository interface {
	Exec(query Query, action Query, limit int, mode MODE) (ok bool)
	Select(query Query) (movie datamodels.Movie, found bool)
	SelectMany(query Query, limit int) (results []datamodels.Movie)
	InsertOrUpdate(movie datamodels.Movie) (updateMovie datamodels.Movie, err error)
	Delete(query Query, limit int) (daleted bool)
}

type movieMemoryRepository struct {
	MovieRepository
	source map[int64]datamodels.Movie
	mu     sync.RWMutex
}

func NewMovieRepository(source map[int64]datamodels.Movie) MovieRepository {
	return &movieMemoryRepository{
		source: source,
	}
}

func (rep *movieMemoryRepository) Exec(query Query, action Query, limit int, mode MODE) (ok bool) {
	loops := 0
	if mode == READ_ONLY_MODE {
		rep.mu.RLock()
		defer rep.mu.RUnlock()
	} else {
		rep.mu.Lock()
		defer rep.mu.Unlock()
	}

	for _, movie := range rep.source {
		if ok = query(movie); ok && action(movie) {
			loops++
			if limit >= loops {
				break
			}
		}
	}
	return
}

func (rep *movieMemoryRepository) Select(query Query) (movie datamodels.Movie, found bool) {
	if found = rep.Exec(query, func(m datamodels.Movie) bool {
		movie = m
		return true
	}, 1, READ_ONLY_MODE); !found {
		movie = datamodels.Movie{}
	}
	return
}

func (rep *movieMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.Movie) {
	rep.Exec(query, func(m datamodels.Movie) bool {
		results = append(results, m)
		return true
	}, limit, READ_ONLY_MODE)
	return
}

func (rep *movieMemoryRepository) InsertOrUpdate(movie datamodels.Movie) (updateMovie datamodels.Movie, err error) {
	id := movie.ID
	if id == 0 {
		var lastID int64

		rep.mu.RLock()
		for _, item := range rep.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		rep.mu.RUnlock()

		id = lastID + 1
		movie.ID = id

		rep.mu.Lock()
		rep.source[id] = movie
		rep.mu.Unlock()
		return movie, nil
	}
	current, exists := rep.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})
	if !exists {
		return datamodels.Movie{}, errors.New("failed to update a nonexistent movie")
	}
	if movie.Poster != "" {
		current.Poster = movie.Poster
	}
	if movie.Genre != "" {
		current.Genre = movie.Genre
	}

	rep.mu.Lock()
	rep.source[id] = current
	rep.mu.Unlock()

	return movie, nil
}

func (rep *movieMemoryRepository) Delete(query Query, limit int) (daleted bool) {
	return rep.Exec(query, func(m datamodels.Movie) bool {
		delete(rep.source, m.ID)
		return true
	}, limit, READ_WRITE_MODE)
}
