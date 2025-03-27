package composite

import (
	"fmt"
	"strings"
)

// Component 接口定义组合中所有对象的公共行为
type Component interface {
	Name() string                 // 获取组件名称
	Path() string                 // 获取组件路径
	SetParent(parent Component)   // 设置父组件
	Parent() Component            // 获取父组件
	Add(component Component)      // 添加子组件
	Remove(component Component)   // 移除子组件
	GetChild(index int) Component // 获取子组件
	Children() []Component        // 获取所有子组件
	IsComposite() bool            // 是否是组合对象
	Print(indent string)          // 打印组件信息
	Size() int                    // 获取组件大小
}

// BaseComponent 为所有组件提供基本实现
type BaseComponent struct {
	name   string
	parent Component
}

// NewBaseComponent 创建基本组件
func NewBaseComponent(name string) BaseComponent {
	return BaseComponent{name: name}
}

// Name 返回组件名称
func (b *BaseComponent) Name() string {
	return b.name
}

// Path 返回组件完整路径
func (b *BaseComponent) Path() string {
	if b.parent == nil {
		return "/" + b.name
	}
	return b.parent.Path() + "/" + b.name
}

// SetParent 设置父组件
func (b *BaseComponent) SetParent(parent Component) {
	b.parent = parent
}

// Parent 获取父组件
func (b *BaseComponent) Parent() Component {
	return b.parent
}

// 以下方法由具体子类重写
func (b *BaseComponent) Add(component Component) {
	// 默认行为：叶子节点不支持添加子组件
	fmt.Printf("无法向 %s 添加子组件：操作不支持\n", b.name)
}

func (b *BaseComponent) Remove(component Component) {
	// 默认行为：叶子节点不支持移除子组件
	fmt.Printf("无法从 %s 移除子组件：操作不支持\n", b.name)
}

func (b *BaseComponent) GetChild(index int) Component {
	// 默认行为：叶子节点没有子组件
	return nil
}

func (b *BaseComponent) Children() []Component {
	// 默认行为：叶子节点返回空数组
	return []Component{}
}

func (b *BaseComponent) IsComposite() bool {
	// 默认行为：基本组件不是组合对象
	return false
}

func (b *BaseComponent) Print(indent string) {
	// 默认行为：只打印自身名称
	fmt.Println(indent + b.name)
}

func (b *BaseComponent) Size() int {
	// 默认行为：基本组件大小为0
	return 0
}

// File 表示文件系统中的文件，是叶子节点
type File struct {
	BaseComponent
	content string
	size    int
}

// NewFile 创建新文件
func NewFile(name string, size int) *File {
	return &File{
		BaseComponent: NewBaseComponent(name),
		size:          size,
	}
}

// IsComposite 文件不是组合对象
func (f *File) IsComposite() bool {
	return false
}

// Print 打印文件信息
func (f *File) Print(indent string) {
	fmt.Printf("%s- %s (%d bytes)\n", indent, f.name, f.size)
}

// Size 返回文件大小
func (f *File) Size() int {
	return f.size
}

// SetContent 设置文件内容
func (f *File) SetContent(content string) {
	f.content = content
	// 更新文件大小
	f.size = len(content)
}

// GetContent 获取文件内容
func (f *File) GetContent() string {
	return f.content
}

// Directory 表示文件系统中的目录，是组合对象
type Directory struct {
	BaseComponent
	children []Component
}

// NewDirectory 创建新目录
func NewDirectory(name string) *Directory {
	return &Directory{
		BaseComponent: NewBaseComponent(name),
		children:      []Component{},
	}
}

// IsComposite 目录是组合对象
func (d *Directory) IsComposite() bool {
	return true
}

// Add 向目录添加子组件
func (d *Directory) Add(component Component) {
	d.children = append(d.children, component)
	component.SetParent(d)
}

// Remove 从目录移除子组件
func (d *Directory) Remove(component Component) {
	for i, child := range d.children {
		if child == component {
			d.children = append(d.children[:i], d.children[i+1:]...)
			component.SetParent(nil)
			return
		}
	}
	fmt.Printf("未找到组件 %s\n", component.Name())
}

// GetChild 获取特定索引的子组件
func (d *Directory) GetChild(index int) Component {
	if index < 0 || index >= len(d.children) {
		return nil
	}
	return d.children[index]
}

// Children 获取所有子组件
func (d *Directory) Children() []Component {
	return d.children
}

// Print 打印目录及其子组件
func (d *Directory) Print(indent string) {
	// 打印当前目录
	fmt.Printf("%s+ %s/\n", indent, d.name)

	// 打印子组件
	childIndent := indent + "  "
	for _, child := range d.children {
		child.Print(childIndent)
	}
}

// Size 计算目录总大小（包括所有子组件）
func (d *Directory) Size() int {
	total := 0
	for _, child := range d.children {
		total += child.Size()
	}
	return total
}

// Find 在目录中查找组件（支持通配符）
func (d *Directory) Find(pattern string) []Component {
	results := []Component{}

	// 递归搜索所有子组件
	for _, child := range d.children {
		// 检查当前组件是否匹配
		if strings.Contains(strings.ToLower(child.Name()), strings.ToLower(pattern)) {
			results = append(results, child)
		}

		// 如果是组合对象，继续递归搜索
		if child.IsComposite() {
			if dir, ok := child.(*Directory); ok {
				childResults := dir.Find(pattern)
				results = append(results, childResults...)
			}
		}
	}

	return results
}

// Count 统计目录中的文件和目录数量
func (d *Directory) Count() (files int, dirs int) {
	for _, child := range d.children {
		if child.IsComposite() {
			dirs++
			if dir, ok := child.(*Directory); ok {
				childFiles, childDirs := dir.Count()
				files += childFiles
				dirs += childDirs
			}
		} else {
			files++
		}
	}
	return
}
