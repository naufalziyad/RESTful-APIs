package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naufalziyad/RESTful-APIs/apis/auth"
	"github.com/naufalziyad/RESTful-APIs/apis/helpers/formaterror"
	"github.com/naufalziyad/RESTful-APIs/apis/helpers/responses"
	"github.com/naufalziyad/RESTful-APIs/apis/models"
)

func (server *Server) CreateAds(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	ads := models.Ads{}
	err = json.Unmarshal(body, &ads)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	ads.Prepare()
	err = ads.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != ads.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	adsCreated, err := ads.SaveAds(server.DB)
	if err != nil {
		formmatedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formmatedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, adsCreated.ID))
	responses.JSON(w, http.StatusCreated, adsCreated)
}

func (server *Server) GetAdsAll(w http.ResponseWriter, r *http.Request) {
	ads := models.Ads{}
	adsAll, err := ads.FindAllAds(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, adsAll)
}

func (server *Server) GetAds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	ads := models.Ads{}

	adsReceived, err := ads.FindAdsByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, adsReceived)
}

func (server *Server) UpdateAds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//check ads valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//check if auth token valid
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//check if ads exist
	ads := models.Ads{}
	err = server.DB.Debug().Model(models.Ads{}).Where("id = ?", pid).Take(&ads).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Ads Not Found"))
		return
	}

	if uid != ads.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//read data  ads
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//start process request data
	adsUpdate := models.Ads{}
	err = json.Unmarshal(body, &adsUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != adsUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unatuhorized"))
		return
	}

	adsUpdate.Prepare()

	err = adsUpdate.Validate()
	if uid != ads.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unatuhorized"))
		return
	}

	adsUpdate.ID = ads.ID

	adsUpdated, err := adsUpdate.UpdateAds(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, adsUpdated)
}

func (server *Server) DeleteAds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//check user authenticated
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//checkk ads exist or not
	ads := models.Ads{}
	err = server.DB.Debug().Model(models.Ads{}).Where("id = ?", pid).Take(&ads).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	if uid != ads.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = ads.DeleteAds(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entitiy", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
