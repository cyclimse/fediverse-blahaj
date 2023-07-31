# TODOs

## Guidelines

A nice website with a list of Fediverse servers and some graphs. Should be helpful to understand the Fediverse and find new servers to join.

We could have global stats about the Fediverse and stats about each software (eg: Mastodon, Pleroma, ...). We could also have a list of servers that are currently down (like a [down detector](https://downdetector.com/)).

The process of crawling the Fediverse should be as transparent as possible. We should be able to explain how we crawl the Fediverse and what we do with the data. It's important to respect the privacy of the users and the servers by excluding personal servers from the API. We should have multiple ways to opt-out of the crawling process (eg: robots.txt, manual opt-out, ...).

## MVP

### Design

* [ ] Ubiquitous language: clarify if we wan't to use Fediverse servers of Fediverse instances
* [ ] Find a name for the project

### Crawler

* [x] Crawl Fediverse servers
  * [x] Acknowledge robots.txt
* [x] Orchestrate multiple crawlers
* [ ] Support various different Fediverse servers
  * [ ] Mastodon and Mastodon likes (Pleroma, Misskey, ...)
  * [ ] Lemmy
  * [ ] Peertube
* [ ] Improve Crawler seeding process (eg: seed from servers that were discovered but not crawled yet)
* [ ] Improve Crawler security (data sanitization, ...)
* [ ] Offer a way to opt-out of crawling (can be manual process for MVP because robots.txt is respected)
* [ ] Exclude ngrok.io from the Crawler and some other domains (eg: localhost, ...)

### API

* [ ] Implement API
  * [ ] Implement pagination
  * [ ] Implement filtering
    * [ ] Filter by server domain
    * [ ] Filter by software (eg: Mastodon, Pleroma, ...)
    * [ ] Filter by number of users
* [ ] Exclude small servers from the API to (eg. personal servers)

### Frontend

Frontend should be kept simple for MVP. With OpenAPI, we can auto-generate a SDK for the frontend.

* [ ] Homepage
* [ ] About page
* [ ] Search bar
* [ ] List of Fediverse servers
* [ ] Detail page for a Fediverse server
  * [ ] Have server description
  * [ ] Last infos from the Crawler
    * [ ] Number of users
    * [ ] Number of posts
* [ ] Graphs
  * [ ] Number of users over time
  * [ ] Number of posts over time

### Deployment

* [ ] Deploy the Crawler (via a cronjob)
* [ ] Deploy the API
* [ ] Deploy the Frontend (bunnyCDN?)

## Post MVP

### Frontend

* [ ] Full-text search

### Crawler

* [ ] Export more metrics related to the Crawling process
  * [ ] Number of unreachable servers
  * [ ] Number of crawls blocked by robots.txt
* [ ] Research on blocking detection for Fediverse servers
  * [ ] Detect blocking
  * [ ] Detect blocking reason
  * [ ] Detect blocking duration

### Mastodon bot

* [ ] Create a Mastodon bot
* [ ] Offer a simple way to opt-out of crawling (eg: send a DM to the bot)
* [ ] Answer to mentions
  * [ ] Is this server down?
  * [ ] Server stats

## Next steps