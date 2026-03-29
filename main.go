/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"decaffeinated/internal/dtime"
	"decaffeinated/internal/dwatchdog"
	"decaffeinated/pkg/measures"
	"log"
	"time"
)
func main(){
	rules := []dwatchdog.Rule{
    	{
    	    AppName:   "alacritty",
    	    TimeLimit: 1 * time.Minute,
    	    Timestamps: []dtime.CallbackTimestamp{
    	        {Timestamp: dtime.InitTimestamp, Callback: func() { measures.Notification("alacritty", time.Now())}},
    	        {Timestamp: dtime.EndTimestamp, Callback: func() { log.Println("10% restantes! Salve o jogo.") }},
    	    },
    	},
    	{
    	    AppName: "firefox",
    	    IsBlocked: true, // Bloqueio direto
    	},
	}
	dog := dwatchdog.NewWatchDog(rules)
	dog.Start()

	select{}
}
