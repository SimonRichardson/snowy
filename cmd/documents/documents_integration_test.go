// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/trussle/snowy/pkg/uuid"
)

func TestStatus(t *testing.T) {
	serverURL := setupDocuments("8080")

	res, err := http.Get(fmt.Sprintf("%s/status/", serverURL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}
}

type documentInput struct {
	Name     string   `json:"name"`
	AuthorID string   `json:"author_id"`
	Tags     []string `json:"tags"`
}

type documentOutput struct {
	ResourceID string `json:"resource_id"`
}

func TestDocumentGet(t *testing.T) {
	var (
		serverURL    = setupDocuments("8081")
		documentsURL = fmt.Sprintf("%s/documents/", serverURL)

		inputModel = documentInput{
			Name:     "document-name",
			AuthorID: uuid.New().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(documentsURL, "application/json", bytes.NewBuffer(input))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var document documentOutput
		if err := json.Unmarshal(output, &document); err != nil {
			t.Fatal(err)
		}

		return document.ResourceID
	}

	get := func(resourceID string) string {
		res, err := http.Get(fmt.Sprintf("%s?resource_id=%s", documentsURL, resourceID))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var document documentInput
		if err := json.Unmarshal(output, &document); err != nil {
			t.Fatal(err)
		}

		return document.Name
	}

	resourceID := post()
	documentName := get(resourceID)

	if expected, actual := inputModel.Name, documentName; expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestDocumentGetMultiple(t *testing.T) {
	var (
		serverURL    = setupDocuments("8082")
		documentsURL = fmt.Sprintf("%s/documents/", serverURL)

		inputModel = documentInput{
			Name:     "document-name",
			AuthorID: uuid.New().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(documentsURL, "application/json", bytes.NewBuffer(input))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var document documentOutput
		if err := json.Unmarshal(output, &document); err != nil {
			t.Fatal(err)
		}

		return document.ResourceID
	}

	getMultiple := func(resourceID string) []string {
		res, err := http.Get(fmt.Sprintf("%smultiple/?resource_id=%s", documentsURL, resourceID))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var documents []documentInput
		if err := json.Unmarshal(output, &documents); err != nil {
			t.Fatal(err)
		}

		var names []string
		for _, v := range documents {
			names = append(names, v.Name)
		}

		return names
	}

	resourceID := post()
	documentNames := getMultiple(resourceID)

	if expected, actual := inputModel.Name, documentNames[0]; expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestDocumentAudit(t *testing.T) {
	var (
		serverURL    = setupDocuments("8083")
		documentsURL = fmt.Sprintf("%s/documents/", serverURL)

		inputModel = documentInput{
			Name:     "document-name",
			AuthorID: uuid.New().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(documentsURL, "application/json", bytes.NewBuffer(input))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var document documentOutput
		if err := json.Unmarshal(output, &document); err != nil {
			t.Fatal(err)
		}

		return document.ResourceID
	}

	put := func(id string, model documentInput) string {
		input, err := json.Marshal(model)
		if err != nil {
			t.Fatal(err)
		}

		res, err := Put(fmt.Sprintf("%s?resource_id=%s", documentsURL, id), "application/json", bytes.NewBuffer(input))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var document documentOutput
		if err := json.Unmarshal(output, &document); err != nil {
			t.Fatal(err)
		}

		return document.ResourceID
	}

	getMultiple := func(resourceID string) []string {
		res, err := http.Get(fmt.Sprintf("%smultiple/?resource_id=%s", documentsURL, resourceID))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
			t.Fatalf("expected: %d, actual: %d", expected, actual)
		}

		output, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var documents []documentInput
		if err := json.Unmarshal(output, &documents); err != nil {
			t.Fatal(err)
		}

		var names []string
		for _, v := range documents {
			names = append(names, v.Name)
		}

		return names
	}

	var (
		resourceID = post()
		models     = make([]documentInput, 10)
	)
	for k := range models {
		model := documentInput{
			Name:     fmt.Sprintf("document-name-%d", k),
			AuthorID: uuid.New().String(),
			Tags:     []string{fmt.Sprintf("tag-%d", k)},
		}
		put(resourceID, model)
		models[k] = model
	}
	documentNames := getMultiple(resourceID)

	if expected, actual := len(models)+1, len(documentNames); expected != actual {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}

	for k := range documentNames {
		if k == 0 {
			if expected, actual := inputModel.Name, documentNames[k]; expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}
			continue
		}

		if expected, actual := models[k-1].Name, documentNames[k]; expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	}
}

func setupDocuments(port string) string {
	var (
		wg          sync.WaitGroup
		serverURL   = fmt.Sprintf("0.0.0.0:%s", port)
		virtualised = []string{
			"-filesystem=virtual",
			"-persistence=virtual",
			"-metrics.registration=false",
			fmt.Sprintf("-api=tcp://%s", serverURL),
		}
	)

	wg.Add(1)

	go func() {
		go func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done()
		}()
		if err := runDocuments(virtualised); err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	return fmt.Sprintf("http://%s", serverURL)
}

func Put(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}
