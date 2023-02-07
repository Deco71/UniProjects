package api

import (
	"encoding/json"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"wasaPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}

	img, err := png.Decode(r.Body)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusUnsupportedMediaType, "Received an unsupported format")
		return
	}

	id, err := rt.db.NewImage(name)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in NewImage function", err)
		return
	}

	f, err := os.Create("./service/database/dbfiles/" + id + ".png")

	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in creating the image file", err)
		_ = rt.db.DeleteImage(id)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in encoding the image", err)
		_ = rt.db.DeleteImage(id)
		return
	}

	var idInt int
	idInt, _ = strconv.Atoi(id)
	_ = json.NewEncoder(w).Encode(Image{Image: idInt})
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "application/json")

	name, err := rt.AutorizeToken(r, w)
	if err != nil {
		return
	}

	id := ps.ByName("img")

	name2, err := rt.db.CheckImage(id)
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Image Not Found")
		return
	}
	if name != name2 {
		rt.HttpErrCodeSender(w, http.StatusForbidden, "You are trying to change other users info. Authenticate with the right account and try again")
		return
	}

	err = os.Remove("./service/database/dbfiles/" + id + ".png")
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in deleting the image file", err)
		return
	}

	err = rt.db.DeleteImage(id)
	if err != nil {
		rt.ErrLoggerAndSender(w, ctx, "Error in deleting the image from the database", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("content-type", "image/png")
	id := ps.ByName("img")

	buf, err := os.ReadFile("./service/database/dbfiles/" + id + ".png")
	if err != nil {
		rt.HttpErrCodeSender(w, http.StatusNotFound, "Image Not Found")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", `attachment;filename="image.png"`)
	_, _ = w.Write(buf)
}
