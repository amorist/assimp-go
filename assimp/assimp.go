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
#include <assimp/texture.h>
#include <assimp/metadata.h>
#include <assimp/commonMetaData.h>
#include <assimp/postprocess.h>
#include <assimp/types.h>

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
	// 获取贴图列表
	fmt.Println(scene.mNumMaterials)
	// get Materials from scene
	for i := 0; i < int(scene.mNumMaterials); i++ {
		mat := C.material_at(scene, C.uint(i))
		if mat == nil {
			continue
		}
		// get material name
		// var name *C.struct_aiString
		// fmt.Println(name)
		// aiReturn := C.aiGetMaterialString(mat, C.AI_MATKEY_NAME, 0, 0, &name)
		// fmt.Println(name, aiReturn)
		// var texFile *C.struct_aiMaterialProperty
		// aiReturn := C.aiGetMaterialProperty(mat, C.CString("$tex.file"), C.aiTextureType_DIFFUSE, 0, &texFile)
		// if aiReturn == C.aiReturn_SUCCESS {
		// 	fmt.Println(texFile)
		// } else {
		// 	fmt.Println("no texture")
		// }
		// defer C.free(unsafe.Pointer(texFile))
		// get diffuse texture
		// tex := C.aiGetMaterialString(mat, C.aiTextureType_DIFFUSE)
		// if tex == nil {
		// 	continue
		// }
		// fmt.Println(tex.mHeight, tex.mWidth, tex.mFilename)
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
	return nil
}
