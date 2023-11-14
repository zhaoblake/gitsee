package main

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"strings"
)

func scan(folder string) {
	repos := ScanFolder(folder)
	filePath := getDotFilePath()
	addRepos(filePath, repos)
}

func ScanFolder(folder string) []string {
	return scanFolders(make([]string, 0), folder)
}

// scanFolders 递归扫描出所有包含 .git 文件的文件夹
func scanFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatalf("%s 打开异常：%s", folder, err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println("关闭文件句柄异常：", err)
		}
	}()

	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				log.Println(path)
				folders = append(folders, path)
				continue
			}
			folders = scanFolders(folders, path)
		}

	}

	return folders
}

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gitseelocalstats"

	return dotFile
}

// addRepos 向文件中添加仓库
func addRepos(filePath string, newRepos []string) {
	repos := getRepos(filePath)

	for _, nr := range newRepos {
		// 去重
		if !contains(repos, nr) {
			repos = append(repos, nr)
		}
	}

	err := os.WriteFile(filePath, []byte(strings.Join(repos, "\n")), 0755)
	if err != nil {
		panic(err)
	}
}

func contains(repos []string, repo string) bool {
	for _, r := range repos {
		if r == repo {
			return true
		}
	}
	return false
}

// getRepos 从文件中读取仓库
func getRepos(filePath string) []string {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println("关闭文件句柄异常：", err)
		}
	}()

	var repos []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		repos = append(repos, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(err)
	}

	return repos
}
