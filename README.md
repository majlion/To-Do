# To-Do
To-Do application

Requirements:
go mod init todo-app
go get -u github.com/gorilla/mux
go get -u github.com/lib/pq


api:
GET http://localhost:8000/tasks


POST http://localhost:8000/tasks
Body: {"title": "My Task", "complete": false}


PUT http://localhost:8000/tasks/{id}
Body: {"title": "Updated Task", "complete": true}


DELETE http://localhost:8000/tasks/{id}

