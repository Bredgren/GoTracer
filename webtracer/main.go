package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/Bredgren/gotracer/webtracer/lib"
	"github.com/nfnt/resize"
)

const path = "/src/github.com/Bredgren/gotracer/webtracer"
const imgType = "jpg"
const thumbnailSize = 128

var goPath = os.Getenv("GOPATH")

var templ *template.Template

var (
	debug      bool
	port       int
	maxHistory int
)

func init() {
	flag.BoolVar(&debug, "D", false, "Debug mode. More logs, use unminified assets, etc.")
	flag.IntVar(&port, "p", 8080, "Set http port")
	flag.IntVar(&maxHistory, "history", 10, "Maximum render history to keep track of")
}

func main() {
	setup()

	http.HandleFunc("/", httpHandler)
	h := &handler{}
	http.HandleFunc("/render", h.renderHandler)
	http.HandleFunc("/history", h.historyHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func setup() {
	e := os.Chdir(filepath.Join(goPath, path))
	if e != nil {
		log.Fatal(e)
	}
	flag.Parse()

	if debug {
		log.Println("Debug mode enabled!")
	}

	templ = template.Must(template.New("templ").Funcs(template.FuncMap{
		"debug":      func() bool { return debug },
		"formatTime": func(t time.Time) string { return t.Format(time.Stamp) },
	}).ParseFiles("./tmpl/page.tmpl"))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("host:", r.Host) // TODO: restrict host?
	if r.RequestURI != "/" && !strings.Contains(r.RequestURI, "?") {
		log.Println("serve file", r.RequestURI)
		http.ServeFile(w, r, "./"+r.RequestURI)
		return
	}
	log.Println("handle", r.RequestURI)

	defer func() {
		if e := r.Body.Close(); e != nil {
			log.Printf("Error closing body: %v\n", e)
		}
	}()

	initItem := &lib.RenderItem{}

	initScene := r.FormValue("initial")
	if initScene != "" {
		allHistory, e := getAllHistory()
		if e != nil {
			msg := fmt.Sprintf("Error getting all history: %v", e)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		found := false
		for _, h := range allHistory {
			if h.Scene == initScene {
				initItem = h
				found = true
				break
			}
		}
		if !found {
			http.Error(w, fmt.Sprintf("Scene not found: %v", initScene), http.StatusNotFound)
			return
		}
	}

	initItemJ, e := json.Marshal(initItem)
	if e != nil {
		msg := fmt.Sprintf("Error marshalling initial item: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	e = templ.ExecuteTemplate(w, "Main", &renderPage{
		InitialItem: string(initItemJ),
	})
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

type renderPage struct {
	InitialItem string
}

type handler struct {
	sync.Mutex
}

func (h *handler) renderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	defer func() {
		if e := r.Body.Close(); e != nil {
			log.Printf("Error closing body: %v\n", e)
		}
	}()

	if r.Method == "GET" {
		log.Println("redirect to /")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Read scene from body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("Error reading body of POST: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	h.Lock()
	defer h.Unlock()

	// Run raytrace command
	tmpImage := "img/tmp." + r.FormValue("format")
	if e := runRaytrace(tmpImage, string(body)); e != nil {
		msg := fmt.Sprintf("Error running raytrace: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Save scene file
	tmpScene := "img/tmp.json"
	if e := saveScene(tmpScene, body); e != nil {
		msg := fmt.Sprintf("Error saving %s: %v", tmpScene, e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Deal with history and get final image name
	imgName, e := saveImageAndScene(tmpImage, tmpScene)
	if e != nil {
		msg := fmt.Sprintf("Error saving/renmaing img/scene files: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Send final image name back to client
	fmt.Fprintln(w, imgName)
}

func (h *handler) historyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	defer func() {
		if e := r.Body.Close(); e != nil {
			log.Printf("Error closing body: %v\n", e)
		}
	}()

	h.Lock()
	defer h.Unlock()

	if r.FormValue("_") == "recent" {
		history, e := getRecentHistory()
		if e != nil {
			msg := fmt.Sprintf("Error getting history: %v", e)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, history)
		return
	}

	start, end := r.FormValue("start"), r.FormValue("end")

	allHistory, e := getAllHistory()
	if e != nil {
		msg := fmt.Sprintf("Error getting all history: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Remember that the oldest item is last
	defaultStartDate := allHistory[len(allHistory)-1].Date
	defaultEndDate := allHistory[0].Date

	// Use default if unspecified
	if start == "" {
		start = defaultStartDate.Format("2006-01-02")
	}
	if end == "" {
		end = defaultEndDate.Format("2006-01-02")
	}

	// Parse dates, if malformed then use default
	startDate, e := time.ParseInLocation("2006-01-02", start, time.Local)
	if e != nil {
		log.Printf("Error parsing start time: %v: %v", start, e)
		startDate = defaultStartDate
	}
	endDate, e := time.ParseInLocation("2006-01-02", end, time.Local)
	if e != nil {
		log.Printf("Error parsing end time: %v: %v\n", end, e)
		endDate = defaultEndDate
	}

	e = templ.ExecuteTemplate(w, "History", &historyPage{
		DefaultStartDate: defaultStartDate.Format("2006-01-02"),
		DefaultEndDate:   defaultEndDate.Format("2006-01-02"),
		StartDate:        start,
		EndDate:          end,
		Items:            historyRange(allHistory, startDate, endDate),
	})
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

type historyPage struct {
	DefaultStartDate, DefaultEndDate string
	StartDate, EndDate               string
	Items                            []*lib.RenderItem
}

func historyRange(history []*lib.RenderItem, start, end time.Time) []*lib.RenderItem {
	if end.Before(start) {
		end = start
	}

	// Add 1 day to include the end day
	end = end.Add(time.Duration(24) * time.Hour)
	firstIndex := sort.Search(len(history), func(i int) bool {
		return !end.Before(history[i].Date)
	})

	lastIndex := sort.Search(len(history), func(i int) bool {
		return !start.Before(history[i].Date)
	})

	return history[firstIndex:lastIndex]
}

func runRaytrace(imgName, scene string) error {
	cmd := exec.Command("raytrace", "-out", imgName)

	// Grab stdin and send scene through it
	cmdIn, e := cmd.StdinPipe()
	if e != nil {
		return fmt.Errorf("getting stdin pipe: %v", e)
	}
	fmt.Fprintln(cmdIn, scene)
	if e := cmdIn.Close(); e != nil {
		return fmt.Errorf("closing stdin to raytrace: %v", e)
	}

	if output, e := cmd.CombinedOutput(); e != nil {
		return fmt.Errorf("running and getting output: %v: %s", e, string(output))
	}

	return nil
}

func saveScene(name string, body []byte) error {
	sceneFile, e := os.Create(name)
	if e != nil {
		return fmt.Errorf("opening %s: %v", name, e)
	}
	_, e = io.Copy(sceneFile, bytes.NewReader(body))
	if e != nil {
		return fmt.Errorf("copying to %s: %v", name, e)
	}
	e = sceneFile.Close()
	if e != nil {
		return fmt.Errorf("closing %s: %v", name, e)
	}
	return nil
}

func saveImageAndScene(tmpImg, tmpScn string) (newImg string, err error) {
	timeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	base := "img/render-" + timeStamp
	newImg = base + filepath.Ext(tmpImg)
	newScn := base + ".json"

	if e := os.Rename(tmpImg, newImg); e != nil {
		return "", fmt.Errorf("renaming %s to %s: %v", tmpImg, newImg, e)
	}

	if e := os.Rename(tmpScn, newScn); e != nil {
		return "", fmt.Errorf("renaming %s to %s: %v", tmpScn, newScn, e)
	}

	if e := saveThumbnail(newImg); e != nil {
		return "", fmt.Errorf("saving thumbnail: %v", e)
	}

	return
}

func saveThumbnail(imgName string) (err error) {
	imgFile, e := os.Open(imgName)
	if e != nil {
		return fmt.Errorf("opening image %s: %v", imgName, e)
	}
	defer func() {
		if e := imgFile.Close(); e != nil && err == nil {
			err = fmt.Errorf("closing img: %v", e)
		}
	}()

	img, _, e := image.Decode(imgFile)
	if e != nil {
		return fmt.Errorf("decoding image %s: %v", imgName, e)
	}

	thumbImg := resize.Thumbnail(thumbnailSize, thumbnailSize, img, resize.Bilinear)

	format := filepath.Ext(imgName)
	thumbFileName := strings.Split(imgName, ".")[0] + "-thumb" + format

	thumbFile, e := os.Create(thumbFileName)
	if e != nil {
		return fmt.Errorf("creating thumbnail file: %v", e)
	}
	defer func() {
		if e := thumbFile.Close(); e != nil && err == nil {
			err = fmt.Errorf("closing thumbnail image: %v", e)
		}
	}()

	switch format[1:] {
	case "jpg", "jpeg":
		jpeg.Encode(thumbFile, thumbImg, nil)
	case "png":
		png.Encode(thumbFile, thumbImg)
	default:
		panic(fmt.Sprintf("Saving image as format %s is not implemented", format[1:]))
	}

	return nil
}

type byDate []*lib.RenderItem

func (h byDate) Len() int {
	return len(h)
}

func (h byDate) Less(i, j int) bool {
	return h[i].Date.Before(h[j].Date)
}

func (h byDate) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func getRecentHistory() (historyJSON string, err error) {
	all, e := getAllHistory()
	if e != nil {
		return "", fmt.Errorf("getting all history: %v", e)
	}

	j, e := json.Marshal(all[:int(math.Min(float64(len(all)), float64(maxHistory)))])
	if e != nil {
		return "", fmt.Errorf("marshal recent history to json: %v", e)
	}

	return string(j), nil
}

// Returns all history sorted by date with newest first
func getAllHistory() (items []*lib.RenderItem, err error) {
	dir, e := os.Open("img")
	defer func() {
		if e := dir.Close(); e != nil && err == nil {
			err = fmt.Errorf("closing img dir: %v", e)
		}
	}()
	if e != nil {
		return nil, fmt.Errorf("opening img dir: %v", e)
	}

	allFiles, e := dir.Readdirnames(-1)
	if e != nil {
		return nil, fmt.Errorf("reading dirnames: %v", e)
	}

	itemMap := make(map[time.Time]*lib.RenderItem)
	for _, f := range allFiles {
		if match := lib.RenderFileRe.FindStringSubmatch(f); match != nil {
			unixTime, e := strconv.ParseInt(match[1], 10, 64)
			if e != nil {
				// Should be impossible since the regex ensures the match is an int
				panic(fmt.Errorf("converting time of %s to int: %v", f, e))
			}
			date := time.Unix(0, unixTime)
			if _, ok := itemMap[date]; !ok {
				itemMap[date] = &lib.RenderItem{Date: date}
			}
			isThumb := match[2] != ""
			format := match[3]
			if isThumb {
				itemMap[date].Thumb = "img/" + f
			} else if format == "json" {
				itemMap[date].Scene = "img/" + f
			} else {
				itemMap[date].Render = "img/" + f
			}
		}
	}

	items = make([]*lib.RenderItem, len(itemMap))
	i := 0
	for _, item := range itemMap {
		items[i] = item
		i++
	}
	sort.Sort(sort.Reverse(byDate(items)))

	return
}
