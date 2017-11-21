package log

import (
	"log"
)

var(
	LOG logger
)

func init(){
	LOG
}

type logger log.Logger{
	
}

func (g *logger) Info(str string){
	g.SetPrefix("INFO ")
}