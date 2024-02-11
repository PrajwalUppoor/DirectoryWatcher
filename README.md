# DirectoryWatcher
This is Go based application which can be used to watch multiple directories concurrently also to save the patterns which is being watched in this directory in the database for each task which is run using REST Api 
Link to the POST man collection for this API is https://www.postman.com/altimetry-engineer-7579503/workspace/watcherfordirectoryusinggo/collection/14080278-82b34288-9dd3-4667-a9f8-5d3663b2b826?action=share&creator=14080278
has the list of APIs with the route, method and sample request body along with sample responses are available in the documentation section. You can also use the curl command to hit the api.

##Prereqisites
Install Docker Desktop 
For Windows https://docs.docker.com/desktop/install/windows-install/ 
For Linux https://docs.docker.com/desktop/install/linux-install/
For Mac https://docs.docker.com/desktop/install/mac-install/
Install Latest Go version 1.20.2 by following link https://go.dev/doc/install
Open terminal and run below command 
```
go version
```
It should show version as follows for windows 
```
go version go1.20.2 windows/amd64
```
Then Start the Docker Services
Go to the Root of the Repo and run the below command 
```
docker compose -f "docker-compose.yml" up -d --build
```
run the below command to verify that postgresql is running at port number 5432
``` 
docker ps
``` 
if you are running something at 5432 just change the port number in the docker-compose.yml or kill that process which is listening to 5432 port
Once your postgres Sql is running succesfully 
run the below command 
```
go run .\main.go
``` 
Applications starts running at port number 8080 you can change this in main.go.

Use Postman Collection for testing this Api 
For Reference of Database Schema open DirectoryWatcherERDiagram.png in root of the repo 

This Application uses Gin Framework along with GORM



###Debugging
```
docker exec -it docker-containerid sh
```
Once the shell opens run the below command
```
psql -U postgres
```
select the database by following command
```
postgres=# \c postgres
```
In order to list the relations
```  
postgres=# \d 
```
Inorder to describe individual tables
``` 
postgres=# \d configurations
```
```
postgres=# \d tasks
```
