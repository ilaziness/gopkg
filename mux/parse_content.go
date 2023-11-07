package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var sourceUrl = "https://raw.githubusercontent.com/avelino/awesome-go/main/README.md"
var sourceFile = "test.md"
var saveFileName = "awesome-go.json"

var matchMenu = regexp.MustCompile(`Contents\n\n((.|\n)+?)\n\n`)
var matchContent = regexp.MustCompile(`## Audio and Music(.|\n)+`)
var matchMenuName = regexp.MustCompile(`\[(.+?)\]`)
var matchMenuLink = regexp.MustCompile(`#[\w-]+`)
var matctMenuStarBlack = regexp.MustCompile(`(\s(\s)+)-`)
var matchGithubStar = regexp.MustCompile(`"repo-stars-counter-star".+?>([\w\.]+)</span>`)
var matchLastUpdateDate = regexp.MustCompile(`<relative-time datetime="(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)"`)
var matchCommitUrl = regexp.MustCompile(`<include-fragment\s=?src="(.+?)"\s+?class`)

type StoreStruct struct {
	Menu    []*MenuItem   `json:"menu"`
	Content []ContentItem `json:"content"`
}

type MenuItem struct {
	Name   string `json:"name"`
	NameZh string `json:"name_zh"`
	Link   string `json:"link"`
	Level  int    `json:"level"`

	Children []*MenuItem `json:"children"`
}

type ContentItem struct {
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	DescriptionZH string    `json:"description_zh"`
	Link          string    `json:"link"`
	Level         int       `json:"level"`
	Projects      []Project `json:"projects"`
}

type Project struct {
	Name           string `json:"name"`
	Star           string `json:"star"`
	Description    string `json:"description"`
	DescriptionZH  string `json:"description_zh"`
	RepositoryUrl  string `json:"repository_url"`
	LastUpdateDate string `json:last_update_date`
}

func getContent() []byte {
	ct, err := os.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}
	return ct
}

func parseMenu() ([]*MenuItem, error) {
	menus := make([]*MenuItem, 0)
	menuData := matchMenu.FindSubmatch(getContent())
	if len(menuData) == 0 {
		return menus, errors.New("not found menu data")
	}
	menuList := make(map[string]*MenuItem)
	menuLines := strings.Split(string(menuData[1]), "\n")
	menuItemStack := make([]*MenuItem, 0)
	for no, line := range menuLines {
		if no == 1 {
			continue
		}
		blank := matctMenuStarBlack.FindStringSubmatch(line)
		level := 0
		if len(blank) > 0 {
			level = len(blank[1])
		}
		name := matchMenuName.FindStringSubmatch(line)
		if len(name) == 0 {
			continue
		}
		item := &MenuItem{
			Name:     name[1],
			Link:     matchMenuLink.FindString(line),
			Level:    level,
			Children: make([]*MenuItem, 0),
		}
		menuList[item.Name] = item
		if item.Level == 0 {
			menuItemStack = menuItemStack[0:0]
			// 顶级
			menus = append(menus, item)
			menuItemStack = append(menuItemStack, item)
			continue
		}
		if menuItemStack[len(menuItemStack)-1].Level < item.Level {
			menuItemStack[len(menuItemStack)-1].Children = append(menuItemStack[len(menuItemStack)-1].Children, item)
			menuItemStack = append(menuItemStack, item)
			continue
		}
		if menuItemStack[len(menuItemStack)-1].Level == item.Level {
			menuItemStack[len(menuItemStack)-2].Children = append(menuItemStack[len(menuItemStack)-2].Children, item)
			menuItemStack[len(menuItemStack)-1] = item
			continue
		}
		if menuItemStack[len(menuItemStack)-1].Level > item.Level {
			popNum := 0
			for _, itemStack := range menuItemStack {
				if itemStack.Level != item.Level {
					popNum++
					continue
				}
				menuItemStack = menuItemStack[0 : len(menuItemStack)-popNum]
				menuItemStack[len(menuItemStack)-2].Children = append(menuItemStack[len(menuItemStack)-2].Children, item)
				menuItemStack[len(menuItemStack)-1] = item
				break
			}
		}
	}
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Name < menus[j].Name
	})
	//printMenus(menus)

	contets := make([]ContentItem, 0)
	content := matchContent.Find(getContent())
	if len(content) == 0 {
		return nil, errors.New("not found content")
	}
	contentLines := strings.Split(string(content), "\n")
	for _, line := range contentLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "**[") {
			continue
		}
		level := 0
		if strings.HasPrefix(line, "##") {
			level = 2
		}
		if strings.HasPrefix(line, "###") {
			level = 3
		}
		if strings.HasPrefix(line, "####") {
			level = 4
		}
		if level > 0 {
			menuName := strings.TrimSpace(line[level:])
			if menuName == "Guided Learning" {
				menuName = "Guided Learning Paths"
			}
			//fmt.Printf("%d(%s)\n", level, menuName)
			cti := ContentItem{
				Title:    menuName,
				Level:    level,
				Link:     menuList[menuName].Link,
				Projects: make([]Project, 0),
			}
			contets = append(contets, cti)
			continue
		}

		// 分类描述
		if strings.HasPrefix(line, "_") {
			contets[len(contets)-1].Description = strings.Trim(line, "_.")
		}
		if strings.HasPrefix(line, "- [") {
			name := strings.Split(line, "](")
			pjt := Project{
				Name:          strings.TrimLeft(name[0], "- ["),
				Description:   strings.TrimRight(strings.TrimSpace(line[strings.LastIndex(line, ") -")+3:]), "."),
				RepositoryUrl: strings.TrimSpace(name[1][0:strings.Index(name[1], ")")]),
			}
			star, updateDate := parseGithubInfo(pjt.RepositoryUrl)
			pjt.Star = star
			pjt.LastUpdateDate = updateDate
			contets[len(contets)-1].Projects = append(contets[len(contets)-1].Projects, pjt)
		}
	}

	// print content
	// for _, ct := range contets {
	// 	fmt.Printf("content: %d - %s - %s\n", ct.Level, ct.Title, ct.Link)
	// 	for _, pj := range ct.Projects {
	// 		fmt.Printf("project: %s - %s - %s\n", pj.Name, pj.RepositoryUrl, pj.Description)
	// 	}
	// }

	slog.Info("save to file")
	saveToFile(menus, contets)

	return menus, nil
}

