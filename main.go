package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/middleware/recover"
	"hdg.com/bulebuff/overview/web/middleware"
	"hdg.com/bulebuff/overview/repositories"
	"hdg.com/bulebuff/overview/datasource"
	"hdg.com/bulebuff/overview/services"
	"hdg.com/bulebuff/overview/web/controllers"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	//加载模板文件
	app.RegisterView(iris.HTML("./web/views", ".html").Reload(true))

	//注册控制器
	//mvc.New(app.Party("/movies")).Handle(new(controllers.MovieController))
	//可以拆分编写的代码配置mvc.Application
	mvc.Configure(app.Party("/movies"), movies)
	//http://localhost:8080/movies
	//http://localhost:8080/movies/1
	app.Run(
		//开启web服务
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		//实现更快的json序列化和更多优化
		iris.WithOptimizations,
	)
}

func movies(app *mvc.Application){
	//添加基本身份验证，用于基于/电影的请求
	app.Router.Use(middleware.BasicAuth)
	//使用数据源中的一些内存，创建我们的电影资源库
	rep := repositories.NewMovieRepository(datasource.Movies)
	//创建我们的电影服务，我们将他绑定到电影应用程序的依赖项
	movieService := services.NewMovieService(rep)

	app.Register(movieService)

	app.Handle(new(controllers.MovieController))
}