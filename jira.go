package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	chglog "github.com/git-chglog/git-chglog"
	"github.com/utrace-ltd/changelogger"
	"github.com/voxelbrain/goptions"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func main() {
	var opt Options
	goptions.ParseAndFail(&opt)

	config, err := NewConfig()
	if err != nil {
		log.Panicf("Error on config load: %s", err)
	}

	jiraClient, err := jira.NewClient(createHttpClient(config), config.Jira.URL)
	if err != nil {
		log.Panicf("Error on parsh base url: %s", err)
	}

	gitVersions, err := changelogger.NewChangeLogger().GetChangeLog(opt.Version)
	if err != nil {
		log.Panicf("Error on get git versions on tag %s: %s", opt.Version, err)
	}

	if len(gitVersions) == 0 {
		log.Panic("Error, no git versions")
	}

	issues := getIssueFromChangeLog(gitVersions[0], config.Jira.Issue.Pattern)

	project, _, err := jiraClient.Project.Get(config.Jira.Project.Id)
	if err != nil {
		log.Panicf("Error on get project by id (%s): %s", config.Jira.Project.Id, err)
	}

	jiraVersion, err := findOrCreateVersion(project, createVersionName(opt.ComponentName, opt.Version), gitVersions[0].Tag.Date, jiraClient)
	if err != nil {
		log.Panicf("Error on create/get version: %s", err)
	}

	updateTasksVersions(jiraVersion, issues, jiraClient)
}

func findOrCreateVersion(p *jira.Project, name string, rd time.Time, c *jira.Client) (*jira.Version, error) {
	for _, ver := range p.Versions {
		if ver.Name == name {
			return &ver, nil
		}
	}
	pid, _ := strconv.Atoi(p.ID)
	n := jira.Version{
		Name:        name,
		Description: "",
		Archived:    false,
		ProjectID:   pid,
		ReleaseDate: rd.String(),
		Released:    false,
	}

	v, _, err := c.Version.Create(&n)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getIssueFromChangeLog(g *chglog.Version, issuePattern string) []string {
	var issues []string
	keys := make(map[string]bool)
	re := regexp.MustCompile(issuePattern)
	for _, c := range g.Commits {
		issue := string(re.Find([]byte(c.Header)))
		if issue != "" {
			if _, value := keys[issue]; !value {
				keys[issue] = true
				issues = append(issues, issue)
			}
		}
	}
	return issues
}

func updateTasksVersions(v *jira.Version, issues []string, c *jira.Client) {
	for _, i := range issues {
		issue, _, err := c.Issue.Get(i, nil)
		if err != nil {
			log.Printf("Error on get issue by name: %s", err)
		}
		data := creatIssueDataFromIssue(issue)
		needSave := syncIssueVersions(data, v)
		if needSave {
			_, _, err := c.Issue.Update(data)
			if err != nil {
				log.Printf("Error on save issue %s: %s", issue.ID,  err)
			}
		}
	}
}

func syncIssueVersions(i *jira.Issue, v *jira.Version) bool {
	for _, ver := range i.Fields.FixVersions {
		if ver.ID == v.ID {
			return false
		}
	}
	i.Fields.FixVersions = append(i.Fields.FixVersions, createFixVersionFromVersion(v))
	log.Printf("In issue %s add version %s", i.Key, v.Name)
	return true
}

func creatIssueDataFromIssue(issue *jira.Issue) *jira.Issue {
	return &jira.Issue{Key:issue.Key, ID: issue.ID, Fields:&jira.IssueFields{FixVersions: issue.Fields.FixVersions}}
}

func createFixVersionFromVersion(v *jira.Version) *jira.FixVersion {
	return &jira.FixVersion{Name:v.Name, ID: v.ID, ProjectID: v.ProjectID}
}

func createHttpClient(config *Config) *http.Client {
	h := jira.BasicAuthTransport{Username: config.Jira.User, Password: config.Jira.Password}
	return h.Client()
}

func createVersionName(component, version string) string  {
	return fmt.Sprintf("%s:%s", component, version)
}
