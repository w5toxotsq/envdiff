// Package loader handles reading .env files from the filesystem.
//
// It wraps the parser package with additional validation steps, such as
// checking that the target file exists and that its name follows the
// conventional .env naming scheme (e.g. .env, .env.local, .env.production).
//
// Basic usage:
//
//	env, err := loader.Load(".env.staging", loader.Options{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// To bypass the extension check (useful in tests or pipelines where files may
// have arbitrary names):
//
//	env, err := loader.Load("custom-config", loader.Options{
//	    SkipExtensionCheck: true,
//	})
package loader
