/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/raihankhan/jwt-auth-system/cmd"
	"github.com/raihankhan/jwt-auth-system/pkg/config"
)

func main() {
	config.Init()
	cmd.Execute()
}
