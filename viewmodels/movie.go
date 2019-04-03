package viewmodels

import (
	"hdg.com/bulebuff/overview/datamodels"
)

type Movie struct {
	datamodels.Movie
}

func (m Movie) IsValid() bool {
	return m.ID > 0
}

//func (m Movie) Dispatch(ctx context.Context){
//	if !m.IsValid() {
//		ctx.NotFound()
//		return
//	}
//	ctx.JSON(m,ctx.JSON{})
//}
