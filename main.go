package main

import (
	"database/sql"
	"net/http"
	"os"

	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/peterlavreniuk/gotodo/docs"
	gtd "github.com/peterlavreniuk/gotodo/src"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var Config gtd.Config
var Log = logrus.New()

func init() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Out = os.Stdout
		Log.Info("Unable to create log file")
	}

	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.DebugLevel)
}

// @title GOTODO Api
// @version 0.1
// @Accept json
// @Produce json
// @description This is a simple typical todo list Go project
// @host 127.0.0.1:8080
// @BasePath /
func main() {
	Config, err := gtd.ReadConfig()
	if err != nil {
		panic(err)
	}

	MigrateTable(Config, Log)

	toDoController := gtd.ToDoController{Config: *Config}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(ErrorHandler(Log))
	router.Use(gin.LoggerWithWriter(Log.Writer()))

	router.GET("/", toDoController.GetAll)
	router.GET("/:id", toDoController.Get)
	router.DELETE("/:id", toDoController.Delete)
	router.POST("/", toDoController.Create)
	router.PUT("/:id", toDoController.Update)

	url := ginSwagger.URL("http://127.0.0.1:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run()
}

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			logger.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
		}
	}
}

func MigrateTable(config *gtd.Config, logger *logrus.Logger) {
	connStr := fmt.Sprintf("%s:%s@/%s", config.MySql.UserName, config.MySql.Password, config.MySql.DatabaseName)
	db, err := sql.Open("mysql", connStr)

	if err != nil {
		logger.Error(err)
		return
	}

	defer db.Close()

	query := `SET sql_notes = 0;  
CREATE TABLE IF NOT EXISTS ` + "`" + `todoitem` + "`" + `(
` + "`" + `id` + "`" + `mediumint NOT NULL AUTO_INCREMENT,
` + "`" + `title` + "`" + `varchar(128) NOT NULL,
` + "`" + `description` + "`" + `varchar(512) DEFAULT NULL,
PRIMARY KEY (` + "`" + `id` + "`" + `)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
SET sql_notes = 1;`

	_, err = db.Exec(query)
	if err != nil {
		logger.Error(err)
		return
	}
}
