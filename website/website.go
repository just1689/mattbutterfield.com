package website

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m-butterfield/mattbutterfield.com/datastore"
	_ "github.com/mattn/go-sqlite3"
)

const (
	adminPathBase     = "/admin/"
	dateDisplayLayout = "January 2006"
	DBFileName        = "app.db"
	imageBaseURL      = "http://images.mattbutterfield.com/"
	imagePathBase     = "/img/"
	port              = "8000"
	templatePath      = "website/templates/"
)

var (
	adminTemplateName = templatePath + "admin.html"
	imageStore        datastore.ImageStore
	imageTemplateName = templatePath + "image.html"
)

type imagePage struct {
	ImageCaption  string
	ImageDate     string
	ImageLocation string
	ImageURL      string
	NextImagePath string
}

func makeImagePage(image *datastore.Image, nextImageID string) imagePage {
	return imagePage{
		ImageCaption:  image.Caption,
		ImageDate:     getImageTimeStr(image),
		ImageLocation: image.Location,
		ImageURL:      imageBaseURL + image.ID,
		NextImagePath: makeImagePath(nextImageID),
	}
}

type adminPage struct {
	imagePage
	PreviousURL string
	NextURL     string
}

func newAdminPage(image *datastore.Image, prevImageID, nextImageID string) adminPage {
	var prevURL, nextURL string
	if prevImageID != "" {
		prevURL = makeAdminPath(prevImageID)
	}
	if nextImageID != "" {
		nextURL = makeAdminPath(nextImageID)
	}
	return adminPage{
		imagePage:   makeImagePage(image, ""),
		PreviousURL: prevURL,
		NextURL:     nextURL,
	}
}

func getImageTimeStr(image *datastore.Image) string {
	t, err := image.TimeFromID()
	if err != nil {
		return ""
	}
	return t.Format(dateDisplayLayout)
}

func Run() error {
	db, err := datastore.InitDB(DBFileName)
	if err != nil {
		return err
	}
	imageStore = datastore.DBImageStore{DB: db}
	fmt.Println("Serving on port: ", port)
	err = http.ListenAndServe(net.JoinHostPort("", port), buildRouter())
	if err != nil {
		return err
	}
	return nil
}

func buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc(imagePathBase+"{id}", img)
	r.HandleFunc(adminPathBase+"{id}", admin)
	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	image, err := imageStore.GetRandomImage()
	if err != nil {
		http.Error(w, "error fetching image", http.StatusInternalServerError)
	}
	http.Redirect(w, r, makeImagePath(image.ID), http.StatusFound)
}

func img(w http.ResponseWriter, r *http.Request) {
	id, err := decodeImageID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid image id", http.StatusInternalServerError)
		return
	}
	image, err := imageStore.GetImage(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "error fetching image", http.StatusInternalServerError)
		return
	}
	nextImage, err := imageStore.GetRandomImage()
	if err != nil {
		http.Error(w, "error fetching next image", http.StatusInternalServerError)
	}
	tmpl, err := template.ParseFiles(imageTemplateName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error fetching template", http.StatusInternalServerError)
		return
	}
	imagePage := makeImagePage(image, nextImage.ID)
	tmpl.Execute(w, imagePage)
}

func admin(w http.ResponseWriter, r *http.Request) {
	id, err := decodeImageID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid image id", http.StatusInternalServerError)
		return
	}
	image, err := imageStore.GetImage(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "error fetching image", http.StatusInternalServerError)
		return
	}
	previous, next, err := imageStore.GetPrevNextImages(image.ID)
	if err != nil {
		http.Error(w, "error fetching previous and next images", http.StatusInternalServerError)
	}
	tmpl, err := template.ParseFiles(adminTemplateName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error fetching template", http.StatusInternalServerError)
		return
	}
	var previousID, nextID string
	if previous != nil {
		previousID = previous.ID
	}
	if next != nil {
		nextID = next.ID
	}
	adminPage := newAdminPage(image, previousID, nextID)
	tmpl.Execute(w, adminPage)
}

func makeImagePath(imageID string) string {
	return imagePathBase + encodeImageID(imageID)
}

func makeAdminPath(imageID string) string {
	return adminPathBase + encodeImageID(imageID)
}

func decodeImageID(encodedID string) (string, error) {
	imageID, err := base64.StdEncoding.DecodeString(encodedID)
	if err != nil {
		return "", err
	}
	return string(imageID), nil
}

func encodeImageID(imageID string) string {
	return base64.StdEncoding.EncodeToString([]byte(imageID))
}
