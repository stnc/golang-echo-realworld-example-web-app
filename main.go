package main


import (
	"github.com/astaxie/beego/utils/pagination"
	"github.com/labstack/echo"
	"net/http"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo/middleware"
	"github.com/siredwin/pongorenderer/renderer" // this package is not publicly available
)

var (
	paginator = &pagination.Paginator{}
	data = pongo2.Context{}
	MainRenderer = renderer.Renderer{Debug:true} // use any renderer
)

//generator
func NewSlice(start, count, step int) []int {
	s := make([]int, count)
	for i := range s {
		s[i] = start
		start += step
	}
	return s
}

func ListAllUsers(c echo.Context) (error){
	// Lets use the Forbes top 7.
	usernames := []string{
		"Larry Ellison",
		"Carlos Slim Helu", 
		"Mark Zuckerberg",
	 	"Amancio Ortega ", 
		 "Jeff Bezos", 
		 " Warren Buffet ",
		  "Bill Gates",
		  "selman tunç",
		  "murat ohdadssd",
		  "john yedfd",
		  "lorem ipsum",
		  "lorem ipsum2",
		  "lorem ipsum3",
		  "lorem ipsum4",
		  "lorem ipsum5",
		  "lorem ipsum6",
		  "lorem ipsum7",
		  "lorem ipsum8",
		}

	// sets paginator with the current offset (from the url query param)
	postsPerPage := 2
	paginator = pagination.NewPaginator(c.Request(), postsPerPage, len(usernames))

	// fetch the next posts "postsPerPage"
	idrange := NewSlice(paginator.Offset(), postsPerPage, 1)

	//create a new page list that shows up on html
	myusernames := []string{}
	for _, num := range idrange {
		//Prevent index out of range errors
		if num <= len(usernames)-1{
			myuser := usernames[num]
			myusernames = append(myusernames, myuser)
		}
	}

	// set the paginator in context
	// also set the page list in context
	// if you also have more data, set it context
	data = pongo2.Context{"paginator":paginator, "posts":myusernames}

	return c.Render(http.StatusOK, "resources/index.html", data)
}

func main() {
	// Echo instance
	e := echo.New()
	e.Renderer = MainRenderer //pongo init

	e.Static("/assets", "resources/assets") //https://echo.labstack.com/guide/static-files
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", ListAllUsers)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
