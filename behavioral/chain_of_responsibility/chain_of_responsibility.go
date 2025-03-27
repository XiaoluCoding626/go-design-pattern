package chain_of_responsibility

import (
	"fmt"
)

// IApprover 接口定义了审批人的方法
type IApprover interface {
	// SetNext 设置下一个审批人
	SetNext(approver IApprover) IApprover
	// Approve 批准金额
	Approve(amount float64) ApprovalResult
	// GetName 获取审批人名称
	GetName() string
}

// ApprovalResult 表示审批结果
type ApprovalResult struct {
	Approved bool   // 是否批准
	Approver string // 审批人
	Message  string // 审批消息
}

// BaseApprover 结构体表示一个基础审批人
type BaseApprover struct {
	name  string
	limit float64
	next  IApprover
}

// NewBaseApprover 创建一个新的基础审批人
func NewBaseApprover(name string, limit float64) *BaseApprover {
	return &BaseApprover{
		name:  name,
		limit: limit,
	}
}

// SetNext 设置下一个审批人
func (a *BaseApprover) SetNext(approver IApprover) IApprover {
	a.next = approver
	return approver
}

// GetName 获取审批人名称
func (a *BaseApprover) GetName() string {
	return a.name
}

// Approve 基础的审批逻辑
func (a *BaseApprover) Approve(amount float64) ApprovalResult {
	// 基类不实现具体逻辑，由子类重写
	return ApprovalResult{
		Approved: false,
		Approver: a.name,
		Message:  "未实现审批逻辑",
	}
}

// TryNext 尝试传递给下一个处理者
func (a *BaseApprover) TryNext(amount float64) ApprovalResult {
	if a.next != nil {
		return a.next.Approve(amount)
	}

	return ApprovalResult{
		Approved: false,
		Approver: a.name,
		Message:  fmt.Sprintf("请求金额 %.2f 超出了所有审批人的权限范围", amount),
	}
}

// Manager 结构体表示经理职级的审批人
type Manager struct {
	*BaseApprover
}

// NewManager 创建一个新的经理审批人
func NewManager(limit float64) *Manager {
	return &Manager{BaseApprover: NewBaseApprover("经理", limit)}
}

// Approve 实现经理审批金额的逻辑
func (m *Manager) Approve(amount float64) ApprovalResult {
	if amount <= m.limit {
		return ApprovalResult{
			Approved: true,
			Approver: m.GetName(),
			Message:  fmt.Sprintf("%s批准了 %.2f 元的请求", m.GetName(), amount),
		}
	}

	return m.TryNext(amount)
}

// Director 结构体表示总监职级的审批人
type Director struct {
	*BaseApprover
}

// NewDirector 创建一个新的总监审批人
func NewDirector(limit float64) *Director {
	return &Director{BaseApprover: NewBaseApprover("总监", limit)}
}

// Approve 实现总监审批金额的逻辑
func (d *Director) Approve(amount float64) ApprovalResult {
	if amount <= d.limit {
		return ApprovalResult{
			Approved: true,
			Approver: d.GetName(),
			Message:  fmt.Sprintf("%s批准了 %.2f 元的请求", d.GetName(), amount),
		}
	}

	return d.TryNext(amount)
}

// CFO 结构体表示CFO职级的审批人
type CFO struct {
	*BaseApprover
}

// NewCFO 创建一个新的CFO审批人
func NewCFO(limit float64) *CFO {
	return &CFO{BaseApprover: NewBaseApprover("CFO", limit)}
}

// Approve 实现CFO审批金额的逻辑
func (c *CFO) Approve(amount float64) ApprovalResult {
	if amount <= c.limit {
		return ApprovalResult{
			Approved: true,
			Approver: c.GetName(),
			Message:  fmt.Sprintf("%s批准了 %.2f 元的请求", c.GetName(), amount),
		}
	}

	// 如果有下一个处理者，尝试传递，否则返回拒绝消息
	if c.next != nil {
		return c.TryNext(amount)
	}

	return ApprovalResult{
		Approved: false,
		Approver: c.GetName(),
		Message:  fmt.Sprintf("请求金额 %.2f 超出了%s的审批权限，无法批准", amount, c.GetName()),
	}
}

// CreateApprovalChain 创建一个完整的审批责任链
func CreateApprovalChain() IApprover {
	manager := NewManager(1000)   // 经理可审批1000元以下
	director := NewDirector(5000) // 总监可审批5000元以下
	cfo := NewCFO(20000)          // CFO可审批20000元以下

	// 构建责任链
	manager.SetNext(director).SetNext(cfo)

	return manager // 返回责任链的首个处理者
}
