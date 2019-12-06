//https://github.com/valyala/fasthttp
//https://gitlab.com/ykyuen/golang-echo-template-example
//https://gist.github.com/mashingan/4212d447f857cfdfbbba4f5436b779ac
/*	 "github.com/jinzhu/gorm/dialects/mysql"*/
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}

func allUsers(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var users []User
		db.Find(&users)
		fmt.Println("{}", users)

		return c.JSON(http.StatusOK, users)
	}
}

func newUser(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		db.Create(&User{Name: name, Email: email})
		return c.String(http.StatusOK, name+" user successfully created")
	}
}

func newUser2(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		db.Create(&User{Name: name, Email: email})
		return c.String(http.StatusOK, name+" user successfully created")
	}
}

func deleteUser(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")

		var user User
		db.Where("name = ?", name).Find(&user)
		db.Delete(&user)

		return c.String(http.StatusOK, name+" user successfully deleted")
	}
}

func updateUser(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		var user User
		db.Where("name=?", name).Find(&user)
		user.Email = email
		db.Save(&user)
		return c.String(http.StatusOK, name+" user successfully updated")
	}
}

func usersByPage(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))
		var result []User
		db.Limit(limit).Offset(limit * (page - 1)).Find(&result)
		return c.JSON(http.StatusOK, result)
	}
}

func handleRequest(db *gorm.DB) {
	e := echo.New()

	e.GET("/users", allUsers(db))
	e.GET("/user", usersByPage(db))
	e.POST("/user/:name/:email", newUser2(db))
	e.POST("/user", newUser(db))
	e.DELETE("/user/:name", deleteUser(db))
	e.PUT("/user/:name/:email", updateUser(db))

	e.Logger.Fatal(e.Start(":3000"))
}

func initialMigration(db *gorm.DB) {

	db.AutoMigrate(&User{})
}

func main() {
	fmt.Println("Go ORM tutorial")
	db, err := gorm.Open("mysql", "root:123456789@/crm?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	initialMigration(db)
	handleRequest(db)
}
