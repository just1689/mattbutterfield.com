package datastore

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	InsertImageRegex       = "^INSERT INTO images \\(id, caption\\) VALUES \\(\\?, \\?\\)$"
	SelectImageByIDRegex   = "^SELECT id, caption FROM images WHERE id = \\?$"
	SelectLatestImageRegex = "^SELECT id, caption FROM images ORDER BY id DESC LIMIT 1$"
	SelectRandomImageRegex = "^SELECT id, caption FROM images WHERE id = \\(SELECT id FROM images ORDER BY RANDOM\\(\\) LIMIT 1\\)$"
)

func TestGetImage(t *testing.T) {
	db, db_mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	id, caption := "1234", "hello"
	db_mock.ExpectQuery(SelectImageByIDRegex).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "caption"}).AddRow(id, caption))

	store := DBImageStore{DB: db}
	image, err := store.GetImage(id)
	if err != nil {
		t.Fatal(err)
	}
	if err := db_mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if image.ID != id {
		t.Errorf("Unexpected image id: %s != %s", id, image.ID)
	}
	if image.Caption != caption {
		t.Errorf("Unexpected image caption: %s != %s", caption, image.Caption)
	}
}

func TestGetLatestImage(t *testing.T) {
	db, db_mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	db_mock.ExpectQuery(SelectLatestImageRegex).
		WillReturnRows(sqlmock.NewRows([]string{"id", "caption"}).AddRow("1234", ""))

	store := DBImageStore{DB: db}
	_, err = store.GetLatestImage()
	if err != nil {
		t.Fatal(err)
	}
	if err := db_mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetRandomImage(t *testing.T) {
	db, db_mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	db_mock.ExpectQuery(SelectRandomImageRegex).
		WillReturnRows(sqlmock.NewRows([]string{"id", "caption"}).AddRow("1234", nil))

	store := DBImageStore{db}
	_, err = store.GetRandomImage()
	if err != nil {
		t.Fatal(err)
	}
	if err := db_mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSaveImage(t *testing.T) {
	db, db_mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	id, caption := "1234", "hello"
	db_mock.ExpectExec(InsertImageRegex).WithArgs(id, caption).WillReturnResult(sqlmock.NewResult(1, 1))

	image := Image{ID: id, Caption: caption}
	store := DBImageStore{DB: db}
	err = store.SaveImage(image)
	if err != nil {
		t.Fatal(err)
	}
	if err := db_mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSaveImageNilCaption(t *testing.T) {
	db, db_mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	id := "1234"
	db_mock.ExpectExec(InsertImageRegex).WithArgs(id, nil).WillReturnResult(sqlmock.NewResult(1, 1))

	image := Image{ID: id}
	store := DBImageStore{DB: db}
	err = store.SaveImage(image)
	if err != nil {
		t.Fatal(err)
	}
	if err := db_mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
