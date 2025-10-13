package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"githup.com/Mo277210/booking/internal/config"
)

func TestRoute(t *testing.T) {
var app config.AppConfig

mux:=routes(&app)

switch v:=mux.(type){

	case *chi.Mux:
	default:
	t.Error(fmt.Printf("type is not *chi.Mux, but is %T",v))

}

}