# TriFont
A file format for storing fonts as 2d meshes.  
This reposetory contains a Go implementation of the format.  
TriFont files can be stored with the .trif extention.  
# The Format
The first 2 bytes specify the amount of characters.  
The next bytes specify the utf32 values of the included characters.  

**The rest of the file will be segments consisting of:**  
- 2 bytes specifying the amount of vertices.
- The vertex positions as 2 float32s per vertex.  
- 4 bytes specifying the amount of indices.  
- The indices as a uint16 per index.
- A float32 specifying how far to advance

The Vertex Y positions should be between 0 and 1.  
Little Endian should be used.  