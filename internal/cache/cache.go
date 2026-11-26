package cache

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"path/filepath"
	"slices"
	"sync"

	gocache "github.com/patrickmn/go-cache"
	"joshuanario.com/hyperfacts/pkg/hydratables"
	"joshuanario.com/hyperfacts/pkg/renderables"
	"joshuanario.com/hyperfacts/pkg/serializables"
)

var (
	dry      bool
	lock     sync.RWMutex
	once     sync.Once
	appCache *gocache.Cache
)

func NewCache(runDry bool) *gocache.Cache {
	dry = runDry
	if dry {
		return appCache
	}
	once.Do(func() {
		appCache = gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	})
	return appCache
}

func MarshalExpressable(id string, documentName string, name string, contextref string) ([]byte, error) {
	h, err := hydratable(id)
	if err != nil {
		return nil, err
	}
	hash := fnv.New128a()
	hash.Write([]byte(id + "/facts/" + name + ":" + contextref))
	cachekey := hex.EncodeToString(hash.Sum([]byte{}))
	lock.RLock()
	if !dry {
		if x, found := appCache.Get(cachekey); found {
			ret := x.([]byte)
			lock.RUnlock()
			return ret, nil
		}
	}
	lock.RUnlock()
	byteArr, err := renderables.MarshalExpressable(documentName, name, contextref, h)
	go func() {
		if dry {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(cachekey, byteArr, gocache.DefaultExpiration)
	}()
	return byteArr, err
}
func MarshalExpressableDocument(id string, documentName string) ([]byte, error) {
	ext := filepath.Ext(documentName)
	if ext == ".xhtml" || ext == ".htm" {
		lock.RLock()
		if !dry {
			if x, found := appCache.Get(id + "/" + documentName); found {
				ret := x.([]byte)
				lock.RUnlock()
				return ret, nil
			}
		}
		lock.RUnlock()
		h, err := hydratable(id)
		if err != nil {
			return nil, err
		}
		index := slices.IndexFunc(h.Folder.Documents, func(d *serializables.Document) bool {
			return d.Filename == documentName
		})
		if index < 0 {
			return nil, nil
		}
		document := h.Folder.Documents[index]
		if document.Filename == documentName {
			data := document.Bytes
			if dry {
				return data, nil
			}
			go func() {
				lock.Lock()
				defer lock.Unlock()
				appCache.Set(id+"/"+documentName, data, gocache.DefaultExpiration)
			}()
			return data, nil
		}
	}
	return nil, nil
}
func MarshalRenderable(id string, hash string) ([]byte, error) {
	lock.RLock()
	if !dry {
		if x, found := appCache.Get(id + "/" + hash); found {
			ret := x.([]byte)
			lock.RUnlock()
			return ret, nil
		}
	}
	lock.RUnlock()
	h, err := hydratable(id)
	if err != nil {
		return nil, err
	}
	byteArr, err := renderables.MarshalRenderable(hash, h)
	go func() {
		if dry {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(id+"/"+hash, byteArr, gocache.DefaultExpiration)
	}()
	return byteArr, err
}

func MarshalCatalog(id string) ([]byte, error) {
	h, err := hydratable(id)
	if err != nil {
		return nil, err
	}
	lock.RLock()
	if !dry {
		if x, found := appCache.Get(id + "/"); found {
			ret := x.([]byte)
			lock.RUnlock()
			return ret, nil
		}
	}
	lock.RUnlock()
	byteArr, err := renderables.MarshalCatalog(h)
	go func() {
		if dry {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(id+"/", byteArr, gocache.DefaultExpiration)
	}()
	return byteArr, err
}

func hydratable(id string) (*hydratables.Hydratable, error) {
	lock.RLock()
	if !dry {
		if x, found := appCache.Get(id); found {
			ret := x.(*hydratables.Hydratable)
			lock.RUnlock()
			return ret, nil
		}
	}
	lock.RUnlock()
	folder, err := serializables.Discover(id)
	if err != nil {
		return nil, fmt.Errorf("failed to discover folder, %v", err)
	}
	ret, err := hydratables.Hydrate(folder)
	go func() {
		if dry {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(id, ret, gocache.DefaultExpiration)
	}()
	return ret, err
}

func MarshalSuperRenderable(id string) ([]byte, error) {
	h, err := hydratable(id)
	if err != nil {
		return nil, err
	}
	hash := fnv.New128a()
	hash.Write([]byte(id + "/superrenderable"))
	cachekey := hex.EncodeToString(hash.Sum([]byte{}))
	lock.RLock()
	if !dry {
		if x, found := appCache.Get(cachekey); found {
			ret := x.([]byte)
			lock.RUnlock()
			return ret, nil
		}
	}
	lock.RUnlock()
	byteArr, err := renderables.MarshalSuperRenderable(h)
	go func() {
		if dry {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		appCache.Set(cachekey, byteArr, gocache.DefaultExpiration)
	}()
	return byteArr, err
}
