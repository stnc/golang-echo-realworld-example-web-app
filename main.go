package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"

	"github.com/labstack/echo/middleware"
	"github.com/siredwin/pongorenderer/renderer" // this package is not publicly available
)

var (
	paginator    = &pagination.Paginator{}
	data         = pongo2.Context{}
	MainRenderer = renderer.Renderer{Debug: true} // use any renderer
)

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

		db.Debug().Limit(postsPerPage).Order("id asc").Offset(offset).Find(&posts)

		data = pongo2.Context{"paginator": paginator, "posts": posts}
		return c.Render(http.StatusOK, "resources/index.html", data)
	}
}

func create(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		data = pongo2.Context{
			"title": "Create Post",
		}
		return c.Render(http.StatusOK, "resources/create.html", data)
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
		return c.Render(http.StatusOK, "resources/show.html", data)
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
		var id_string = strconv.FormatUint(uint64(lastID.ID), 10)

		// for flash message
		var status string
		if id != 0 {
			status = "success"
		} else {
			status = "alert"
		}
		c.Redirect(http.StatusMovedPermanently, "/post/edit/"+id_string)

		data = pongo2.Context{
			"title":        "Create",
			"flashMsg":     status,
			"flashMsgType": status,
			"dataType":     dataType,
		}
		return c.Render(http.StatusOK, "resources/create.html", data)
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
		return c.Render(http.StatusOK, "resources/edit.html", data)
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
		/* //https://forum.golangbridge.org/t/unable-to-check-nil-for-a-struct-referenced-as-pointer/9226/3
		// for flash message
		var status string
		if nil ==  postData{}  {
			status = "success"
		} else {
			status = "alert"
		}
		*/status := "success"
		c.Redirect(http.StatusMovedPermanently, "/post/edit/"+id)

		data = pongo2.Context{
			"title":        "Create",
			"id":           id,
			"flashMsg":     status,
			"flashMsgType": status,
			"dataType":     dataType,
		}
		return c.Render(http.StatusOK, "resources/edit.html", data)
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

	// Echo instance
	e := echo.New()
	e.Renderer = MainRenderer //pongo init

	// link https://echo.labstack.com/middleware/method-override
	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	e.Static("/assets", "resources/assets") //https://echo.labstack.com/guide/static-files

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialMigration(dbConn)

	// Route => handler
	e.GET("/", index(dbConn))
	e.GET("/create", create(dbConn))
	e.POST("/store", store(dbConn))

	e.GET("/post/edit/:id", edit(dbConn))

	e.GET("/post/show/:id", show(dbConn))
	e.DELETE("/post/:id", delete(dbConn))

	e.POST("/post/update/:id", update(dbConn))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
