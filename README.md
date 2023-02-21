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

- [ ] Do not fetch github data with regular http client and github api client.
  In this case we should use only GitHub API client.
- [ ] Parallel fact gathering
- [x] Gather facts for big awesome lists (rate limiters, auto-continue process)
- [ ] Render Markdown
- [ ] Render HTML
- [ ] Render Sitemap
- [ ] Custom HTML template
    - [ ] Custom CSS
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
- [ ] Microformats and meta
    - [ ] Twitter cards
    - [ ] OpenGraph
    - [ ] Keywords, description
- [ ] Advertisement
    - [ ] Link, Title, Description, Logo, URL, Priority
    - [ ] Links in head/links in footer/links everywhere/
    - [ ] Google Ads
- [ ] Trackers and pixels
    - [ ] Google, FB, etc
    - [ ] Auto-add utm_source to external links
    - [ ] rel=nofollow
- [ ] Manifest.json
- [ ] Create separate pages

