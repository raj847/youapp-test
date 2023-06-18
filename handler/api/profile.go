package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"youapp/entity"
	"youapp/service"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ProfileAPI struct {
	profileService *service.ProfileService
	minioClient    *minio.Client
}

func NewProfileAPI(
	profileService *service.ProfileService,
	minioClient *minio.Client,
) *ProfileAPI {
	return &ProfileAPI{
		profileService: profileService,
		minioClient:    minioClient,
	}
}

func (i *ProfileAPI) AddProfile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get form file from request
	file, header, err := r.FormFile("foto")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file.Close()

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header.Filename, file, header.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	fileName1 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header.Filename)

	// get form value from request
	displayname := r.FormValue("display_name")
	gender := r.FormValue("gender")
	birthday := r.FormValue("birthday")
	horoscope := r.FormValue("horoscope")
	zodiac := r.FormValue("zodiac")
	height := r.FormValue("height")
	weight := r.FormValue("weight")

	h, _ := strconv.Atoi(height)
	we, _ := strconv.Atoi(weight)

	id := r.Context().Value("id").(uint)

	profile := entity.ProfileReq{
		Displayname: displayname,
		Gender:      gender,
		Birthday:    birthday,
		Horoscope:   horoscope,
		Zodiac:      zodiac,
		Height:      h,
		Weight:      we,
		PhotoURL1:   fileName1,
	}

	err = r.ParseMultipartForm(4096) // parsing request dengan size maksimal 4096 bytes
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, _ := time.Parse("2006-01-02", profile.Birthday)
	res, err := i.profileService.AddProfile(r.Context(), entity.Profile{
		UserID:      id,
		Displayname: displayname,
		Gender:      gender,
		Birthday:    t,
		Horoscope:   horoscope,
		Zodiac:      zodiac,
		Height:      h,
		Weight:      we,
		PhotoURL1:   fileName1,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"profileid": res.ID,
		"message":   "add profile success",
	}

	WriteJSON(w, http.StatusCreated, response)
}
func (i *ProfileAPI) GetAllProfile(w http.ResponseWriter, r *http.Request) {

	inv := r.URL.Query()
	invID, foundJenisId := inv["profile_id"]
	id := r.Context().Value("id").(uint)
	if foundJenisId {
		jID, _ := strconv.Atoi(invID[0])
		invByID, err := i.profileService.GetProfileByID(r.Context(), jID)
		if err != nil {
			if invByID.ID == 0 {
				WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error profile not found"))
				return
			}
			if invByID.UserID != id {
				WriteJSON(w, http.StatusUnauthorized, entity.NewErrorResponse("error unauthorized user id"))
			}

			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		WriteJSON(w, http.StatusOK, invByID)
		return
	}

	list, err := i.profileService.GetAllProfile(r.Context(), id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (i *ProfileAPI) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	invID := r.URL.Query().Get("profile_id")
	jID, _ := strconv.Atoi(invID)
	err := i.profileService.DeleteProfile(r.Context(), jID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"inv_id":  jID,
		"message": "success delete profile",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (i *ProfileAPI) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get form file from request
	file, header, err := r.FormFile("foto")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file.Close()

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header.Filename, file, header.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	fileName1 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header.Filename)

	// get form value from request
	displayname := r.FormValue("display_name")
	gender := r.FormValue("gender")
	birthday := r.FormValue("birthday")
	horoscope := r.FormValue("horoscope")
	zodiac := r.FormValue("zodiac")
	height := r.FormValue("height")
	weight := r.FormValue("weight")

	ids := r.Context().Value("id").(uint)

	id := r.URL.Query().Get("profile_id")
	idInt, _ := strconv.Atoi(id)

	h, _ := strconv.Atoi(height)
	we, _ := strconv.Atoi(weight)

	profile := entity.ProfileReq{
		Displayname: displayname,
		Gender:      gender,
		Birthday:    birthday,
		Horoscope:   horoscope,
		Zodiac:      zodiac,
		Height:      h,
		Weight:      we,
		PhotoURL1:   fileName1,
	}

	err = r.ParseMultipartForm(4096) // parsing request dengan size maksimal 4096 bytes
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, _ := time.Parse("2006-01-02", profile.Birthday)

	_, err = i.profileService.UpdateProfile(r.Context(), entity.Profile{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		UserID:      ids,
		Displayname: profile.Displayname,
		Gender:      profile.Gender,
		Birthday:    t,
		Horoscope:   profile.Horoscope,
		Zodiac:      profile.Zodiac,
		Height:      profile.Height,
		Weight:      profile.Weight,
		PhotoURL1:   profile.PhotoURL1,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"id":      idInt,
		"message": "update profile success",
	}

	WriteJSON(w, http.StatusCreated, response)
}
