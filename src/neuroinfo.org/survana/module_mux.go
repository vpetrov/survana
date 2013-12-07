package survana

import (
        "net/http"
        "log"
        "strings"
        "sync"
        "time"
       )

//decides which Module handles which Request
type ModuleMux struct {
    lock sync.RWMutex
    modules map[string]*Module
}

//Returns a new ModuleMux
func NewModuleMux() *ModuleMux {
    return &ModuleMux{
        modules: make(map[string]*Module, 0),
    }
}

//Finds the requested module and hands off the request to its router
func (mux *ModuleMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()

    //detect which module this request should go to
    mount_point, rurl := mux.GetModule(r.URL.Path)

    //lock modules mutex for reading (ensures that the map won't be changed
    //while we're reading from it)
    mux.lock.RLock()
        //fetch the module
        module, ok := mux.modules[mount_point]
    mux.lock.RUnlock()

    if !ok {
        //fetch the "/" module
        mux.lock.RLock()
            module, ok = mux.modules["/"]
        mux.lock.RUnlock()
        // if no default module found, return a 500 Internal Server Error
        if !ok {
            http.Error(w,
                       "Internal Server Error",
                       http.StatusInternalServerError)
            return
        }
    }

    //create an application-specific request object
    request := NewRequest(r, rurl, module)

    //route the request
    module.Router.Route(w, request)

    log.Printf("[%s] %s", time.Since(startTime).String(), r.URL.Path)
}

//Registers a new Module for a URL path
func (mux *ModuleMux) Mount(m *Module, path string) {
    //atempt to lock the modules mutex before we write to the map
    mux.lock.Lock()
        // add the module pointer to the map
        mux.modules[path] = m;
    // unlock the mutex
    mux.lock.Unlock()

    log.Println("Mounting ", m.Name, "on", path)

    m.MountPoint = path
}

//Unregisters the module that handles the path
func (mux *ModuleMux) Unmount(path string) {
    mux.lock.Lock()
        //remove the module from the path
        delete(mux.modules, path)
    mux.lock.Unlock()
}

// searches for a module by the path it's been mounted on
func (mux *ModuleMux) GetModule(path string) (module, mpath string) {

    if len(path) <= 1 {
        return "/", path
    }

    // find the second "/" in path
    slash := strings.Index(path[1:], "/") + 1

    //if no slashes were found, return original path and default module
    if slash == 0 {
        return "/", path
    }

    // the module is the first item between the 2 slashes
    module = path[:slash]

    // the module's URL 'path' is everything that follows it
    mpath = path[slash:]

    return
}

//the default module mux
var Modules *ModuleMux = NewModuleMux()
