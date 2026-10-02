Build-from scratch doc

following below docs <br>
-> https://github.com/golang-standards/project-layout <br>
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

create a open api spec and download tool to generate code,
```
> wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.25.0/openapi-generator-cli-7.25.0.jar -O openapi-generator-cli.jar
--2026-10-02 17:40:58--  https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.25.0/openapi-generator-cli-7.25.0.jar
Resolving repo1.maven.org (repo1.maven.org)... 104.18.18.12, 104.18.19.12, 2606:4700:9c63:fa77:15c6:4cc:c262:b56
Connecting to repo1.maven.org (repo1.maven.org)|104.18.18.12|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 31942042 (30M) [application/java-archive]
Saving to: ‘openapi-generator-cli.jar’

openapi-generator-cli.jar     100%[=================================================>]  30.46M  6.46MB/s    in 5.2s

2026-10-02 17:41:03 (5.88 MB/s) - ‘openapi-generator-cli.jar’ saved [31942042/31942042]
```

created a make file to run command easily
```
>make generate-server
java -jar ./api/tool/openapi-generator-cli.jar generate \
        -i api/url-shortener.yaml \
        -g go-server \
        -o gen \
        --additional-properties=router=mux,packageName=gen,sourceFolder=""
[main] INFO  o.o.codegen.DefaultGenerator - Generating with dryRun=false
.
.
.
############################################################################################
# Thanks for using OpenAPI Generator.                                                      #
# We appreciate your support! Please consider donating to help us maintain this project.   #
# https://opencollective.com/openapi_generator/donate                                      #
############################################################################################
```

