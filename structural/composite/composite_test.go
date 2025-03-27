package composite

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv" // 添加strconv包
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 捕获标准输出的辅助函数
func captureOutput(f func()) string {
	// 保存原始的标准输出
	oldStdout := os.Stdout

	// 创建一个管道
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 执行函数
	f()

	// 恢复原始的标准输出
	os.Stdout = oldStdout
	w.Close()

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// 测试文件基本功能
func TestFile(t *testing.T) {
	t.Run("File basic properties", func(t *testing.T) {
		assert := assert.New(t)

		file := NewFile("document.txt", 100)
		assert.Equal("document.txt", file.Name())
		assert.Equal(100, file.Size())
		assert.False(file.IsComposite())

		// 测试设置内容
		file.SetContent("Hello, World!")
		assert.Equal("Hello, World!", file.GetContent())
		assert.Equal(13, file.Size()) // 内容长度应为13个字符
	})

	t.Run("File cannot add children", func(t *testing.T) {
		file := NewFile("document.txt", 100)
		childFile := NewFile("child.txt", 50)

		output := captureOutput(func() {
			file.Add(childFile)
		})

		assert.Contains(t, output, "无法向 document.txt 添加子组件")
	})
}

// 测试目录基本功能
func TestDirectory(t *testing.T) {
	t.Run("Directory basic properties", func(t *testing.T) {
		assert := assert.New(t)

		dir := NewDirectory("projects")
		assert.Equal("projects", dir.Name())
		assert.Equal(0, dir.Size()) // 空目录大小为0
		assert.True(dir.IsComposite())
		assert.Empty(dir.Children())
	})

	t.Run("Directory add and remove children", func(t *testing.T) {
		assert := assert.New(t)

		dir := NewDirectory("projects")
		file1 := NewFile("file1.txt", 100)
		file2 := NewFile("file2.txt", 200)

		// 添加文件
		dir.Add(file1)
		dir.Add(file2)

		assert.Len(dir.Children(), 2)
		assert.Equal(file1, dir.GetChild(0))
		assert.Equal(file2, dir.GetChild(1))

		// 移除文件
		dir.Remove(file1)
		assert.Len(dir.Children(), 1)
		assert.Equal(file2, dir.GetChild(0))

		// 尝试访问无效索引
		assert.Nil(dir.GetChild(99))
	})

	t.Run("Directory parent-child relationship", func(t *testing.T) {
		assert := assert.New(t)

		root := NewDirectory("root")
		home := NewDirectory("home")
		user := NewDirectory("user")
		file := NewFile("profile.txt", 150)

		root.Add(home)
		home.Add(user)
		user.Add(file)

		// 检查父子关系
		assert.Equal(root, home.Parent())
		assert.Equal(home, user.Parent())
		assert.Equal(user, file.Parent())

		// 检查路径
		assert.Equal("/root", root.Path())
		assert.Equal("/root/home", home.Path())
		assert.Equal("/root/home/user", user.Path())
		assert.Equal("/root/home/user/profile.txt", file.Path())
	})
}

// 测试目录大小计算
func TestDirectorySize(t *testing.T) {
	assert := assert.New(t)

	// 创建一个复杂的目录结构
	root := NewDirectory("root")

	docs := NewDirectory("documents")
	root.Add(docs)

	file1 := NewFile("doc1.txt", 100)
	file2 := NewFile("doc2.txt", 200)
	docs.Add(file1)
	docs.Add(file2)

	photos := NewDirectory("photos")
	root.Add(photos)

	photo1 := NewFile("photo1.jpg", 1000)
	photo2 := NewFile("photo2.jpg", 2000)
	photos.Add(photo1)
	photos.Add(photo2)

	// 验证递归大小计算
	assert.Equal(300, docs.Size())
	assert.Equal(3000, photos.Size())
	assert.Equal(3300, root.Size())

	// 测试修改文件大小后的更新
	file1.SetContent("This is a longer content that should increase file size")
	assert.Equal(55, file1.Size())
	assert.Equal(255, docs.Size())
	assert.Equal(3255, root.Size())
}

// 测试查找功能
func TestDirectoryFind(t *testing.T) {
	assert := assert.New(t)

	// 创建示例文件系统
	root := NewDirectory("root")

	docs := NewDirectory("documents")
	root.Add(docs)

	report := NewFile("report.doc", 150)
	letter := NewFile("letter.doc", 75)
	notes := NewFile("notes.txt", 50)
	docs.Add(report)
	docs.Add(letter)
	docs.Add(notes)

	photos := NewDirectory("photos")
	root.Add(photos)

	vacation := NewDirectory("vacation")
	photos.Add(vacation)

	beach := NewFile("beach.jpg", 1000)
	sunset := NewFile("sunset.jpg", 1200)
	vacation.Add(beach)
	vacation.Add(sunset)

	// 测试查找功能
	results := root.Find("doc")
	assert.Len(results, 3) // documents目录和两个.doc文件

	results = root.Find("jpg")
	assert.Len(results, 2)

	results = root.Find("vacation")
	assert.Len(results, 1)

	results = root.Find("nonexistent")
	assert.Empty(results)
}

// 测试目录计数功能
func TestDirectoryCount(t *testing.T) {
	assert := assert.New(t)

	// 创建示例文件系统
	root := NewDirectory("root")

	docs := NewDirectory("documents")
	photos := NewDirectory("photos")
	music := NewDirectory("music")
	root.Add(docs)
	root.Add(photos)
	root.Add(music)

	docs.Add(NewFile("doc1.txt", 100))
	docs.Add(NewFile("doc2.txt", 200))

	photos.Add(NewFile("photo1.jpg", 1000))

	vacationDir := NewDirectory("vacation")
	photos.Add(vacationDir)
	vacationDir.Add(NewFile("beach.jpg", 2000))
	vacationDir.Add(NewFile("mountain.jpg", 1500))

	// 验证计数
	files, dirs := root.Count()
	assert.Equal(5, files) // 总共5个文件
	assert.Equal(4, dirs)  // 不含root的4个目录

	files, dirs = docs.Count()
	assert.Equal(2, files)
	assert.Equal(0, dirs)

	files, dirs = photos.Count()
	assert.Equal(3, files)
	assert.Equal(1, dirs)
}

// 测试打印功能
func TestPrint(t *testing.T) {
	t.Run("Print file", func(t *testing.T) {
		file := NewFile("document.txt", 100)

		output := captureOutput(func() {
			file.Print("")
		})

		assert.Contains(t, output, "- document.txt (100 bytes)")
	})

	t.Run("Print directory structure", func(t *testing.T) {
		// 创建一个简单的目录结构
		root := NewDirectory("root")
		docs := NewDirectory("documents")
		root.Add(docs)
		docs.Add(NewFile("report.doc", 150))
		docs.Add(NewFile("letter.txt", 75))
		root.Add(NewFile("config.json", 30))

		output := captureOutput(func() {
			root.Print("")
		})

		// 验证输出包含预期的格式和内容
		assert.Contains(t, output, "+ root/")
		assert.Contains(t, output, "  + documents/")
		assert.Contains(t, output, "    - report.doc (150 bytes)")
		assert.Contains(t, output, "    - letter.txt (75 bytes)")
		assert.Contains(t, output, "  - config.json (30 bytes)")
	})
}

// 测试移除功能
func TestRemove(t *testing.T) {
	t.Run("Remove nonexistent component", func(t *testing.T) {
		dir := NewDirectory("home")
		nonexistentFile := NewFile("nonexistent.txt", 100)

		output := captureOutput(func() {
			dir.Remove(nonexistentFile)
		})

		assert.Contains(t, output, "未找到组件")
	})

	t.Run("Remove from file is not allowed", func(t *testing.T) {
		file := NewFile("document.txt", 100)
		childFile := NewFile("child.txt", 50)

		output := captureOutput(func() {
			file.Remove(childFile)
		})

		assert.Contains(t, output, "无法从 document.txt 移除子组件")
	})
}

// 示例测试：展示组合模式的典型使用
func ExampleDirectory_Print() {
	// 创建文件系统结构
	root := NewDirectory("home")

	bob := NewDirectory("bob")
	root.Add(bob)

	alice := NewDirectory("alice")
	root.Add(alice)

	// 添加Bob的文件
	bobDocs := NewDirectory("documents")
	bob.Add(bobDocs)
	bobDocs.Add(NewFile("resume.pdf", 1024))
	bobDocs.Add(NewFile("photo.jpg", 2048))

	// 添加Alice的文件
	alice.Add(NewFile("notes.txt", 256))
	aliceProjects := NewDirectory("projects")
	alice.Add(aliceProjects)
	aliceProjects.Add(NewFile("project1.doc", 4096))

	// 打印整个结构
	root.Print("")

	// Output:
	// + home/
	//   + bob/
	//     + documents/
	//       - resume.pdf (1024 bytes)
	//       - photo.jpg (2048 bytes)
	//   + alice/
	//     - notes.txt (256 bytes)
	//     + projects/
	//       - project1.doc (4096 bytes)
}

func ExampleDirectory_Find() {
	// 创建一个含有各种文档的文件系统
	root := NewDirectory("docs")

	technical := NewDirectory("technical")
	root.Add(technical)
	technical.Add(NewFile("manual.pdf", 2048))
	technical.Add(NewFile("guide.pdf", 1024))

	personal := NewDirectory("personal")
	root.Add(personal)
	personal.Add(NewFile("resume.docx", 512))
	personal.Add(NewFile("photo.jpg", 3072))

	// 查找所有PDF文件
	results := root.Find("pdf")

	// 打印结果
	for _, file := range results {
		fmt.Println(file.Name(), "-", file.Path())
	}

	// Output:
	// manual.pdf - /docs/technical/manual.pdf
	// guide.pdf - /docs/technical/guide.pdf
}

func ExampleDirectory_Size() {
	// 创建一个文件系统
	root := NewDirectory("projects")

	projectA := NewDirectory("projectA")
	root.Add(projectA)
	projectA.Add(NewFile("main.go", 1024))
	projectA.Add(NewFile("README.md", 256))

	projectB := NewDirectory("projectB")
	root.Add(projectB)
	projectB.Add(NewFile("index.html", 512))
	projectB.Add(NewFile("styles.css", 128))

	// 计算并打印每个目录的大小
	fmt.Printf("projectA大小: %d bytes\n", projectA.Size())
	fmt.Printf("projectB大小: %d bytes\n", projectB.Size())
	fmt.Printf("总大小: %d bytes\n", root.Size())

	// Output:
	// projectA大小: 1280 bytes
	// projectB大小: 640 bytes
	// 总大小: 1920 bytes
}

// 测试边缘情况
func TestEdgeCases(t *testing.T) {
	t.Run("Empty directory", func(t *testing.T) {
		assert := assert.New(t)
		dir := NewDirectory("empty")

		assert.Equal(0, dir.Size())
		assert.Empty(dir.Children())

		files, dirs := dir.Count()
		assert.Equal(0, files)
		assert.Equal(0, dirs)

		results := dir.Find("anything")
		assert.Empty(results)
	})

	t.Run("Deep nesting", func(t *testing.T) {
		assert := assert.New(t)

		// 创建一个深层嵌套的目录结构
		root := NewDirectory("level1")
		current := root

		// 创建10层深的结构
		for i := 2; i <= 10; i++ {
			// 修正这一行，使用正确的整数到字符串的转换
			next := NewDirectory("level" + strconv.Itoa(i))
			current.Add(next)
			current = next
		}

		// 在最深层添加一个文件
		deepFile := NewFile("deep.txt", 100)
		current.Add(deepFile)

		// 检查路径是否正确
		path := deepFile.Path()
		assert.Equal(11, strings.Count(path, "/")) // 路径中应有11个斜杠
		assert.True(strings.HasSuffix(path, "deep.txt"))

		// 目录应只有一个文件和9个目录
		files, dirs := root.Count()
		assert.Equal(1, files)
		assert.Equal(9, dirs)
	})
}
