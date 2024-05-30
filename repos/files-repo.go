package repos

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"

	dbModels "dev.farukh/copy-close/models/db_models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var writeFlag = os.O_CREATE | os.O_APPEND | os.O_RDWR
var readFlag = os.O_RDONLY
var permissionMode fs.FileMode = 0644

func NewFileRepo(dsn string) FileRepo {
	db, _ := openConnection(dsn)
	return &FileRepoImpl{
		Mutex:    &sync.Mutex{},
		sessions: make(map[uuid.UUID]*fileUploadMeta, 0),
		db:       db,
	}
}

type fileUploadMeta struct {
	lastReceived time.Time
	done         bool
	path         string
	file         *os.File
	uploaded     int64
	fileSize     int64
}

func (meta *fileUploadMeta) release() {
	meta.done = true
	meta.file.Close()
}

func (meta *fileUploadMeta) remove() {
	meta.file.Close()
	os.Remove(meta.path)
}

func (meta *fileUploadMeta) name() string {
	splittedPath := strings.Split(meta.path, "/")
	return splittedPath[len(splittedPath)-1]
}

type FileRepo interface {
	CreateSession(pathToFile string, filSize int64) uuid.UUID
	WriteToSession(sessionId uuid.UUID, chunk []byte) (uuid.UUID, error)
	SessionExists(sessionID uuid.UUID) bool
	GetDocument(documentID uuid.UUID, offset int64) (*os.File, error)
	GetDocumentPath(documentID uuid.UUID) (string, string, error)
}

type FileRepoImpl struct {
	*sync.Mutex
	sessions map[uuid.UUID]*fileUploadMeta
	db       *gorm.DB
}

func (repo *FileRepoImpl) InitGarbageScan() {
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			repo.Lock()
			checkTime := time.Now()

			keysToDelete := make([]uuid.UUID, 0, len(repo.sessions))
			for sessionID, meta := range repo.sessions {
				switch {
				case meta.done:
					keysToDelete = append(keysToDelete, sessionID)
					continue
				case checkTime.Sub(meta.lastReceived).Minutes() >= time.Hour.Minutes()*10:
					meta.remove()
					keysToDelete = append(keysToDelete, sessionID)
				}
			}

			for _, session := range keysToDelete {
				delete(repo.sessions, session)
			}

			repo.Unlock()
		}
	}()
}

func (repo *FileRepoImpl) CreateSession(pathToFile string, totalChunks int64) uuid.UUID {
	repo.Lock()
	defer repo.Unlock()

	println(pathToFile)
	newSessionID := uuid.NewV4()
	repo.sessions[newSessionID] = createFileMeta(pathToFile, totalChunks)

	return newSessionID
}

func (repo *FileRepoImpl) WriteToSession(sessionId uuid.UUID, chunk []byte) (uuid.UUID, error) {
	repo.Lock()
	defer repo.Unlock()

	meta := repo.sessions[sessionId]
	meta.lastReceived = time.Now()
	written, err := meta.file.Write(chunk)

	if err != nil {
		return uuid.Nil, err
	}

	meta.uploaded += int64(written)

	if	meta.uploaded == meta.fileSize {
		return repo.createDocument(meta)
	}
	return uuid.Nil, nil
}

func (repo *FileRepoImpl) SessionExists(sessionID uuid.UUID) bool {
	repo.Lock()
	defer repo.Unlock()

	_, ok := repo.sessions[sessionID]
	return ok
}

func (repo *FileRepoImpl) GetDocument(documentID uuid.UUID, offset int64) (*os.File, error) {
	var document dbModels.Document
	err := repo.db.Where("id = ?", document).Find(&document).Error
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(document.Path, readFlag, permissionMode)
	if err != nil {
		return nil, err
	}

	seeked, err := file.Seek(offset, 0)
	if err != nil {
		return nil, err
	} else if seeked != offset {
		return nil, errors.New("unable to seek to specified offset")
	}

	return file, nil
}

func (repo *FileRepoImpl) GetDocumentPath(documentID uuid.UUID) (string, string, error) {
	var document dbModels.Document
	err := repo.db.Where("id = ?", documentID).Find(&document).Error
	if err != nil {
		return "", "", err
	}
	return document.Path, document.Name, nil
}

func (repo *FileRepoImpl) createDocument(meta *fileUploadMeta) (uuid.UUID, error) {
	meta.release()
	document := dbModels.Document{
		Name: meta.name(),
		Path: meta.path,
	}
	err := repo.db.Create(&document).Error
	return document.ID, err
}

func createFileMeta(pathToFile string, fileSize int64) *fileUploadMeta {
	file, _ := os.OpenFile(pathToFile, writeFlag, permissionMode)
	return &fileUploadMeta{
		done: false,
		file: file,
		path: pathToFile,
		fileSize: fileSize,
		uploaded: 0,
	}
}
