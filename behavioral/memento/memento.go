package memento

import "fmt"

// Memento 定义备忘录接口
// 在备忘录模式中，该接口用于保存对象的内部状态
type Memento interface {
	// GetState 是一个空方法，仅用于限制访问
	// 真正的状态获取方法只对Originator可见
}

// documentMemento 具体的备忘录实现，保存文档的完整状态
type documentMemento struct {
	title string
	body  string
}

// 确保documentMemento实现了Memento接口
var _ Memento = &documentMemento{}

// GetState 是Memento接口要求的方法
func (dm *documentMemento) GetState() {
	// 这是一个空方法，仅用于接口实现
}

// Document 定义文档结构体（Originator）
// 在备忘录模式中，Originator可以创建备忘录并从中恢复状态
type Document struct {
	title string // 文档标题
	body  string // 文档内容
}

// Title 获取文档标题
func (d *Document) Title() string {
	return d.title
}

// SetTitle 设置文档标题
func (d *Document) SetTitle(title string) {
	d.title = title
}

// Body 获取文档内容
func (d *Document) Body() string {
	return d.body
}

// SetBody 设置文档内容
func (d *Document) SetBody(body string) {
	d.body = body
}

// NewDocument 创建新的文档对象实例
func NewDocument(title string) *Document {
	return &Document{title: title}
}

// CreateMemento 创建文档的备忘录
func (d *Document) CreateMemento() Memento {
	// 返回包含完整文档状态的备忘录
	return &documentMemento{
		title: d.title,
		body:  d.body,
	}
}

// RestoreFromMemento 从备忘录恢复文档状态
func (d *Document) RestoreFromMemento(m Memento) {
	// 类型断言确保传入的是正确的备忘录类型
	if memento, ok := m.(*documentMemento); ok {
		d.title = memento.title
		d.body = memento.body
	}
}

// Caretaker 定义备忘录管理者结构体
// 在备忘录模式中，Caretaker负责保存备忘录对象，但不会操作或检查备忘录的内容
type Caretaker struct {
	document    *Document // 当前编辑的文档
	mementos    []Memento // 历史备忘录数组
	currentPos  int       // 当前备忘录位置
	maxMementos int       // 最大保存的备忘录数量
}

// NewCaretaker 创建新的备忘录管理者实例
func NewCaretaker(document *Document, maxMementos int) *Caretaker {
	fmt.Println("打开文档：" + document.Title())

	if maxMementos <= 0 {
		maxMementos = 10 // 默认保存10个备忘录
	}

	c := &Caretaker{
		document:    document,
		mementos:    make([]Memento, 0),
		currentPos:  -1,
		maxMementos: maxMementos,
	}

	// 保存初始状态
	c.SaveState()
	return c
}

// SaveState 保存当前文档状态到备忘录
func (c *Caretaker) SaveState() {
	// 如果当前不是在最新状态后面添加，需要删除当前位置之后的所有历史
	if c.currentPos < len(c.mementos)-1 {
		c.mementos = c.mementos[:c.currentPos+1]
	}

	// 添加新的备忘录
	c.mementos = append(c.mementos, c.document.CreateMemento())
	c.currentPos = len(c.mementos) - 1

	// 如果历史记录超过上限，删除最早的记录
	if len(c.mementos) > c.maxMementos {
		c.mementos = c.mementos[1:]
		c.currentPos--
	}
}

// Undo 撤销操作，恢复到上一个状态
func (c *Caretaker) Undo() bool {
	if c.currentPos <= 0 {
		fmt.Println("已经是最早的状态，无法撤销")
		return false
	}

	c.currentPos--
	c.document.RestoreFromMemento(c.mementos[c.currentPos])

	fmt.Println("===> 撤销操作，文档内容如下：")
	c.ShowDocument()
	return true
}

// Redo 重做操作，恢复到下一个状态
func (c *Caretaker) Redo() bool {
	if c.currentPos >= len(c.mementos)-1 {
		fmt.Println("已经是最新的状态，无法重做")
		return false
	}

	c.currentPos++
	c.document.RestoreFromMemento(c.mementos[c.currentPos])

	fmt.Println("===> 重做操作，文档内容如下：")
	c.ShowDocument()
	return true
}

// Append 在文档末尾添加文本内容
func (c *Caretaker) Append(text string) {
	c.document.SetBody(c.document.Body() + text + "\n")
	c.SaveState()

	fmt.Println("===> 插入操作，文档内容如下：")
	c.ShowDocument()
}

// Delete 删除文档内容
func (c *Caretaker) Delete() {
	c.document.SetBody("")
	c.SaveState()

	fmt.Println("===> 删除操作，文档内容如下：")
	c.ShowDocument()
}

// ShowDocument 显示文档内容
func (c *Caretaker) ShowDocument() {
	fmt.Printf("标题: %s\n", c.document.Title())
	fmt.Printf("内容:\n%s\n", c.document.Body())
}

// Save 保存文档内容
func (c *Caretaker) Save() {
	fmt.Println("===> 存盘操作")
}

// CanUndo 检查是否可以撤销
func (c *Caretaker) CanUndo() bool {
	return c.currentPos > 0
}

// CanRedo 检查是否可以重做
func (c *Caretaker) CanRedo() bool {
	return c.currentPos < len(c.mementos)-1
}
