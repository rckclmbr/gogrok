package grok

/*
#cgo LDFLAGS: -lgrok
#include <grok.h>
*/
import "C"
import (
	"unsafe"
	"sync"
)

type Grok struct {
	grok *C.grok_t
	lock sync.Mutex
}

// MyError is an error implementation that includes a time and message.
type GrokError struct {
	GrokErrorCode	C.int
}

func (e GrokError) Error() string {
	
	errors := []string{
		"GROK_OK",
		"GROK_ERROR_FILE_NOT_ACCESSIBLE",
		"GROK_ERROR_PATTERN_NOT_FOUND",
		"GROK_ERROR_UNEXPECTED_READ_SIZE",
		"GROK_ERROR_COMPILE_FAILED",
		"GROK_ERROR_UNINITIALIZED",
		"GROK_ERROR_PCRE_ERROR",
		"GROK_ERROR_NOMATCH",
	}

	return errors[e.GrokErrorCode]
}

func New() (*Grok) {
	grok := C.grok_new()
	return &Grok{grok:grok}
}

func (grok *Grok) Cleanup() {
	C.grok_free(grok.grok)
}

func (grok *Grok) Compile(str string) (error) {

	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	ret := C.grok_compile(grok.grok, cstr)
	if ret != C.GROK_OK {
		return GrokError{ret}
	}
	return nil
}

func (grok *Grok) Match(str string) (map[string]string, error) {


	var gm C.grok_match_t
	//log.Print("Attemping to match '", str, "'")
	
	grok.lock.Lock()
	ret := C.grok_exec(grok.grok, C.CString(str), &gm)
	grok.lock.Unlock()
	if ret != C.GROK_OK {
		return nil, GrokError{ret}
	}
	
	var pdata, pname *C.char
	var pname_len, pdata_len C.int
	
	result := make(map[string]string)

	grok.lock.Lock()	
	C.grok_match_walk_init(&gm)
	for C.grok_match_walk_next(&gm, &pname, &pname_len, &pdata, &pdata_len) == 0 {
		name := C.GoStringN(pname, pname_len)
		data := C.GoStringN(pdata, pdata_len)
		
		result[name] = data
	}
	C.grok_match_walk_end(&gm)
	grok.lock.Unlock()
	return result, nil

	//grok_match_t
}

func (grok *Grok) AddPattern(word string, regexp string) (error) {
	//const char *regexp = NULL;
	
	word_len := C.size_t(len(word))
	regexp_len := C.size_t(len(regexp))
	
	cword := C.CString(word)
	defer C.free(unsafe.Pointer(cword))
	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))
	
	ret := C.grok_pattern_add(grok.grok, cword, word_len, cregexp, regexp_len) 
	if ret != C.GROK_OK {
		return GrokError{ret}
	}
	return nil
}

func (grok *Grok) AddPatternsFromFile(filename string) (error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	ret := C.grok_patterns_import_from_file(grok.grok, cfilename)
	if ret != C.GROK_OK {
		return GrokError{ret}
	}
	return nil
}
