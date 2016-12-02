package main

import (
	"flag"
	"fmt"
	"github.com/octokit/go-octokit/octokit"
	"gopkg.in/libgit2/git2go.v24"
	"log"
	"os"
	"regexp"
)

func main() {
	repositoryPath := flag.String("repo", ".", "Path to the git repository")
	flag.Parse()

	repo, err := git.OpenRepository(*repositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	upstreamRemote, err := repo.Remotes.Lookup("upstream")
	if upstreamRemote != nil {
		log.Println("Remote 'upstream' already exists.")
		os.Exit(1)
	}

	originRemote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		fmt.Println("Could not find remote 'origin'.")
		os.Exit(1)
	}

	r := regexp.MustCompile(`github.com\/(.*)/(.*)\.git`)
	matched := r.FindAllStringSubmatch(originRemote.Url(), -1)
	matchedOwner := matched[0][1]
	matchedRepo := matched[0][2]

	client := octokit.NewClient(nil)
	remoteRepo, _ := client.Repositories().One(nil, octokit.M{"owner": matchedOwner, "repo": matchedRepo})
	if remoteRepo == nil {
		fmt.Println("Could not find the repository. This is probably a private repo.")
		os.Exit(1)
	}

	if !remoteRepo.Fork {
		fmt.Println("This repository is not forked repository.")
		os.Exit(1)
	}

	repo.Remotes.Create("upstream", remoteRepo.Parent.CloneURL)
	fmt.Printf("git remote add upstream %s\n", remoteRepo.Parent.CloneURL)
}
