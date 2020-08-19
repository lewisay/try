// Package try provides retry functionality.
//     err := try.Do(context.TODO(), func(attempt int) (retry bool, err error) {
//       retry = attempt < 3 // try 3 times
//       err = doSomeThing()
//       return retry, err
//     })
//     if err != nil {
//       log.Fatalln("error:", err)
//     }
package documentation
