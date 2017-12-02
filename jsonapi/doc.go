/*
Package jsonapi is simple wrapper for buildin net/http package.
It aims to let developers build json-based web api easier.

Create an api handler is so easy:

    // HelloArgs is data structure for arguments passed by POST body.
    type HelloArgs struct {
            Name string
            Title string
    }

    // HelloReply defines data structure this api will return.
    type HelloReply struct {
            Message string
    }

    // HelloHandler greets user with hello
    func HelloHandler(dec *json.Decoder, r *http.Request, w http.ResponseWriter) (res interface{}, err error) {
            // Read json objct from request.
            var args HelloArgs
            if err = dec.Decode(&args); err != nil {
                    // The arguments are not passed in JSON format, do error handling here.
                    return
            }

            res = HelloReply{fmt.Sprintf("Hello, %s %s", args,Title, args.Name)}
            return
    }

And this is how we do in main function:

    // If you used to write http.HandleFunc("/api/hello", HelloHandler)
    http.Handle("/api/hello", jsonapi.Handler(HelloHandler))

    // Batch processing
    jsonapi.Register(myServerMux, []jsonapi.API{
            {"/api/hello", HelloHandler},
    })

*/
package jsonapi
