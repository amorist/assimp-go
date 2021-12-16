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
#include <assimp/mesh.h>
#include <assimp/matrix4x4.h>
#include <assimp/material.h>
#include <assimp/texture.h>
#include <assimp/metadata.h>
#include <assimp/commonMetaData.h>
#include <assimp/postprocess.h>
#include <assimp/types.h>

struct aiNode* get_child(struct aiNode* n, unsigned int index)
{
	return n->mChildren[index];
}
struct aiMesh* get_mesh(struct aiScene* s, struct aiNode* n,
	unsigned int index)
{
	return s->mMeshes[n->mMeshes[index]];
}
struct aiVector3D* mesh_vertex_at(struct aiMesh* m, unsigned int index)
{
	return &(m->mVertices[index]);
}
struct aiVector3D* mesh_normal_at(struct aiMesh* m, unsigned int index)
{
	return &(m->mNormals[index]);
}
_Bool has_tex_coords(struct aiMesh* m) {
	return m->mTextureCoords[0];
}
struct aiVector3D* mesh_texture_at(struct aiMesh* m, unsigned int index)
{
	return &(m->mTextureCoords[0][index]);
}
struct aiVector3D* mesh_tangent_at(struct aiMesh* m, unsigned int index)
{
	return &(m->mTangents[index]);
}
struct aiVector3D* mesh_bitangent_at(struct aiMesh* m, unsigned int index)
{
	return &(m->mBitangents[index]);
}
struct aiFace* get_face(struct aiMesh* m, unsigned int index)
{
	return &(m->mFaces[index]);
}
unsigned int get_face_indices(struct aiFace* f, unsigned int index)
{
	return f->mIndices[index];
}
struct aiMaterial* get_material(struct aiScene* s, struct aiMesh* m)
{
	return s->mMaterials[m->mMaterialIndex];
}

*/
import "C"
import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

// Model .
type Model struct {
	TexturesLoaded  []Texture
	Meshes          []*Mesh
	directory       string
	gammaCorrection bool
}

// Import model .
func Import(path string) (scene *C.struct_aiScene, err error) {
	now1 := time.Now()
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

	model := &Model{}
	model.directory = path[0:strings.LastIndex(path, "/")]
	textures := model.processNode(scene.mRootNode, scene)
	fmt.Println("textures", len(textures), model.directory)
	now2 := time.Now()
	fmt.Println("time", now2.Sub(now1))
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

func (model *Model) processNode(aiNode *C.struct_aiNode, aiScene *C.struct_aiScene) []Texture {
	// Process the current node
	for i := 0; i < int(aiNode.mNumMeshes); i++ {
		// Get mesh just does scene->mMeshes[node->mMeshes[i]]
		mesh := C.get_mesh(aiScene, aiNode, C.uint(i))
		model.TexturesLoaded = append(model.TexturesLoaded, model.processMesh(mesh, aiScene)...)
	}

	// Call process node on all the children nodes
	for i := 0; i < int(aiNode.mNumChildren); i++ {
		model.processNode(C.get_child(aiNode, C.uint(i)), aiScene)
	}
	return model.TexturesLoaded
}

func (model *Model) processMesh(aiMesh *C.struct_aiMesh, aiScene *C.struct_aiScene) []Texture {
	var textures []Texture

	// Process materias
	material := C.get_material(aiScene, aiMesh)

	// 1. diffuse maps
	diffuseMaps := model.loadMaterialTextures(material,
		C.aiTextureType_DIFFUSE, "texture_diffuse")
	// TODO make sure this isnt overwriting slice values
	textures = append(textures, diffuseMaps...)
	// 2. specular maps
	speculareMaps := model.loadMaterialTextures(material,
		C.aiTextureType_SPECULAR, "texture_specular")
	textures = append(textures, speculareMaps...)
	// 3. normal maps
	normalMaps := model.loadMaterialTextures(material,
		C.aiTextureType_HEIGHT, "texture_normal")
	textures = append(textures, normalMaps...)
	// 4. height maps
	heightMaps := model.loadMaterialTextures(material,
		C.aiTextureType_AMBIENT, "texture_height")
	textures = append(textures, heightMaps...)
	// 5. displacement maps
	displacementMaps := model.loadMaterialTextures(material,
		C.aiTextureType_DISPLACEMENT, "texture_displacement")
	textures = append(textures, displacementMaps...)
	// 6. color maps
	colorMaps := model.loadMaterialTextures(material,
		C.aiTextureType_EMISSIVE, "texture_color")
	textures = append(textures, colorMaps...)
	// 7. reflection maps
	reflectionMaps := model.loadMaterialTextures(material,
		C.aiTextureType_REFLECTION, "texture_reflection")
	textures = append(textures, reflectionMaps...)
	// 8. light maps
	lightMaps := model.loadMaterialTextures(material,
		C.aiTextureType_LIGHTMAP, "texture_light")
	textures = append(textures, lightMaps...)
	// 9. opacity maps
	opacityMaps := model.loadMaterialTextures(material,
		C.aiTextureType_OPACITY, "texture_opacity")
	textures = append(textures, opacityMaps...)
	return textures
}

func (model *Model) loadMaterialTextures(mat *C.struct_aiMaterial, textType uint32 /**C.enum_aiTextureType*/, typeName string) []Texture {

	var textures []Texture

	textCount := C.aiGetMaterialTextureCount(mat, textType)
	for i := uint32(0); i < uint32(textCount); i++ {
		var path C.struct_aiString

		C.aiGetMaterialTexture(
			mat,       // Material
			textType,  // Type of texture
			C.uint(i), // Index
			&path,     // Path to string
			nil,       // Texture mapping
			nil,       // UV index
			nil,       // Blend
			nil,       // Texture op
			nil,       // Map mode
			nil)       // Flags
		pathAsGoString := C.GoString(&path.data[0])

		haveLoaded := false
		for j := 0; j < len(model.TexturesLoaded); j++ {
			if model.TexturesLoaded[j].Path == pathAsGoString {
				haveLoaded = true
				break
			}
		}

		if !haveLoaded {
			var texture Texture
			texture.TextureType = typeName
			texture.Path = pathAsGoString
			textures = append(textures, texture)
			model.TexturesLoaded = append(model.TexturesLoaded, texture)
			fmt.Println(typeName, pathAsGoString)
		}
	}
	return textures
}

// TextureFromFile .
func TextureFromFile(path string, directory string, gamma bool) uint32 {
	filePath := directory + "/" + path

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	data := ImageLoad(filePath)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(data.Rect.Size().X),
		int32(data.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Set texture parameters for wrapping
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER,
		gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return textureID
}
