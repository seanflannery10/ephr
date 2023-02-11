package main

// func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Name     string `json:"name"`
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//
//	err := helpers.ReadJSON(w, r, &input)
//	if err != nil {
//		httperrors.BadRequest(w, r, err)
//		return
//	}
//
//	params := queries.CreateUserParams{
//		Name:      input.Name,
//		Email:     input.Email,
//		Activated: false,
//	}
//
//	v := validator.New()
//}
