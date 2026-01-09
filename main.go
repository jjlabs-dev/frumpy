package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

// TODO TRANNY CREATES and exits i guess
var repository, nexus_url, nexus_username, nexus_password, base_path, filter string
var filterRe *regexp.Regexp

type Asset struct {
	DownloadUrl      string            `json:"downloadUrl"`
	Path             string            `json:"path"`
	Checksum         map[string]string `json:"checksum"`
	ContentType      string            `json:"contentType"`
	LastModified     string            `json:"lastModified"`
	LastModifiedTime time.Time         `json:"-"`
	LastDownloaded   string            `json:"lastDownloaded"`
	Uploader         string            `json:"uploader"`
	FileSize         int               `json:"fileSize"`
}
type Folder struct {
	Name         string
	Size         int64
	LastModified time.Time
}
type Component struct {
	Id      string  `json:"id"`
	Group   string  `json:"group"`
	Version string  `json:"version"`
	Name    string  `json:"name"`
	Format  string  `json:"format"`
	Assets  []Asset `json:"assets"`
}

func normalizeSpacing(lines []string, delim string) string {
	columnLenghts := []int{}
	retLine := ""
	if len(lines) == 1 {
		return strings.ReplaceAll(lines[0], delim, "    ")
	}
	for _, line := range lines {
		for i, column := range strings.Split(line, delim) {
			if len(columnLenghts) > i {
				columnLenghts[i] = max(columnLenghts[i], len(column))
			} else {
				columnLenghts = append(columnLenghts, len(column))
			}
		}
	}
	for _, line := range lines {
		addLine := ""
		for i, column := range strings.Split(line, delim) {
			re := regexp.MustCompile(`Ě|Š|Č|Ř|Ž|Ý|Á|Í|É|Ó|Ú|Ů|Ť|Ď|Ň|ě|š|č|ř|ž|ý|á|í|é|ó|ú|ů|ť|ď|ň`)
			countSpec := 4
			countSpec += len(re.FindAll([]byte(column), -1))
			re = regexp.MustCompile(`–`)
			countSpec += 2 * len(re.FindAll([]byte(column), -1))
			addLine += column + strings.Repeat(" ", columnLenghts[i]-len(column)+countSpec)
		}
		retLine += addLine + "\n"

	}
	return retLine
}
func (c Component) String() string {
	ret := ""
	for _, a := range c.Assets {
		ret += fmt.Sprintf("%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%d KB", color.HiYellowString(c.Name), color.MagentaString(c.Format), color.RedString(c.Version), color.HiBlueString(a.DownloadUrl), color.CyanString(a.Uploader), color.HiRedString(a.Checksum["sha1"]), color.BlackString(a.LastModified), color.HiBlackString(a.LastDownloaded), a.FileSize/1024)
	}
	return ret
}

type ComponentList struct {
	Items             []Component `json:"items"`
	ContinuationToken *string     `json:"continuationToken"`
}

func fetchCompontents(continuationToken string) ([]Component, error) {
	var comps, retcomps []Component
	var cl ComponentList
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/service/rest/v1/components?continuationToken=%s&repository=%s", nexus_url, continuationToken, repository), nil)
	if err != nil {
		return nil, err
	}
	if len(nexus_username) > 0 && len(nexus_password) > 0 {
		req.SetBasicAuth(nexus_username, nexus_password)
	}
	req.Header.Set("Accept", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &cl)
	if err != nil {
		return nil, err
	}
	comps = cl.Items
	if cl.ContinuationToken != nil {
		tc, err := fetchCompontents(*cl.ContinuationToken)
		if err != nil {
			return nil, err
		}
		comps = append(comps, tc...)
	}

	for _, v := range comps {
		matchedAsset := false
		for _, n := range v.Assets {

			if filterRe.Match([]byte(n.DownloadUrl)) {
				matchedAsset = true
				break
			}
		}
		if filterRe.Match([]byte(v.Name)) || matchedAsset {
			v.Assets = func(a []Asset) []Asset {
				for u := range a {
					x := a[u]
					x.LastModifiedTime, err = time.Parse("2006-01-02T15:04:05.000Z07:00", x.LastModified)

					if err != nil {
						fmt.Println(err)
						return a
					}
					a[u] = x
				}
				return a
			}(v.Assets)
			split := strings.Split(v.Name, "/")
			v.Name = strings.TrimSuffix(split[len(split)-1], ".zip")
			retcomps = append(retcomps, v)
		}
	}
	return retcomps, nil
}
func main() {
	repository = os.Getenv("FRUMPY_REPOSITORY")
	nexus_url = os.Getenv("FRUMPY_URL")
	nexus_username = os.Getenv("FRUMPY_USERNAME")
	nexus_password = os.Getenv("FRUMPY_PASSWORD")
	filter = os.Getenv("FRUMPY_FILTER")
	if len(os.Args) > 1 {
		repository = os.Args[1]
	}
	if len(os.Args) > 2 {
		filter = os.Args[2]
	}
	var err error
	filterRe, err = regexp.Compile(filter)
	if err != nil {
		log.Fatalln(err)
	}
	comps, err := fetchCompontents("")
	if err != nil {
		log.Fatalln(err)
	}
	strs := []string{fmt.Sprintf("%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s#DELIM#%s", color.WhiteString("Name"), color.WhiteString("Format"), color.WhiteString("Version"), color.WhiteString("DownloadUrl"), color.WhiteString("Uploader"), color.WhiteString("SHA1"), color.WhiteString("LastModified"), color.WhiteString("LastDownloaded"), color.WhiteString("FileSize"))}
	for _, v := range comps {
		strs = append(strs, v.String())
	}
	fmt.Println(normalizeSpacing(strs, "#DELIM#"))
}
