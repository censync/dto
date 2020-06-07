Package for conversation different request forms to Service layer DTO (Data Transfer Object). 

Structure tags (default `dto`) available for linking fields with different names in source or destination structure, also _snake_case_ tags will be assigned to _CamelCase_. 

For example:
```go
type SignUpStep1Form struct {
	Firstname string `json:"firstname" validate:"attr=firstname,min=1,max=50" dto:"name"`
	Lastname  string `json:"lastname" validate:"attr=lastname,max=50"`
}

type SignUpStep2Form struct {
	Username  string `json:"username" validate:"attr=username,min=3,max=32"`
	Email     string `json:"email" validate:"attr=email,email" dto:"contact_email"`
}

```

```go
type SignUpDTO struct {	
	Name             string
	Lastname         string
    Path             string `dto:"username"`
	ContactEmail     string  
}
```

```go
	formDTO := SignUpDTO{}
	err := dto.RequestToDTO(&formDTO, &requestStep1, &requestStep2)
```


Controller
```go
package controllers

import (
	"github.com/censync/go-api-structure/service"
	"github.com/censync/go-dto"
	"github.com/censync/go-validator"
	"github.com/gin-gonic/gin"
)

// Input JSON form with validation rules
type SignUpForm struct {
	Username  string `json:"username" validate:"attr=username,min=3,max=32"`
	Firstname string `json:"firstname" validate:"attr=firstname,min=1,max=50"`
	Lastname  string `json:"lastname" validate:"attr=lastname,max=50"`
	Email     string `json:"email" validate:"attr=email,email"`
}

type SignUpDTO struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
}

func TestMethod(ctx *gin.Context) {
	request := SignUpForm{}
	
    ctx.BindJSON(&request)

	errs := validator.Validate(&request)

	if !errs.IsEmpty() {
		ctx.JSON(400, errs)
		return
	}

	formDTO := SignUpDTO{}

	err := dto.RequestToDTO(&formDTO, &request)

	if err != nil {
		ctx.JSON(400, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = service.AppService().SignUp(&formDTO)

	if err != nil {
		ctx.JSON(400, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(200)
	return
}
```

Service
```go
func (s *Service) SignUp(stx *cs.ContextService) (dto *SignUpDTO) error {
	// Process user sign up
}
```