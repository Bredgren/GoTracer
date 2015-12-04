package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"text/template"
)

const path = "/src/github.com/Bredgren/gotracer/webtracer"

var goPath = os.Getenv("GOPATH")

var templ *template.Template

var (
	// debug is set by the -D command line flag.
	debug      bool
	port       int
	maxHistory int
)

func init() {
	flag.BoolVar(&debug, "D", false, "Debug mode. More logs, use unminified assets, etc.")
	flag.IntVar(&port, "p", 8080, "Set http port")
	flag.IntVar(&maxHistory, "history", 20, "Maximum render history to keep track of")
}

func main() {
	setup()

	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/render", renderHandler)
	http.HandleFunc("/history", historyHandler)
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
		"debug": func() bool { return debug },
	}).ParseFiles("./tmpl/page.tmpl"))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
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

	renderTmpl(w, &page{})
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
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
		msg := fmt.Sprintf("Error reading post body: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fileMu.Lock()
	defer fileMu.Unlock()

	// Run raytrace command
	tmpImage := "img/render.jpg"
	if e := runRaytrace(tmpImage, string(body)); e != nil {
		msg := fmt.Sprintf("Error running raytrace: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Save scene file
	tmpScene := "img/scene.json"
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

func historyHandler(w http.ResponseWriter, r *http.Request) {
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

	fileMu.Lock()
	defer fileMu.Unlock()

	history, e := getHistory()
	if e != nil {
		msg := fmt.Sprintf("Error getting history: %v", e)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	fmt.Fprintln(w, history)
}

func renderTmpl(w http.ResponseWriter, p *page) {
	e := templ.ExecuteTemplate(w, "Main", p)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "Not implemented")
}

type page struct {
	// Nothing here for now...
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

var renderFileRe = regexp.MustCompile(`render(\d+).jpg`)
var sceneFileRe = regexp.MustCompile(`scene(\d+).jpg`)
var fileMu sync.Mutex

// This helper figures out what name to give new files and shifts files around so that
// they're consistently numbered, 0 is oldest and maxHistory-1 is newest.
func saveImageAndScene(tmpImg, tmpScn string) (newImgName string, err error) {
	// Combine file pairs and remove pairs that have at least one missing
	pairs, e := getFilePairs()
	if e != nil {
		return "", fmt.Errorf("getting file pairs: %v", e)
	}
	// Add the new ones
	pairs = append(pairs, [2]string{tmpImg, tmpScn})
	if len(pairs) == maxHistory+1 {
		// We're full, make room by clearing the first entry, compress will remove it
		pairs[0][0] = ""
	}
	log.Println("pairs before", pairs)
	pairs = compress(pairs)
	log.Println("pairs after", pairs)

	// Rename files based on their index
	for i, p := range pairs {
		newImgName = "img/render" + strconv.Itoa(i) + ".jpg"
		e := os.Rename(p[0], newImgName)
		if e != nil {
			return "", fmt.Errorf("renaming %s to %s: %v", p[0], newImgName, e)
		}
		newSceneName := "img/scene" + strconv.Itoa(i) + ".jpg"
		e = os.Rename(p[1], newSceneName)
		if e != nil {
			return "", fmt.Errorf("renaming %s to %s: %v", p[1], newSceneName, e)
		}
	}
	return newImgName, nil
}

func getHistory() (history string, err error) {
	pairs, e := getFilePairs()
	if e != nil {
		return "", fmt.Errorf("getting file pairs: %v", e)
	}
	pairs = compress(pairs)

	j, e := json.Marshal(pairs)
	if e != nil {
		return "", fmt.Errorf("marshal pairs to json: %v", e)
	}

	return string(j), nil
}

func getFilePairs() (pairs [][2]string, err error) {
	dir, e := os.Open("img")
	defer func() {
		// Only reporet this error if no others have happened
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

	renderFiles := make([]string, maxHistory)
	sceneFiles := make([]string, maxHistory)
	for _, f := range allFiles {
		if i, match := getFileIndex(renderFileRe, f); match && i < maxHistory {
			renderFiles[i] = "img/" + f
		} else if i, match := getFileIndex(sceneFileRe, f); match && i < maxHistory {
			sceneFiles[i] = "img/" + f
		}
	}

	return zip(renderFiles, sceneFiles), nil
}

func getFileIndex(re *regexp.Regexp, file string) (int, bool) {
	if match := re.FindStringSubmatch(file); match != nil {
		i, e := strconv.Atoi(match[1])
		if e != nil {
			// Should be impossible since the regex ensures the match is an int
			panic(fmt.Errorf("converting index of %s to int: %v", file, e))
		}
		return i, true
	}
	return -1, false
}

func zip(s1, s2 []string) (zipped [][2]string) {
	// Assumes len(s1) == len(s2)
	zipped = make([][2]string, len(s1))
	for i := range s1 {
		zipped[i][0] = s1[i]
		zipped[i][1] = s2[i]
	}
	return
}

func compress(pairs [][2]string) [][2]string {
	// Removes elements that are missing or aren't paired
	var toRemove []int
	for i, p := range pairs {
		if p[0] == "" || p[1] == "" {
			toRemove = append(toRemove, i)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(toRemove)))
	for _, i := range toRemove {
		pairs = append(pairs[:i], pairs[i+1:]...)
	}
	return pairs
}
