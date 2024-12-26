# Homarr iFrames

This project was created as a fork of [Homarr iFrames](https://github.com/diogovalentte/homarr-iframes) to adapt it to my personal needs.

I have simply modified some code related to **Vikunja**'s iFrame and the way to generate the Docker image. The modifications in the code for sure are improvable because I have no knowledge of Go.

I have added to the begining of this README the modifications on the original documentation of the Homarr iFrames project. If you want to know more, follow this [link](https://github.com/diogovalentte/homarr-iframes).

I started from tag [v1.2.0](https://github.com/diogovalentte/homarr-iframes/releases/tag/v1.2.0) and made the following changes:

## Vikunja iFrame

The following modifications were made on the Vikunka iFrame:

- A task can have a due date and, in addition, be repeated. I adapted the nesting of conditionals so that, in addition to showing the due date, it also shows whether it has a repetition (in the original project, if the task has a due date, it does not show information about repetitions).

- All tasks can be marked as "Done", not only those without repetition. To make this work correctly, changes are made to both the API and the iFrame. In addition to sending the identifier of the task to be marked as "Done", two additional parameters are sent: "repeat_mode" and "repeat_after". This ensures that the task keeps its repetition cycle. Documentation of these parameters is added.

- I added a new parameter to the iFrame URL called "showCompact" which defaults to false. If it is sent as true, the CSS is slightly adapted to show more tasks per page. Practically twice as many tasks are displayed, although it is certainly less aesthetically pleasing than the original version. Documentation for this parameter is added.

Example URL:

`[PROTOCOL]://[SERVER]:[PORT]/v1/iframe/vikunja?showCreated=false&showCompact=true&api_url=[PROTOCOL]://[SERVER]:[PORT]`

## Docker

When installing Homarr iFrames with Docker on a Raspberry Pi 5, the Docker image weighs 1.3Gb. The reason for this is the Golang image it is based on. Looking for the reason why a Golang image is so heavy, I found the article [Smallest Golang Docker Image](https://klotzandrew.com/blog/smallest-golang-docker-image/) and, with slight modifications, I applied what is indicated there (see `Dockerfile`). This results in a 27.7Mb image (as stated above, on a Raspberry Pi 5).

In the code of this repository, the `Dockerfile` is prepared to run only on Linux arm64 systems (such as Raspberry Pi 5) because of the following line:

```go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main .
```

If you want to run on a different type of processor, modify the `GOOS` and `GOARCH` parameters.

The `docker-compose.yml` file has also been given an external network called "npm" and a label so that [Watchtower](https://containrrr.dev/watchtower/) does not try to update it. It also has a mapping from port 8080 of the container to port 38080 of the host. If you don't have an external network with that name or you want to map the container to a different local port, you have to modify this configuration.

## How to run:

Like the original project, this can be run manually or with Docker. The recommended way is to use the Docker Compose file that comes with the repository and build a new image.

Although it depends on the exact installation of Docker on the machine, the following command (or similar) should work:

```sh
sudo docker-compose up --build -d
```

Once the image creation is finished, it is advisable to do a clean-up of the images that were downloaded during the generation, such as `golang:1.23.4`.

## IMPORTANT!

- This project does not have any authentication system, be careful with its use. In my case it has been put behind a [Nginx Proxy Manager](https://nginxproxymanager.com/) with basic authentication, but any other method can be used.

# Original README

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/8df579cb-9cc9-4bad-a1da-f0cf015e741b)

This project connects to multiple self-hosted applications (called **sources** here) and creates an iFrame to be used in any dashboard (_not only [Homarr](https://github.com/ajnart/homarr), despite the project's name_).

The iFrames will be available under the project's API routes, like `/v1/iframe/linkwarden`. These routes accept query parameters to change the iFrame, like limiting the number of items or specifying whether you want the iFrames to check for updates automatically (_the iframe reloads if the source contents change (like adding new bookmarks on Linkwarden)_).

- You can check all query parameters in the API docs.

# Sources

The API can create iFrames for multiple sources, like the **Vikunja** source that creates an iFrame with your tasks, or the **Linkwarden** source that creates an iFrame with your bookmarks.

The sources may require environment variables with specific information like the application address or credentials. The way you provide these environment variables depends on how you run the API.

- A list of the sources can be found [here](/docs/SOURCES.md).

# API docs

The API docs are under the path `/v1/swagger/index.html`, like `http://192.168.1.44/v1/swagger/index.html` or `https://sub.domain.com/v1/swagger/index.html`, depending on how you access the API.

# Notes

## How to add the iframe to your dashboard

1. In your Homarr dashboard, click on **Enter edit mode -> Add a tile -> Widgets -> iFrame**.
2. Click to edit the new iFrame widget.
3. Add the API URL (`http://192.168.1.15:8080`) + the source path (`/v1/iframe/linkwarden`) + query parameters, like `http://192.168.1.15:8080/v1/iframe/linkwarden?collectionId=1&limit=3&theme=dark`.

## How accessing the iframes works

When you add an iFrame to your dashboard, it's **>your<** web browser that fetches the iFrame from the API and shows it to you, not your dashboard application running on your server. So your browser needs to be able to access the API, that's how an iFrame works.

- **Examples**:
  - If you run this project on your server under port 5000, your browser needs to use your server's IP address + port 5000.
  - If you're accessing your dashboard with a domain and using HTTPS, you also need to access this project's API with a domain and using HTTPS. If you try to use HTTP + HTTPS, your browser will likely block the iFrame.

## No built-in authentication system

This project doesn't have any built-in authentication system, so anyone who can access the API will be able to get all information from the API routes, like your Vikunja tasks, Linkwarden bookmarks, etc. You can add an authentication portal like [Authelia](https://github.com/authelia/authelia) or [Authentik](https://github.com/goauthentik/authentik) in front of the project to secure it, that's how I do it.

# How to run:

- **For Docker and Docker Compose**: by default, the API will be available on port `8080` and is not accessible by other machines. To be accessible by other machines, you need to run the API behind a reverse proxy or run the container in [host network mode](https://docs.docker.com/network/drivers/host/).

- You can change the API port using the environment variable `PORT`.
  - Depending on the port you choose, you need to run the container with user `root` instead of the user `1000` used in the examples and the `docker-compose.yml` file.

## Using Docker:

1. Run the latest version:

```sh
docker run --name homarr-iframes -p 8080:8080 -e VARIABLE_NAME=VARIABLE_VALUE -e VARIABLE_NAME=VARIABLE_VALUE ghcr.io/diogovalentte/homarr-iframes:latest
```

## Using Docker Compose:

1. There is a `docker-compose.yml` file in this repository. You can clone this repository to use this file or create one yourself.
2. Create a `.env` file with the environment variables you want to provide to the API. It should be like the `.env.example` file and be in the same directory as the `docker-compose.yml` file.
3. Start the container by running:

```sh
docker compose up
```

## Manually:

1. Install the dependencies:

```sh
go mod download
```

2. Export the environment variables.
3. Run:

```sh
go run main.go
```
