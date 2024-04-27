package seed

import (
	"bufio"
	"fmt"
	"os"
	taskServices "pompom/go/src/services"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			return true
		} else if response == "n" {
			return false
		}
	}
}

func Seed(db *sqlx.DB) {
	taskService := taskServices.NewTaskService(db)
	seed := NewSeedService(taskService, db)
	switch os.Args[1] {
	case "seed":
		var val int
		if len(os.Args) > 2 {
			var err error
			val, err = strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("Error converting argument to integer:", err)
				val = 5
			}
		} else {
			val = 5 // Default to 5
		}

		// seed.CreateTags(val)
		seed.CreateTasks(val)
	case "deleteTasks":
		seed.DeleteTasks()
	case "deleteBase":
		if askForConfirmation("Are you sure you want to drop all tables in the database? This action cannot be undone") {
			fmt.Println("Proceeding with dropping all tables...")
			seed.DeleteBd()
		} else {
			fmt.Println("Operation canceled.")
		}
	case "createDatabase":
		seed.CreateDatabase()
	default:
		fmt.Println("unsupported arg")
	}
}
