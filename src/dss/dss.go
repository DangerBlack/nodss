package dss

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Constants to be used in the search
const (
	UserAgent                      = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0"
	LocalizedNoResultForSearchTerm = "Nessun risultato trovato per i termini di ricerca"
	ResultStatsDivIDPattern        = `<div id="result\-stats">([a-zA-Z 0-9\.]*)<nobr>`
)

// Dss is the class containing methods for counting events
type Dss struct{}

// CountEvent fetches the Google search result for a query and returns the count of results
func (d *Dss) CountEvent(query string, after *time.Time) (int, error) {

	if after != nil {
		query = fmt.Sprintf("%s after:%s", query, after.Format("2006/01/02"))
	}

	searchKey := url.QueryEscape(`"` + strings.ReplaceAll(strings.TrimSpace(query), " ", "+") + `"`)
	url := fmt.Sprintf("https://www.google.it/search?q=%s", searchKey)

	// Create HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("user-agent", UserAgent)
	req.Close = true

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Check if result stats div is present
	bodyStr := string(body)
	if !strings.Contains(bodyStr, `<div id="result-stats">`) {
		println(bodyStr)
		return 0, nil
	}

	// Match the result count using regex
	re := regexp.MustCompile(ResultStatsDivIDPattern)
	matched := re.FindStringSubmatch(bodyStr)
	if len(matched) == 0 {
		return 0, nil
	}

	// Check if no result localized string is present
	if strings.Contains(bodyStr, LocalizedNoResultForSearchTerm) {
		return 0, nil
	}

	// Parse the number from the matched text
	fields := strings.Fields(matched[1])

	pos := 0

	if len(fields) > 2 {
		pos = 1
	}

	value := strings.ReplaceAll(fields[pos], ".", "") // Remove dots from the number
	count, err := strconv.Atoi(value)
	if err != nil {
		println("Error parsing count", value)
		return 0, err
	}

	return count, nil
}

// CountEvents returns the number of results for each of the provided search keys combined with a head
func (d *Dss) CountEvents(head string, keys []string, after *time.Time) (map[string]int, error) {
	results := make(map[string]int)
	var wg sync.WaitGroup
	mu := &sync.Mutex{} // Mutex to prevent concurrent map writes

	// Goroutine for each query to make it concurrent
	for _, key := range keys {
		println("Start query for", head, " : ", key)
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			query := fmt.Sprintf("%s %s", head, k)
			println("Querying for", query)
			count, err := d.CountEvent(query, after)
			if err != nil {
				fmt.Printf("Error for query %s: %v\n", query, err)
				return
			}
			println("Count for", query, " : ", count)
			mu.Lock()
			results[k] = count
			mu.Unlock()
		}(key)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return results, nil
}
