package controllers

import "github.com/robfig/revel"

func init() {
	revel.InterceptMethod(App.init, revel.BEFORE)
}
