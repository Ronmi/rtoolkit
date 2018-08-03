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
    func HelloHandler(q jsonapi.Request) (res interface{}, err error) {
            // Read json objct from request.
            var args HelloArgs
            if err = q.Decode(&args); err != nil {
                    // The arguments are not passed in JSON format, do error
                    // handling here.
                    return
            }

            res = HelloReply{fmt.Sprintf("Hello, %s %s", args,Title, args.Name)}
            return
    }

And this is how we do in main function:

    // Suggested usage
    apis := []jsonapi.API{
        {"/api/hello", HelloHandler},
    }
    jsonapi.Register(http.DefaultMux, apis)

    // old-school
    http.Handle("/api/hello", jsonapi.Handler(HelloHandler))


Call API with TypeScript

There's a `fetch.ts` providing `grab<T>()` as simple wrapping around `fetch()`.
With following Go code:

    type MyStruct struct {
        X int  `json:"x"`
    	Y bool `json:"y"
    }

    func MyAPI(q jsonapi.Request) (ret interface{}, err error) {
        return []MyStruct{
    	    {X: 1, Y: true},
    		{X: 2},
    	}, nil
    }

    function main() {
        apis := []jsonapi.API{
    	    {"/my-api", MyAPI},
        }
    	jsonapi.Register(http.DefaultMux, apis)
    	http.ListenAndServe(":80", nil)
    }

You might write TypeScript code like this:

    export interface MyStruct {
      x?: number;
      y?: boolean;
    }

    export function getMyApi(): Promise<MyStruct[]> {
      return grab<MyStruct[]>('/my-api');
    }

    export function postMyApi(): Promise<MyStruct[]> {
      return grab<MyStruct[]>('/my-api', {
        method: 'POST',
    	headers: {'Content-Type': 'application/json'},
    	body: JSON.stringify('my data')
      });
    }

*/
package jsonapi
