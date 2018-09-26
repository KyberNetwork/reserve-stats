package users

import (
	"log"

	"go.uber.org/zap"
)

func configLog(stdoutLog bool) {
	_, err := zap.NewProduction()
	if err != nil {
		log.Printf("Cannot init zap logger")
	}
}

func main() {

}
