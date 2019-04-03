package controllers

import (
	"hdg.com/bulebuff/overview/services"
	"hdg.com/bulebuff/overview/datamodels"
	"github.com/go-errors/errors"
	"github.com/kataras/iris"
)

type MovieController struct {
	Service services.MovieService
}

func (c *MovieController) Get() (results []datamodels.Movie) {
	return c.Service.GetAll()
}

/**
获取一部电影
 */
func (c *MovieController) GetBy(id int64) (movie datamodels.Movie, found bool) {
	return c.Service.GetByID(id)
}

func (c *MovieController) PutBy(ctx iris.Context, id int64) (datamodels.Movie, error) {
	file, info, err := ctx.FormFile("poster")
	if err != nil {
		return datamodels.Movie{}, errors.New("")
	}

	file.Close()

	poster := info.Filename

	genre := ctx.FormValue("genre")

	return c.Service.UpdatePosterAndGenreByID(id, poster, genre)
}

func (c *MovieController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)

	if wasDel {
		return iris.Map{"deleted": id}
	}

	return iris.StatusBadRequest
}
