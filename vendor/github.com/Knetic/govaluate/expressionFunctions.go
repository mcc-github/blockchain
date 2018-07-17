package govaluate


type ExpressionFunction func(arguments ...interface{}) (interface{}, error)
