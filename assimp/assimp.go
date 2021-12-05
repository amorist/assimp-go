package assimp

/*
#cgo CFLAGS: -I${SRCDIR}/assimplib/include -std=c99
#cgo LDFLAGS: -L${SRCDIR}/assimplib/bin/macos -lassimp -lz -lstdc++

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <assimp/cimport.h>
#include <assimp/cexport.h>
#include <assimp/scene.h>
#include <assimp/material.h>
#include <assimp/postprocess.h>
#include <assimp/types.h>

struct aiMesh* mesh_at(struct aiScene* s, unsigned int index)
{
	return s->mMeshes[index];
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Assimp .
type Assimp struct {
}

// NewAssimp .
func NewAssimp() *Assimp {
	return &Assimp{}
}

// Import model .
func (assimp *Assimp) Import(filename string) (scene *C.struct_aiScene, err error) {
	fmt.Println("filename:", filename)
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	// C.aiProcess_GlobalScale
	scene = C.aiImportFile(cfilename, C.aiProcess_Triangulate|C.aiProcess_ValidateDataStructure|C.aiProcess_GlobalScale)
	if scene == nil {
		err = fmt.Errorf("Failed to import %s", filename)
		return
	}
	if uintptr(unsafe.Pointer(scene)) == 0 {
		return nil, fmt.Errorf("Unable to load %s", filename)
	}
	// make sure we have at least one mesh
	if scene.mNumMeshes < 1 {
		return nil, fmt.Errorf("%s no meshes were found", filename)
	}
	// for m := 0; m < int(scene.mMetaData.mNumProperties); m++ {
	// 	mMetaData := C.metadata_at(scene, C.uint(m))
	// 	fmt.Println(mMetaData)
	// }
	// get scene.mMetaData kay value
	// fmt.Println("scene.mMetaData:", scene.mMetaData)
	// fmt.Println("scene.mMetaData.mNumProperties:", scene.mMetaData.mNumProperties)
	// fmt.Println(scene.mMetaData.mKeys)

	// for m := 0; m < int(scene.mNumMeshes); m++ {

	// 	// get the current mesh
	// 	mesh := C.mesh_at(scene, C.uint(m))
	// 	fmt.Println(mesh.mVertices.x, mesh.mVertices.y, mesh.mVertices.z)
	// }
	// C.aiReleaseImport(scene)
	return
}

// Export model.
func (assimp *Assimp) Export(path string) (err error) {
	// filename := path[strings.LastIndex(path, "/")+1:]
	// // get file path with out file name
	// filepath := path[:strings.LastIndex(path, "/")+1]
	// filename = filename[:strings.LastIndex(filename, ".")]
	// fmt.Println("filename:", filename, filepath)
	scene, err := assimp.Import(path)
	if err != nil {
		return
	}
	// gltf2
	aiReturn := C.aiExportScene(scene, C.CString("obj"), C.CString("/Users/amor/Desktop/new/中文/aaa.obj"), C.aiProcess_Triangulate)
	if aiReturn != C.aiReturn_SUCCESS {
		err = fmt.Errorf("Failed to export %s", path)
		return
	}
	return nil
}
