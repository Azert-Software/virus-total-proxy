package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// SendFormData sends a file and list of parameters as a multi part form
func SendFormData(url string, file string, params map[string]string, data interface{}) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer

	// add a file if one has been passed
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		fw, err := w.CreateFormFile("file", file)
		if err != nil {
			return err
		}

		if _, err = io.Copy(fw, f); err != nil {
			return err
		}
	}

	// Add the other fields
	for key, value := range params {
		var err error

		if fw, err = w.CreateFormField(key); err != nil {
			return err
		}

		if _, err = fw.Write([]byte(value)); err != nil {
			return err
		}

	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)

	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	buf := bytes.NewBuffer(make([]byte, 0, res.ContentLength))
	_, readErr := buf.ReadFrom(res.Body)

	if readErr != nil {
		return readErr
	}

	json.Unmarshal(buf.Bytes(), data)

	return nil
}
