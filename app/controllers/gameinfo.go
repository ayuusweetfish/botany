package controllers

import (
	"fmt"
	"net/http"
)

func gamelistHandler (w http.ResponseWriter, r *http.Request) {
	
}

func init () {
	registerFunc("/gamelist", gamelistHandler);
}
