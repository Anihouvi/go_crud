package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"

    _ "github.com/lib/pq"
    "strconv"
)

type Customer struct {
   ID int
   Name string
   Email string
   Phone string
   Address string

}

func connectDB() (*sql.DB, error) {
    host := "localhost"
    port := 5432
    user := ""
    password := ""
    dbname := "customer"

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}


func Index(w http.ResponseWriter, r *http.Request) {
    db, err := connectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Query the count of customers
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
    if err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query("SELECT id, name, email, phone, address FROM customers")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var customers []Customer
    for rows.Next() {
        var customer Customer
        err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address)
        if err != nil {
            log.Fatal(err)
        }
        customers = append(customers, customer)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    tmpl := template.Must(template.ParseFiles("index.html"))
    data := struct {
        Customers []Customer
        Count     int
    }{
        Customers: customers,
        Count:     count,
    }
    tmpl.Execute(w, data)
}


    

func Create(w http.ResponseWriter, r *http.Request) {
      db, err := connectDB()
      if err != nil {
         log.Fatal(err)     
      }
      defer db.Close()
      
      if r.Method == "POST" {
         name :=  r.FormValue("name")
         email := r.FormValue("email")
         phone := r.FormValue("phone")
         address := r.FormValue("address")

          if name == "" || email == "" || phone == "" || address == "" {
            http.Error(w, "Please fill in all the required fields", http.StatusBadRequest)
            return
        }
      
         _, err := db.Exec("INSERT INTO customers (name, email, phone, address) VALUES ($1, $2, $3, $4)", name, email, phone, address)
         if err != nil {
            log.Fatal(err)
         }
         http.Redirect(w, r, "/", http.StatusSeeOther)
         return
      }
      tmpl := template.Must(template.ParseFiles("create.html"))
      tmpl.Execute(w, nil)
     }

func Update(w http.ResponseWriter, r *http.Request) {
    db, err := connectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    id := r.URL.Query().Get("id")
    if id == "" {
        http.NotFound(w, r)
        return
    }

    if r.Method == "POST" {
        name := r.FormValue("name")
        email := r.FormValue("email")
        phone := r.FormValue("phone")
        address := r.FormValue("address")

        _, err := db.Exec("UPDATE customers SET name=$1, email=$2, phone=$3, address=$4 WHERE id=$5", name, email, phone, address, id)
        if err != nil {
            log.Fatal(err)
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    rows, err := db.Query("SELECT id, name, email, phone, address FROM customers WHERE id=$1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var customer Customer
    for rows.Next() {
        err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    tmpl := template.Must(template.ParseFiles("update.html"))
    tmpl.Execute(w, customer)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db, err := connectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    id := r.URL.Query().Get("id")
    if id == "" {
        http.NotFound(w, r)
        return
    }

    // Add this line to check the value of 'id'
    fmt.Println("ID:", id)

    if r.Method == "POST" {
        confirmation := r.FormValue("confirmation")
        if confirmation == "yes" {
            _, err = db.Exec("DELETE FROM customers WHERE id=$1", id)
            if err != nil {
                log.Fatal(err)
            }
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/delete.html"))
    var customer Customer
    customer.ID, err = strconv.Atoi(id)
    if err != nil {
        log.Fatal(err)
    }
    tmpl.Execute(w, customer)
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    query := r.URL.Query().Get("query")
    if query == "" {
        http.Error(w, "A valid input is required to search against", http.StatusBadRequest)
        return
    }

    db, err := connectDB()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    customers, err := searchCustomers(db, query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/search.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data := struct {
        Customers []Customer
        Query     string
    }{
        Customers: customers,
        Query:     query,
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


func searchCustomers(db *sql.DB, query string) ([]Customer, error) {
    rows, err := db.Query(`SELECT id, name, email FROM customers WHERE name LIKE '%' || $1 || '%' OR name LIKE '%' || $1 || '%' OR id::text = $1`, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    customers := []Customer{}
    for rows.Next() {
        var c Customer
        err := rows.Scan(&c.ID, &c.Name, &c.Email)
        if err != nil {
            return nil, err
        }
        customers = append(customers, c)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return customers, nil
}


func main() {
     // Serve static files from the 'static' directory
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    dbname := "customer"
    local_mode := true

    http.HandleFunc("/", Index)
    http.HandleFunc("/create", Create)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.HandleFunc("/search", searchHandler)
    fmt.Printf("Using database: %s\n", dbname)
    fmt.Printf("Server running in local mode: %t\n", local_mode)
    fmt.Printf("Server listening on http://localhost:8080\n")
    log.Fatal(http.ListenAndServe(":8080", nil))
    
}




      


