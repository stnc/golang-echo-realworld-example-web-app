package main

/*
license : www.selmantunc.com.tr

*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stnc/pongo-renderer-echo4/renderer"
)

var (
	paginator = &pagination.Paginator{}
	data      = pongo2.Context{}
)

//Posts struct
type Posts struct {
	gorm.Model
	// ID      int    `gorm:"type:int(11); NULL;index"`
	Title   string `gorm:"type:varchar(255);" json:"title_"`
	Content string `gorm:"type:text" json:"Content"  `
}

func index(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {

		posts := []Posts{} // a slice

		var total int

		db.Find(&posts).Count(&total)
		postsPerPage := 5
		paginator = pagination.NewPaginator(c.Request(), postsPerPage, total)
		offset := paginator.Offset()

		db.Debug().Limit(postsPerPage).Order("id desc").Offset(offset).Find(&posts)

		data = pongo2.Context{"paginator": paginator, "posts": posts}
		return c.Render(http.StatusOK, "resources/blog/index.html", data)
	}
}

func create(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		data = pongo2.Context{
			"title": "Create Post",
		}
		return c.Render(http.StatusOK, "resources/blog/create.html", data)
	}
}

func show(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		dataType := c.Request().Method
		post := Posts{}
		db.Debug().Where("id = ?", id).Take(&post)
		data = pongo2.Context{
			"title":    "Post Edit",
			"post":     post,
			"id":       id,
			"dataType": dataType,
		}
		return c.Render(http.StatusOK, "resources/blog/show.html", data)
	}
}

func store(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {

		dataType := c.Request().Method

		title := c.FormValue("Title")
		Content := c.FormValue("Content")

		//posts := &Posts{title, description, time.Now(), time.Now(), nil}

		lastID := Posts{
			// CreatedAt: time.Now(),
			// UpdatedAt: time.Now(),
			// DeletedAt: nil,
			Title:   title,
			Content: Content,
		}
		db.Debug().Create(&lastID)

		var id = lastID.ID
		var idString = strconv.FormatUint(uint64(lastID.ID), 10)

		// for flash message
		var status string
		if id != 0 {
			status = "success"
		} else {
			status = "alert"
		}
		c.Redirect(http.StatusMovedPermanently, "/post/edit/"+idString)

		data = pongo2.Context{
			"title":        "Create",
			"flashMsg":     status,
			"flashMsgType": status,
			"dataType":     dataType,
		}
		return c.Render(http.StatusOK, "resources/blog/create.html", data)
	}
}

func edit(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		dataType := c.Request().Method
		post := Posts{}
		db.Debug().Where("id = ?", id).Take(&post)
		data = pongo2.Context{
			"title":    "Post Edit",
			"post":     post,
			"id":       id,
			"dataType": dataType,
		}
		return c.Render(http.StatusOK, "resources/blog/edit.html", data)
	}
}

func update(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {

		id := c.Param("id")

		dataType := c.Request().Method

		title := c.FormValue("Title")
		Content := c.FormValue("Content")

		//posts := &Posts{title, description, time.Now(), time.Now(), nil}

		postData := Posts{
			// CreatedAt: time.Now(),
			// UpdatedAt: time.Now(),
			// DeletedAt: nil,
			Title:   title,
			Content: Content,
		}
		db.Debug().Model(postData).Where("id = ?", id).Update(postData)
		fmt.Println(postData)

		status := "success"
		c.Redirect(http.StatusMovedPermanently, "/post/edit/"+id)

		data = pongo2.Context{
			"title":        "Create",
			"id":           id,
			"flashMsg":     status,
			"flashMsgType": status,
			"dataType":     dataType,
		}
		return c.Render(http.StatusOK, "resources/blog/edit.html", data)
	}
}

func delete(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		var post Posts
		db.Where("id = ?", id).Take(&post)
		db.Delete(&post)
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}

//InitialMigration migratein init
func initialMigration(db *gorm.DB) {

	db.AutoMigrate(&Posts{})
	//dummy data
	for i := 0; i < 30; i++ {
		post := Posts{Title: "hello" + strconv.Itoa(i), Content: "lorem ipsumm lorem ipsummlorem ipsummlorem ipsummlorem ipsummlorem ipsumm" + strconv.Itoa(i)}
		db.Create(&post)
	}
}

func main() {

	dbConn, err := gorm.Open("sqlite3", "blog.db")
	defer dbConn.Close()
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer dbConn.Close()

	// sql init
	dbConn.DB()
	dbConn.DB().Ping()
	dbConn.DB().SetMaxIdleConns(10)
	dbConn.DB().SetMaxOpenConns(100)

	mainRenderer := renderer.Renderer{Debug: true}

	server := echo.New()

	server.Renderer = mainRenderer

	// Middleware

	server.Use(middleware.Logger())

	server.Use(middleware.Recover())

	server.Use(session.Middleware(sessions.NewCookieStore([]byte("KoronaSecret"))))

	//method
	server.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	server.Static("/assets", "resources/assets") //https://echo.labstack.com/guide/static-files

	//migration run
	// initialMigration(dbConn)

	// Route => handler
	server.GET("/", index(dbConn))
	server.GET("/create", create(dbConn))
	server.POST("/store", store(dbConn))

	server.GET("/post/edit/:id", edit(dbConn))

	server.GET("/post/show/:id", show(dbConn))
	server.DELETE("/post/:id", delete(dbConn))

	server.POST("/post/update/:id", update(dbConn))

	// Start server
	server.Logger.Fatal(server.Start(":8000"))
}
