Visit Logger
============

Visit Logger is made up of three componanats:

 1. Chrome extension "Visit Logger"
 2. Go based web service "httplogger"
 3. Vue.js frontend "Chrome History Report"


# Deployment
 
## Visit Logger

The extension must be zipped and uploaded to the [Chrome Web Store Devloper Dashboard](https://chrome.google.com/webstore/devconsole).


## httplogger
 
The service includes a makefile to compile the Go app and a Docker Compose file (compose.yml) to deploy the service.

External Docker volumes are needed to contain Caddy configs and httplogger API key. The following process populates these on first deploy:

	docker volume create caddy
	docker volume create caddy_data
	docker volume create secrets

  docker compose up -d

  docker cp Caddyfile httplogger-caddy:/ect/caddy/
  docker cp apikey httplogger-httplogger:/secrets/

  docker compose restart


## Chrome History Report

The frontend app is hosted on Netlify and a link exists between the GitHub repo and Netlify. Commits to main will trigger a deployment.
