package dto

type ConsultRequest struct {
    NumeroIdentificacion string `json:"numeroIdentificacion" binding:"required,min=10,max=25"`
    Token                string `json:"token" binding:"required"`
}

type ConsultResponse struct {
    NumeroIdentificacion 	string `json:"numeroIdentificacion"`
    Status               	string `json:"status"`
	Message 				string `json:"message"`
}
