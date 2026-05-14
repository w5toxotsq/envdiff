// Package parser provides functionality for reading and parsing .env files
// into key-value maps.
//
// Supported formats:
//
//	KEY=VALUE
//	KEY="VALUE"
//	KEY='VALUE'
//
// Lines beginning with '#' are treated as comments and ignored.
// Blank lines are also ignored.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["DB_HOST"])
package parser
