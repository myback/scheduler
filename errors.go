package scheduler

import "fmt"

var ErrorJobAlreadyExists = fmt.Errorf("the job with the same name already exists")
var ErrorJobNotFound = fmt.Errorf("job not found")
