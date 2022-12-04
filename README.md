# Awesome tool for Awesome repos

![Awesome Tool](./cover.svg)

This tool is allow to create and manage awesome lists easily.

## How it works

## Quickstart

- Write YAML file with links and groups
- Run the tool and get resulted Markdown and HTML


### Recommended process

- Create YAML file, which contains all links, groups and other meta-information.
- Run the tool on source YAML file, which create a resulted build-file (YAML) 
  with all datae, you written and additional metadata (url-checks, github repo 
  info, etc)
- Check the changes between new-generation of generated file
- Run the tool on generated file to produce outputs in Markdown or HTML
- Commit changes

This process is like npm or go package managers, which use separate file to 
store calculated data. This is allows you to exactly known what you get and 
track changes across commits.

## Features

- [ ] Gather facts for big awesome lists (rate limiters, auto-continue process)
- [ ] Render Markdown
- [ ] Render HTML
- [ ] Custom HTML template
- [ ] Track all changes with `.sum` file
- [ ] Check links (http statuses)
- [ ] Fetch data for GitHub repos 
  - [ ] Archived or not
  - [ ] Presented or not
  - [ ] Fetch counters: stars, followers, etc
  - [ ] Fetch other significant data: contributors, tags, lang, last commit, 
        etc.
- [ ] Set restrictions for the links
  - [ ] On GitHub counters 
  - [ ] On last commit
  - [ ] On http statuses

