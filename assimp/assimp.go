package assimp

/*
#cgo CFLAGS: -I/usr/local/lib/assimp/include -std=c99
#cgo LDFLAGS: -L/usr/local/lib/assimp/bin -lassimp -lstdc++

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <assimp/cimport.h>
#include <assimp/cexport.h>
#include <assimp/scene.h>
#include <assimp/material.h>
#include <assimp/postprocess.h>

struct aiMesh* mesh_at(struct aiScene* s, unsigned int index)
{
	return s->mMeshes[index];
}

struct aiMaterial* material_at(struct aiScene* s, unsigned int index)
{
	return s->mMaterials[index];
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Import model .
func Import(path string) (scene *C.struct_aiScene, err error) {
	if path == "" {
		err = fmt.Errorf("file is empty")
		return
	}
	cfilename := C.CString(path)
	defer C.free(unsafe.Pointer(cfilename))
	// C.aiProcess_GlobalScale
	scene = C.aiImportFile(cfilename, C.aiProcess_Triangulate|C.aiProcess_ValidateDataStructure|C.aiProcess_GlobalScale)
	if scene == nil {
		err = fmt.Errorf("Failed to import %s", path)
		return
	}
	if uintptr(unsafe.Pointer(scene)) == 0 {
		return nil, fmt.Errorf("Unable to load %s", path)
	}
	// make sure we have at least one mesh
	if scene.mNumMeshes < 1 {
		return nil, fmt.Errorf("%s no meshes were found", path)
	}
	return
}

// Export model.
func Export(path, outPath, format string) (err error) {
	if format == "" {
		format = "obj"
	}
	cfilename := C.CString(path)
	defer C.free(unsafe.Pointer(cfilename))
	// C.aiProcess_GlobalScale
	scene := C.aiImportFile(cfilename, C.aiProcess_ValidateDataStructure|C.aiProcess_GlobalScale)
	if scene == nil {
		err = fmt.Errorf("Failed to import %s", path)
		return
	}
	aiReturn := C.aiExportScene(scene, C.CString(format), C.CString(outPath), C.aiProcess_ValidateDataStructure)
	if aiReturn != C.aiReturn_SUCCESS {
		err = fmt.Errorf("Failed to export %s", path)
		return
	}
	C.aiReleaseImport(scene)
	return nil
}
