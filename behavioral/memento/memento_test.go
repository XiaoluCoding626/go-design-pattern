package memento

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMementoPattern(t *testing.T) {
	// 创建文档
	document := NewDocument("设计模式示例")

	// 创建备忘录管理者
	editor := NewCaretaker(document, 10)

	// 进行一系列操作
	editor.Append("这是第一行内容。")
	editor.Append("这是第二行内容。")
	editor.Append("这是第三行内容。")

	// 测试撤销
	assert.True(t, editor.Undo(), "应该可以撤销")

	// 测试重做
	assert.True(t, editor.Redo(), "应该可以重做")

	// 测试在撤销后添加新内容
	editor.Undo()
	editor.Append("这是新的第三行内容。")

	// 此时应该无法重做，因为历史已分叉
	assert.False(t, editor.CanRedo(), "不应该可以重做")

	// 测试删除操作
	editor.Delete()
	assert.Empty(t, document.Body(), "删除操作后文档应为空")

	// 测试撤销删除操作
	editor.Undo()
	assert.NotEmpty(t, document.Body(), "撤销删除操作后文档不应为空")
}
