Build-from scratch doc

following below docs
-> https://github.com/golang-standards/project-layout
-> https://go.dev/doc/modules/layout

```
> git init

Initialized empty Git repository in C:/Users/user-name/folder/path/system design/url-shortener/.git/

> git branch -M main
> git remote add origin https://github.com/Ankush-Hegde/url-shortener.git
> git add .
> git commit -m "initial commit"
> git push --set-upstream origin main

> go mod init url-shortener

go: creating new go.mod: module url-shortener

> docker init

Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? Go                                                                 
? What version of Go do you want to use? 
1.23.4                                                                       
? What's the relative directory (with a leading .) of your main package?
.
? What port does your server listen on?
8089


✔ Created → .dockerignore
✔ Created → Dockerfile
✔ Created → compose.yaml
✔ Created → README.Docker.md

→ Your Docker files are ready!
  Review your Docker files and tailor them to your application.
  Consult README.Docker.md for information about using the generated files.

What's next?
  Start your application by running → docker compose up --build
  Your application will be available at http://localhost:8089

> go get "github.com/gorilla/mux"

go: downloading github.com/gorilla/mux v1.8.1
go: added github.com/gorilla/mux v1.8.1

> go mod tidy
go: warning: "all" matched no packages

```

