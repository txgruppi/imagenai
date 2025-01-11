package api

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type FileUpload interface {
	Name() string
	Open() (io.ReadSeekCloser, error)
}

func NewFileUploadFromLocalFS(filepath string) FileUpload {
	return &localFileUpload{filepath: filepath}
}

type localFileUpload struct {
	filepath string
}

func (t *localFileUpload) Name() string {
	return path.Base(t.filepath)
}

func (t *localFileUpload) Open() (io.ReadSeekCloser, error) {
	return os.Open(t.filepath)
}

type getTemporaryUploadLinksRequestFile struct {
	Name string `json:"file_name"`
	MD5  string `json:"md5"`
}

type getTemporaryUploadLinksRequest struct {
	Files []getTemporaryUploadLinksRequestFile `json:"files_list"`
}

type getTemporaryUploadLinksResponse struct {
	Data struct {
		Files []struct {
			Name string `json:"file_name"`
			Link string `json:"upload_link"`
		} `json:"files_list"`
	} `json:"data"`
}

func (t *client) prepareFile(file FileUpload, buffer *bytes.Buffer) error {
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file.Name(), err)
	}
	defer reader.Close()

	_, err = io.Copy(buffer, reader)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
	}

	return nil
}

func (t *client) getBase64EncodedMD5Hash(md5 hash.Hash, buffer *bytes.Buffer) string {
	md5.Write(buffer.Bytes())
	return base64.StdEncoding.EncodeToString(md5.Sum(nil))
}

func (t *client) getTemporaryUploadLink(project Project, file FileUpload, hash string, encoded *bytes.Buffer) (string, error) {
	var data getTemporaryUploadLinksRequest
	data.Files = append(data.Files, getTemporaryUploadLinksRequestFile{
		Name: file.Name(),
		MD5:  hash,
	})

	if err := json.NewEncoder(encoded).Encode(data); err != nil {
		return "", fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := t.newPostRequest(fmt.Sprintf("/projects/%s/get_temporary_upload_links", project.ID), encoded)
	if err != nil {
		return "", fmt.Errorf("failed to create request to get temporary upload link: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get temporary upload link: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get temporary upload link, expected HTTP 200: %s", resp.Status)
	}

	var envelope getTemporaryUploadLinksResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", fmt.Errorf("failed to decode temporary upload link: %w", err)
	}

	if len(envelope.Data.Files) != 1 {
		return "", fmt.Errorf("expected 1 file in response, got %d", len(envelope.Data.Files))
	}

	return envelope.Data.Files[0].Link, nil
}

func (t *client) uploadFile(file FileUpload, hash, link string, buffer *bytes.Buffer) error {
	// we use http.NewRequest here because this request is sent to
	// AWS S3, and not to the API server
	req, err := http.NewRequest(http.MethodPut, link, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request to upload file %s: %w", file.Name(), err)
	}
	req.Header.Set("Content-Type", "")
	req.Header.Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	req.Header.Set("Content-MD5", hash)

	resp, err := t.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %w", file.Name(), err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file %s, expected HTTP 200: %s", file.Name(), resp.Status)
	}

	return nil
}

func (t *client) UploadFiles(project Project, files []FileUpload) error {
	var prepared bytes.Buffer
	var encoded bytes.Buffer
	md5 := md5.New()
	for _, file := range files {
		prepared.Reset()
		md5.Reset()
		encoded.Reset()

		if err := t.prepareFile(file, &prepared); err != nil {
			return err
		}

		hash := t.getBase64EncodedMD5Hash(md5, &prepared)

		link, err := t.getTemporaryUploadLink(project, file, hash, &encoded)
		if err != nil {
			return err
		}

		err = t.uploadFile(file, hash, link, &prepared)
		if err != nil {
			return err
		}
	}
	return nil
}
