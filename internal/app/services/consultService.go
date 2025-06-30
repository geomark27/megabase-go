package services

import (
	"encoding/json"
	"fmt"
	"io"
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/logger"
	"megabaseGo/internal/models"
	"megabaseGo/internal/database"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type ConsultService struct{}

func NewConsultService() *ConsultService {
	return &ConsultService{}
}

func (s *ConsultService) GetCitizenByNumeroIdentificacion(req *dto.ConsultRequest) (interface{}, error) {
	id := req.NumeroIdentificacion
	length := len(id)
	logger.Debug.WithFields(logrus.Fields{"id": id, "length": length}).Debug("Iniciando validación de identificación")

	if length != 10 && length != 13 {
		msg := "El número debe tener 10 o 13 dígitos"
		logger.Debug.WithFields(logrus.Fields{"id": id}).Warn(msg)
		return &dto.ConsultResponse{NumeroIdentificacion: id, Status: "invalid", Message: msg}, nil
	}

	idType := "cedula"
	if length == 13 {
		idType = "ruc"
		suffix := id[10:]
		logger.Debug.WithFields(logrus.Fields{"id": id, "suffix": suffix}).Debug("Validando sufijo de RUC")
		if suffix != "001" {
			msg := "Los últimos 3 dígitos del RUC son inválidos"
			logger.Debug.WithFields(logrus.Fields{"id": id}).Warn(msg)
			return &dto.ConsultResponse{NumeroIdentificacion: id, Status: "invalid", Message: msg}, nil
		}
	}

	logger.Debug.WithFields(logrus.Fields{"id": id}).Info("Consultando API externa")
	raw, err := s.fetchExternalData(idType, id)
	if err != nil {
		return nil, err
	}

	if err := s.saveOrUpdateDB(idType, raw); err != nil {
		logger.Debug.WithError(err).Error("Error guardando o actualizando en la base de datos")
	}

	return raw, nil
}

func (s *ConsultService) fetchExternalData(idType, id string) (interface{}, error) {
	baseURL		:= os.Getenv("API_URL")
	apiKey 		:= os.Getenv("API_KEY")
	endpoint	:= fmt.Sprintf("%s/%s/%s", baseURL, idType, id)

	logger.Debug.WithField("endpoint", endpoint).Info("Construyendo petición HTTP")
	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Debug.WithError(err).Error("Error creando la petición HTTP")
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Debug.WithError(err).Error("Error ejecutando petición HTTP")
		return nil, err
	}
	defer resp.Body.Close()

	logger.Debug.WithField("status", resp.StatusCode).Info("Respuesta HTTP obtenida")

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("API externa status %d: %s", resp.StatusCode, string(body))
		logger.Debug.WithFields(logrus.Fields{"status": resp.StatusCode}).Error(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Debug.WithError(err).Error("Error decodificando JSON")
		return nil, err
	}

	logger.Debug.WithField("data", result).Info("JSON decodificado exitosamente")
	return result, nil
}

func (s *ConsultService) saveOrUpdateDB(idType string, raw interface{}) error {
	siType := map[string]string{"cedula": "05", "ruc": "04"}[idType]

	payload, ok := raw.(map[string]interface{})["resultado"]
	if !ok {
		return fmt.Errorf("formato inesperado: missing 'resultado'")
	}
	data, ok := payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("formato inesperado de 'resultado'")
	}

	cit := models.Citizen{
		NumeroIdentificacion:        dataValue(data, "NumeroRuc", "Cedula"),
		TipoIdentificacion:          siType,
		Email:                       dataValue(data, "Email"),
		Celular:                     dataValue(data, "Celular"),
		Convencional:                dataValue(data, "Convencional"),
		DireccionPrincipal:          dataValue(data, "DireccionContribuyente", "Domicilio"),
		Pais:                        "ECUADOR",
		Provincia:                   nestedString(data, "DPA_DireccionContribuyente", "Provincia"),
		Ciudad:                      nestedString(data, "DPA_DireccionContribuyente", "Canton"),
		TipoContribuyente:           dataValue(data, "TipoContribuyente"),
		EstadoContribuyente:         dataValue(data, "EstadoContribuyente"),
		Regimen:                     dataValue(data, "Regimen"),
		Categoria:                   dataValue(data, "Categoria"),
		ObligadoContabilidad:        dataValue(data, "ObligadoContabilidad"),
		AgenteRetencion:             ptrString(dataValue(data, "AgenteRetencion")),
		ContribuyenteEspecial:       ptrString(dataValue(data, "ContribuyenteEspecial")),
		ActividadEconomicaPrincipal: dataValue(data, "ActividadEconomicaPrincipal"),
		MotivoCancelacionSuspension: dataValue(data, "MotivoCancelacionSuspension"),
	}

	if idType == "ruc" {
		cit.RazonSocial = ptrString(dataValue(data, "RazonSocial"))
		cit.NombreComercial = ptrString(dataValue(data, "NombreComercial"))
		rep, _ := json.Marshal(data["RepresentantesLegales"])
		cit.RepresentantesLegales = datatypes.JSON(rep)
		suc, _ := json.Marshal(data["Sucursales"])
		cit.Sucursales = datatypes.JSON(suc)
	}

	if idType == "cedula" {
		cit.Nombre = ptrString(dataValue(data, "NombreCiudadano"))
		cit.Genero = ptrString(dataValue(data, "Sexo"))
		cit.EstadoCivil = ptrString(dataValue(data, "EstadoCivil"))
		if dob := dataValue(data, "FechaNacimiento"); dob != "" {
			parsed, err := time.Parse("02/01/2006", dob)
			if err == nil {
				cit.FechaNacimiento = &parsed
			}
		}
		cit.Nacionalidad = ptrString(dataValue(data, "Nacionalidad"))
	}

	var existing models.Citizen
	err := database.DB.Where("numero_identificacion = ?", cit.NumeroIdentificacion).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		if err := database.DB.Create(&cit).Error; err != nil {
			return err
		}
		logger.Debug.WithFields(logrus.Fields{"citizen_id": cit.ID}).Info("Citizen creado con éxito")
	} else {
		cit.ID = existing.ID
		if err := database.DB.Save(&cit).Error; err != nil {
			return err
		}
		logger.Debug.WithFields(logrus.Fields{"citizen_id": cit.ID}).Info("Citizen actualizado con éxito")
	}

	return nil
}

func dataValue(data map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := data[k]; ok {
			if s, ok2 := v.(string); ok2 {
				return s
			}
		}
	}
	return ""
}

func nestedString(data map[string]interface{}, objKey, field string) string {
	if nested, ok := data[objKey]; ok {
		if m, ok2 := nested.(map[string]interface{}); ok2 {
			if v, ok3 := m[field]; ok3 {
				if s, ok4 := v.(string); ok4 {
					return s
				}
			}
		}
	}
	return ""
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
