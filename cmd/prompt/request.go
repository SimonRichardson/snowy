package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/trussle/snowy/pkg/models"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	greenf = color.New(color.FgGreen).SprintfFunc()
	red    = color.New(color.FgRed).SprintFunc()
	redf   = color.New(color.FgRed).SprintfFunc()
	white  = color.New(color.FgWhite).SprintFunc()
	whitef = color.New(color.FgWhite).SprintfFunc()
)

type client struct {
	base string
	out  io.Writer
}

func (c client) health(statusOnly bool) {
	endpoint := "status/health"
	if statusOnly {
		endpoint = "status/ready"
	}
	res, err := http.Get(fmt.Sprintf("%s/%s", c.base, endpoint))
	if err != nil {
		c.output(redf("Error response %q with error %s", endpoint, white(err.Error())))
		return
	}

	defer res.Body.Close()
	if c.errors(res, err) {
		return
	}

	c.output(greenf("valid %s", white(endpoint)))
}

func (c client) ledger(id string) {
	res, err := http.Get(fmt.Sprintf("%s/ledgers/?resource_id=%s", c.base, id))
	if err != nil {
		defer res.Body.Close()
	}
	if c.errors(res, err) {
		return
	}

	endpoint := res.Request.URL.String()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.output(redf("Error response %q with error %s", endpoint, white(err.Error())))
		return
	}

	var ledger models.Ledger
	if err := json.Unmarshal(body, &ledger); err != nil {
		c.output(redf("Error response %q with error %s", endpoint, white(err.Error())))
		return
	}

	c.outputLedger(ledger)
}

func (c client) ledgers(id string) {
	res, err := http.Get(fmt.Sprintf("%s/ledgers/revisions/?resource_id=%s", c.base, id))
	if err != nil {
		defer res.Body.Close()
	}
	if c.errors(res, err) {
		return
	}

	endpoint := res.Request.URL.String()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.output(redf("Error response %q with error %s", endpoint, white(err.Error())))
		return
	}

	var ledgers []models.Ledger
	if err := json.Unmarshal(body, &ledgers); err != nil {
		c.output(redf("Error response %q with error %s", endpoint, white(err.Error())))
		return
	}

	c.outputLedgers(ledgers)
}

func (c client) output(s string) {
	fmt.Fprintln(c.out, s)
}

func (c client) errors(res *http.Response, err error) bool {
	endpoint := res.Request.URL.String()
	if err != nil {
		c.output(redf("Error requesting %q with error %s", endpoint, white(err.Error())))
		return true
	}
	if code := res.StatusCode; code != 200 {
		c.output(redf("Error requesting %q with status code %s", endpoint, white(code)))
		return true
	}
	return false
}

func (c client) outputLedger(ledger models.Ledger) {
	c.outputLedgers([]models.Ledger{
		ledger,
	})
}

func (c client) outputLedgers(ledgers []models.Ledger) {
	data := [][]string{}
	for _, ledger := range ledgers {
		data = append(data, []string{
			ledger.Name(),
			ledger.ResourceID().String(),
			ledger.ResourceAddress(),
			strconv.Itoa(int(ledger.ResourceSize())),
			ledger.ResourceContentType(),
			ledger.AuthorID(),
			strings.Join(ledger.Tags(), ","),
			ledger.CreatedOn().Format(time.RFC3339),
		})
	}

	c.output(green("valid"))

	table := tablewriter.NewWriter(c.out)
	table.SetHeader([]string{
		"Name",
		"Resource ID",
		"Resource Address",
		"Resource Size",
		"Resource ContentType",
		"Author ID",
		"Tags",
		"Created On",
	})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
