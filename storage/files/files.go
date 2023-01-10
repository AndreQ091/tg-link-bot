package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	e "tg-link-bot/lib/error"
	"tg-link-bot/storage"
	"time"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = e.WrapIfErr("can't save page", err)
	}()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)

	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)

	if err != nil {
		return err
	}

	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil

}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		err = e.WrapIfErr("can't pick random page", err)
	}()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))

}

func (s *Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)

	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	fPath := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fPath)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)

	if err != nil {
		return false, e.Wrap("can't check file", err)
	}

	fPath := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file exists %s", fPath)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
