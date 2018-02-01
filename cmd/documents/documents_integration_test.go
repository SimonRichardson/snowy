// +build integration

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/trussle/uuid"
)

func TestStatus(t *testing.T) {
	serverURL := setupDocuments("8080")

	res, err := http.Get(fmt.Sprintf("%s/status/health", serverURL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if expected, actual := http.StatusOK, res.StatusCode; expected != actual {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}
}

type ledger struct {
	Name                string   `json:"name"`
	ResourceAddress     string   `json:"resource_address"`
	ResourceSize        int64    `json:"resource_size"`
	ResourceContentType string   `json:"resource_content_type"`
	AuthorID            string   `json:"author_id"`
	Tags                []string `json:"tags"`
}

type ledgerInput struct {
	Name     string   `json:"name"`
	AuthorID string   `json:"author_id"`
	Tags     []string `json:"tags"`
}

type ledgerOutput struct {
	ResourceID string `json:"resource_id"`
}

type contentOutput struct {
	Address     string `json:"address"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

func TestLedgerGet(t *testing.T) {
	var (
		serverURL  = setupDocuments("8081")
		ledgersURL = fmt.Sprintf("%s/ledgers/", serverURL)

		inputModel = ledgerInput{
			Name:     "ledger-name",
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(ledgersURL, "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID
	}

	get := func(resourceID string) string {
		res, err := http.Get(fmt.Sprintf("%s?resource_id=%s", ledgersURL, resourceID))
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

		var ledger ledgerInput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.Name
	}

	resourceID := post()
	ledgerName := get(resourceID)

	if expected, actual := inputModel.Name, ledgerName; expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestLedgerSelectRevisions(t *testing.T) {
	var (
		serverURL  = setupDocuments("8082")
		ledgersURL = fmt.Sprintf("%s/ledgers/", serverURL)

		inputModel = ledgerInput{
			Name:     "ledger-name",
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(ledgersURL, "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID
	}

	getRevisions := func(resourceID string) []string {
		res, err := http.Get(fmt.Sprintf("%srevisions/?resource_id=%s", ledgersURL, resourceID))
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

		var ledgers []ledgerInput
		if err := json.Unmarshal(output, &ledgers); err != nil {
			t.Fatal(err)
		}

		var names []string
		for _, v := range ledgers {
			names = append(names, v.Name)
		}

		return names
	}

	resourceID := post()
	ledgerName := getRevisions(resourceID)

	if expected, actual := inputModel.Name, ledgerName[0]; expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestLedgerAudit(t *testing.T) {
	var (
		serverURL  = setupDocuments("8083")
		ledgersURL = fmt.Sprintf("%s/ledgers/", serverURL)

		inputModel = ledgerInput{
			Name:     "ledger-name",
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	post := func() string {
		input, err := json.Marshal(inputModel)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(ledgersURL, "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID
	}

	put := func(id string, model ledgerInput) string {
		input, err := json.Marshal(model)
		if err != nil {
			t.Fatal(err)
		}

		res, err := Put(fmt.Sprintf("%s?resource_id=%s", ledgersURL, id), "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID
	}

	getRevisions := func(resourceID string) []string {
		res, err := http.Get(fmt.Sprintf("%srevisions/?resource_id=%s", ledgersURL, resourceID))
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

		var ledgers []ledgerInput
		if err := json.Unmarshal(output, &ledgers); err != nil {
			t.Fatal(err)
		}

		var names []string
		for _, v := range ledgers {
			names = append(names, v.Name)
		}

		return names
	}

	var (
		resourceID = post()
		models     = make([]ledgerInput, 10)
	)
	for k := range models {
		model := ledgerInput{
			Name:     fmt.Sprintf("ledger-name-%d", k),
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{fmt.Sprintf("tag-%d", k)},
		}
		put(resourceID, model)
		models[k] = model
	}
	ledgerName := getRevisions(resourceID)

	if expected, actual := len(models)+1, len(ledgerName); expected != actual {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}

	for k := range ledgerName {
		if k == 0 {
			if expected, actual := inputModel.Name, ledgerName[k]; expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}
			continue
		}

		if expected, actual := models[k-1].Name, ledgerName[k]; expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	}
}

func TestContentsAudit(t *testing.T) {
	var (
		serverURL   = setupDocuments("8084")
		ledgersURL  = fmt.Sprintf("%s/ledgers/", serverURL)
		contentsURL = fmt.Sprintf("%s/contents/", serverURL)

		inputModel = ledgerInput{
			Name:     "ledger-name",
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{"abc", "def", "g"},
		}
	)

	insert := func(body []byte) contentOutput {
		res, err := http.Post(contentsURL, "application/octet-stream", bytes.NewBuffer(body))
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

		var content contentOutput
		if err := json.Unmarshal(output, &content); err != nil {
			t.Fatal(err)
		}

		return content
	}

	post := func(body []byte) (string, ledger) {
		content := insert(body)

		payload := ledger{
			Name:                inputModel.Name,
			ResourceAddress:     content.Address,
			ResourceSize:        content.Size,
			ResourceContentType: content.ContentType,
			AuthorID:            inputModel.AuthorID,
			Tags:                inputModel.Tags,
		}
		input, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(ledgersURL, "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID, payload
	}

	put := func(id string, model ledgerInput, body []byte) (string, ledger) {
		content := insert(body)

		payload := ledger{
			Name:                model.Name,
			ResourceAddress:     content.Address,
			ResourceSize:        content.Size,
			ResourceContentType: content.ContentType,
			AuthorID:            model.AuthorID,
			Tags:                model.Tags,
		}
		input, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		res, err := Put(fmt.Sprintf("%s?resource_id=%s", ledgersURL, id), "application/json", bytes.NewBuffer(input))
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

		var ledger ledgerOutput
		if err := json.Unmarshal(output, &ledger); err != nil {
			t.Fatal(err)
		}

		return ledger.ResourceID, payload
	}

	getRevisions := func(resourceID string) []byte {
		res, err := http.Get(fmt.Sprintf("%srevisions/?resource_id=%s", contentsURL, resourceID))
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

		return output
	}

	content := make([]byte, 1024)
	if _, err := rand.Read(content); err != nil {
		t.Fatal(err)
	}

	var (
		resourceID, seed = post(content)
		models           = make([]ledger, 11)
		contents         = make([][]byte, 11)
	)

	models[0] = seed
	contents[0] = content

	for i := 1; i < len(models); i++ {
		model := ledgerInput{
			Name:     fmt.Sprintf("ledger-name-%d", i),
			AuthorID: uuid.MustNew().String(),
			Tags:     []string{fmt.Sprintf("tag-%d", i)},
		}
		content := make([]byte, 1024)
		if _, err := rand.Read(content); err != nil {
			t.Fatal(err)
		}
		_, result := put(resourceID, model, content)
		models[i] = result
		contents[i] = content
	}

	fileBytes := getRevisions(resourceID)

	reader, err := zip.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		t.Fatal(err)
	}

	indexOf := func(models []ledger, name string) int {
		for k, v := range models {
			if v.ResourceAddress == name {
				return k
			}
		}
		return -1
	}

	for _, v := range reader.File {
		index := indexOf(models, v.Name)
		if actual := index; actual == -1 {
			t.Fatalf("expected: >= 0, actual: %d", actual)
		}

		reader, err := v.Open()
		if err != nil {
			t.Fatal(err)
		}
		defer reader.Close()

		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := contents[index], bytes; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
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
