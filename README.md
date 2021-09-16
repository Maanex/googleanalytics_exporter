[![Build Status](https://travis-ci.org/paha/googleanalytics_exporter.svg?branch=master)](https://travis-ci.org/paha/googleanalytics_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/paha/googleanalytics_exporter)](https://goreportcard.com/report/github.com/paha/googleanalytics_exporter)
[![Docker Repository on Quay](https://quay.io/repository/paha/ga-prom/status "Docker Repository on Quay")](https://quay.io/repository/paha/ga-prom)

# Google Real Time Analytics to Prometheus

Obtains Google Analytics RealTime metrics, and presents them to prometheus for scraping.

This fork makes slight changes to the configuration by allowing you to configure everything through environment variables and docker secrets so you don't have to compile the container every time you want to change some settings. Check out the docker-compose file for more info.

---

## Quick start

1. Edit the properties in the docker-compose.yml file.

### ViewID for the Google Analytics

From your Google Analytics Web UI: *Admin (Low left) ==> View Settings (far right tab, named VIEW)'*

*View ID* should be among *Basic Settings*. Prefix `ga:` must be added to the ID, e.g. `ga:1234556` while adding it to the config.

### Google creds

[Google API manager][2] allows to create OAuth 2.0 credentials for Google APIs. Use *Service account key* credentials type, upon creation a json creds file will be provided. Project RO permissions should be sufficient.

>*The email from GA API creds must be added to analytics project metrics will be obtained from.*>

## Authors

Original Author - Pavel Snagovsky, pavel@snagovsky.com

Modifications - Andreas May, andreas@maanex.me

## License

Licensed under the terms of [MIT license][4], see [LICENSE][5] file

[1]: https://github.com/Masterminds/glide
[2]: https://console.developers.google.com/apis/credentials
[3]: https://hub.docker.com/_/alpine/
[4]: https://choosealicense.com/licenses/mit/
[5]: ./LICENSE
