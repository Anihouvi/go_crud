
# go_crud
Basic CRUD in Go Lang (Create , Update, Delete) and  search functionality

# Set up your environment

Make sure that you have your Go environment ready and  use the Go command to download all the packages in the .mod file 


## Use go get command  to install  packages

`go get <package_name>`

`go run main.go`



## Customer Table
Use the following  script to creat the table

```
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT
);

 ```


#### Schema


| id | email | name |   phone    |   address
| :-------- | :------- | :------| :------| :------|
| `integer` | `string` |  `string`|`string`| `string`   





## Acknowledgements

 - [All go packages](https://pkg.go.dev/awesome-README-templates)
 - [More about CRUD Operations](https://stackify.com/what-are-crud-operations/)
 - [Build a web server in Go lang](https://blog.logrocket.com/creating-a-web-server-with-golang/)




## 🚀 About Me
I'm a passionate about solving problems and writing interesting codes.Introvert and a Jesus believer. Follow me for more interesting projects. 

