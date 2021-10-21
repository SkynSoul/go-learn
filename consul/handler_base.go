package consul

import (
	"github.com/SkynSoul/go-learn/utils"
	"log"
	"net/http"
)

var WorkingPath = ""

func init() {
	var err error
	WorkingPath, err = utils.GetWorkingPath()
	if err != nil {
		log.Println("get cur working path failed, err: ", err)
		WorkingPath = "."
	}
}

func HandlerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error"))
}
