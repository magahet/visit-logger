Visit Logger
============



Visit Logger is made up of three components:

 1. Chrome extension "Visit Logger"
 2. Go based web service "httplogger"
 3. Vue.js frontend "Chrome History Report"


# Deployment
 
## Visit Logger

The Visit Logger extension watches Chrome tab changes and sends page viewing history to a user configured remote server.

This allow users to record their web browsing history on a remote server. The intention is to allow you to review, and perform analytics on, your browsing history data.


### Deployment

The extension must be zipped and uploaded to the [Chrome Web Store Developer Dashboard](https://chrome.google.com/webstore/devconsole).


## httplogger

The httplogger service captures POSTS made by the extension and populates an LRU cache to keep track of recently visited sites. It includes basic auth to ensure privicy.
 
### Deployment

The service includes a makefile to compile the Go app and a Docker Compose file (compose.yml) to deploy the service.

External Docker volumes are needed to contain Caddy configs and httplogger API key. The following process populates these on first deploy:

```
docker volume create caddy
docker volume create caddy_data
docker volume create secrets

docker compose up -d

docker cp Caddyfile httplogger-caddy:/ect/caddy/
docker cp apikey httplogger-httplogger:/secrets/

docker compose restart
```


## Chrome History Report

The frontend app makes calls to the httplogger service to provide basic auth and to render page viewing history reports for each user.

### Deployment

The frontend app is hosted on Netlify and a link exists between the GitHub repo and Netlify. Commits to main will trigger a deployment.
