package main


import (
	"github.com/astaxie/beego/utils/pagination"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo/middleware"
	"github.com/siredwin/pongorenderer/renderer" // this package is not publicly available
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	paginator = &pagination.Paginator{}
	data = pongo2.Context{}
	MainRenderer = renderer.Renderer{Debug:true} // use any renderer
)


type Posts struct {
	gorm.Model
	// ID       int    `gorm:"type:int(11); NULL;index"`
	Title    string  `gorm:"type:varchar(255);" json:"title_"`
	Content string `gorm:"type:text" json:"Content"  `
}


func index(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {

		posts := []Posts{} // a slice

		var total int

		db.Find(&posts).Count(&total)
		postsPerPage := 10
		paginator = pagination.NewPaginator(c.Request(), postsPerPage, total)
		offset := paginator.Offset()

		db.Debug().Limit(postsPerPage).Order("id asc").Offset(offset).Find(&posts)

		data = pongo2.Context{"paginator": paginator, "posts": posts}
		return c.Render(http.StatusOK, "resources/create.html", data)
	}
}

//InitialMigration migratein init
func initialMigration(db *gorm.DB) {

	db.AutoMigrate(&Posts{})
		//dummy data 
	for i := 0; i < 30; i++ {
		post := Posts{Title: "hello"+ strconv.Itoa(i), Content: "lorem ipsumm lorem ipsummlorem ipsummlorem ipsummlorem ipsummlorem ipsumm"+strconv.Itoa(i)}
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

	e.Static("/assets", "resources/assets") //https://echo.labstack.com/guide/static-files
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	initialMigration(dbConn)

	// Route => handler
	e.GET("/", index(dbConn))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