func printMenus(menus []*MenuItem) {
	for _, item := range menus {
		fmt.Printf("menu: %d - [%s] - %s\n", item.Level, item.Name, item.Link)
		if len(item.Children) > 0 {
			printMenus(item.Children)
		}
	}
}

func saveToFile(menus []*MenuItem, content []ContentItem) {
	f, err := os.OpenFile(saveFileName, os.O_CREATE|os.O_TRUNC, os.FileMode(755))
	if err != nil {
		panic(err)
	}
	p, _ := os.Getwd()
	fmt.Print(p)
	defer func() {
		err = f.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("close file error: %s", err))
		}
	}()
	store := StoreStruct{
		Menu:    menus,
		Content: content,
	}
	res, err := json.Marshal(store)
	if err != nil {
		slog.Error(fmt.Sprintf("json Marshal error: %s", err))
		return
	}
	_, err = f.Write(res)
	if err != nil {
		slog.Error(fmt.Sprintf("save to file error: %s", err))
	}
}

func parseGithubInfo(urlStr string) (star, lastUpdateDate string) {
	slog.Info(urlStr)
	if !strings.HasPrefix(urlStr, "https://github.com") {
		return
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		slog.Warn(fmt.Sprintf("new req error: %s", err))
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0")
	if http.DefaultClient.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}
		http.DefaultClient.Jar = jar
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn(fmt.Sprintf("req error: %s", err))
		return
	}
	resContent, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Warn(fmt.Sprintf("read response body error: %s", err))
		return
	}
	matchRes := matchGithubStar.FindSubmatch(resContent)
	if len(matchRes) == 2 {
		star = string(matchRes[1])
	}

	retryCount := 0
retry:
	updatedate := matchLastUpdateDate.FindSubmatch(resContent)
	if len(updatedate) == 2 {
		lastUpdateDate = string(updatedate[1])
	} else if retryCount == 0 {
		commitUrl := matchCommitUrl.FindSubmatch(resContent)
		if len(commitUrl) == 2 {
			urlObj, err := url.Parse(fmt.Sprintf("https://github.com%s", commitUrl[1]))
			if err != nil {
				slog.Error(fmt.Sprintf("format uri (%s) error: %s", commitUrl[1], err))
			}
			req.URL = urlObj
			fmt.Printf("%s\n", commitUrl[1])
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				slog.Warn(fmt.Sprintf("req commit url error: %s", err))
			} else {
				resContent, err = io.ReadAll(res.Body)
				if err != nil {
					slog.Warn(fmt.Sprintf("read commit url response body error: %s", err))
					return
				} else {
					retryCount++
					goto retry
				}
			}
		}
	}
	if star == "" {
		slog.Warn("not found star text: " + urlStr)
	}
	if lastUpdateDate == "" {
		slog.Warn("not found last update date text: " + urlStr)
	}
	slog.Info(fmt.Sprintf("%s - %s", star, lastUpdateDate))
	return
}
