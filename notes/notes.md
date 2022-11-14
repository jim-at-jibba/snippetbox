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
