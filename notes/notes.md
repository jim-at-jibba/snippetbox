# Notes

## **fix path** and **subtree paths**

- fix path dont end in a slash `/snippets/view` etc
- subtree paths match when the _start_ of the path match. A bit like a wild card

## Default serveMux()

It is possible to not create a `NewServeMux` and just do `HandleFunc`. This creates a `DefaultServeMux` which is global. This is not recommended to prod applications as other 3rd party packages would also be able to add routes to your app. If they are compromised this would be bad

## Writing headers

- Its only possible to call `w.WriteHeader` once per reponse
- If not called explicitly, the first `w.Write` will send 200 OK

## Application organisation

- `cmd` - application specific code for the execuable of the project
- `internal` - non application specific code like validation and sql models
- `ui` for the UI assets

## The http.Handler interface

- a handler is an object that satifies the `http.Handler` interface, basically the object must have the `ServeHTTP()` method

```go
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  [..]
}
```

- This is really long winded and so we make use of some built in methods that allow us to reduce the amount we type

```go
  // this is syntactic sugar that transforms  a functions to a handler
	mux.HandleFunc("/snippet/view", snippetView)
```

- Because `serveMux` also has `ServeHTTP()` it also satifies the interface so it can help to thnk of it as a _special kind of handler_
- All requests are handled concurrently. This makes it super fast but you will need to guard against race conditions
