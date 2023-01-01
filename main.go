package main
import (
    "fmt"
    "github.com/itchyny/gojq"
)
var exp = `$username +" ::: "+ $password`
func main() {
    var1 := "$username"
    val1 := "bar"
    var2 := "$password"
    val2 := "admin"
    q, err := gojq.Parse(exp)
    if err != nil {
        panic(err)
    }
    code, err := gojq.Compile(q, gojq.WithVariables([]string{var1, var2}))
    if err != nil {
        panic(err)
    }
    iter := code.Run(nil, val1, val2)
    vv, ok := iter.Next()
    if !ok {
        panic("iter returned nothing")
    }
    fmt.Println(vv)
}