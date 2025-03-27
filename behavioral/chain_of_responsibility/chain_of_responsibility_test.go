package chainofresponsibility

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试各级审批人的审批功能
func TestApproverLevels(t *testing.T) {
	manager := NewManager(1000)
	director := NewDirector(5000)
	cfo := NewCFO(20000)

	tests := []struct {
		name         string
		approver     IApprover
		amount       float64
		wantApproved bool
		wantApprover string
	}{
		{
			name:         "经理批准小额请求",
			approver:     manager,
			amount:       800,
			wantApproved: true,
			wantApprover: "经理",
		},
		{
			name:         "经理无法批准超出权限的请求",
			approver:     manager,
			amount:       2000,
			wantApproved: false,
			wantApprover: "经理", // 因为没有下一个处理者，所以还是经理
		},
		{
			name:         "总监批准中等额度请求",
			approver:     director,
			amount:       3000,
			wantApproved: true,
			wantApprover: "总监",
		},
		{
			name:         "总监无法批准超出权限的请求",
			approver:     director,
			amount:       8000,
			wantApproved: false,
			wantApprover: "总监",
		},
		{
			name:         "CFO批准大额请求",
			approver:     cfo,
			amount:       18000,
			wantApproved: true,
			wantApprover: "CFO",
		},
		{
			name:         "CFO无法批准超出权限的请求",
			approver:     cfo,
			amount:       25000,
			wantApproved: false,
			wantApprover: "CFO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.approver.Approve(tt.amount)

			assert := assert.New(t)
			assert.Equal(tt.wantApproved, result.Approved, "批准状态错误")
			assert.Equal(tt.wantApprover, result.Approver, "审批人错误")
			t.Logf("测试结果: %s", result.Message)
		})
	}
}

// 测试责任链的正确传递
func TestApprovalChain(t *testing.T) {
	// 创建责任链
	manager := NewManager(1000)
	director := NewDirector(5000)
	cfo := NewCFO(20000)

	manager.SetNext(director)
	director.SetNext(cfo)

	tests := []struct {
		name         string
		amount       float64
		wantApproved bool
		wantApprover string
	}{
		{
			name:         "500元由经理审批",
			amount:       500,
			wantApproved: true,
			wantApprover: "经理",
		},
		{
			name:         "3000元由总监审批",
			amount:       3000,
			wantApproved: true,
			wantApprover: "总监",
		},
		{
			name:         "15000元由CFO审批",
			amount:       15000,
			wantApproved: true,
			wantApprover: "CFO",
		},
		{
			name:         "25000元无法被审批",
			amount:       25000,
			wantApproved: false,
			wantApprover: "CFO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.Approve(tt.amount)

			assert := assert.New(t)
			assert.Equal(tt.wantApproved, result.Approved, "批准状态错误")
			assert.Equal(tt.wantApprover, result.Approver, "审批人错误")
			t.Logf("测试结果: %s", result.Message)
		})
	}
}

// 测试工厂方法创建的责任链
func TestCreateApprovalChain(t *testing.T) {
	chain := CreateApprovalChain()
	assert := assert.New(t)

	amounts := []float64{500, 1000, 1001, 5000, 5001, 20000, 20001}
	expectedApprovers := []string{"经理", "经理", "总监", "总监", "CFO", "CFO", "CFO"}
	expectedApproved := []bool{true, true, true, true, true, true, false}

	for i, amount := range amounts {
		t.Run(fmt.Sprintf("金额%.2f元", amount), func(t *testing.T) {
			result := chain.Approve(amount)

			assert.Equal(expectedApprovers[i], result.Approver, "审批人错误")
			assert.Equal(expectedApproved[i], result.Approved, "批准状态错误")

			// 检查消息内容是否合理
			if result.Approved {
				assert.Contains(result.Message, "批准", "批准消息不正确")
			} else {
				assert.Contains(result.Message, "超出", "拒绝消息不正确")
			}
		})
	}
}

// 测试边界值
func TestBoundaryValues(t *testing.T) {
	chain := CreateApprovalChain()
	assert := assert.New(t)

	tests := []struct {
		name         string
		amount       float64
		wantApproved bool
		wantApprover string
	}{
		{
			name:         "经理权限边界 - 999.99",
			amount:       999.99,
			wantApproved: true,
			wantApprover: "经理",
		},
		{
			name:         "经理权限边界 - 1000",
			amount:       1000,
			wantApproved: true,
			wantApprover: "经理",
		},
		{
			name:         "经理权限边界 - 1000.01",
			amount:       1000.01,
			wantApproved: true,
			wantApprover: "总监",
		},
		{
			name:         "总监权限边界 - 4999.99",
			amount:       4999.99,
			wantApproved: true,
			wantApprover: "总监",
		},
		{
			name:         "总监权限边界 - 5000",
			amount:       5000,
			wantApproved: true,
			wantApprover: "总监",
		},
		{
			name:         "总监权限边界 - 5000.01",
			amount:       5000.01,
			wantApproved: true,
			wantApprover: "CFO",
		},
		{
			name:         "CFO权限边界 - 19999.99",
			amount:       19999.99,
			wantApproved: true,
			wantApprover: "CFO",
		},
		{
			name:         "CFO权限边界 - 20000",
			amount:       20000,
			wantApproved: true,
			wantApprover: "CFO",
		},
		{
			name:         "CFO权限边界 - 20000.01",
			amount:       20000.01,
			wantApproved: false,
			wantApprover: "CFO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chain.Approve(tt.amount)

			assert.Equal(tt.wantApproved, result.Approved, "批准状态错误")
			assert.Equal(tt.wantApprover, result.Approver, "审批人错误")
		})
	}
}

// 测试自定义责任链
func TestCustomApprovalChain(t *testing.T) {
	assert := assert.New(t)

	// 创建一个具有特殊顺序的责任链：CFO -> 经理 -> 总监
	cfo := NewCFO(20000)
	manager := NewManager(1000)
	director := NewDirector(5000)

	cfo.SetNext(manager).SetNext(director)

	// 这种情况下，无论金额大小，都会先由CFO处理
	result := cfo.Approve(500)

	assert.True(result.Approved, "CFO应该批准500元的请求")
	assert.Equal("CFO", result.Approver, "应该由CFO处理请求")

	// 测试完整链条, 应该传递给经理
	cfo = NewCFO(100) // 修改CFO权限为100元
	manager = NewManager(1000)
	director = NewDirector(5000)

	cfo.SetNext(manager).SetNext(director)

	result = cfo.Approve(500)

	assert.True(result.Approved, "经理应该批准500元的请求")
	assert.Equal("经理", result.Approver, "应该由经理处理请求")
}

// 示例用法
func ExampleCreateApprovalChain() {
	chain := CreateApprovalChain()

	// 不同金额的请求
	amounts := []float64{500, 3000, 12000, 30000}

	for _, amount := range amounts {
		result := chain.Approve(amount)
		fmt.Printf("请求金额: %.2f, 结果: %s\n", amount, result.Message)
	}

	// Output:
	// 请求金额: 500.00, 结果: 经理批准了 500.00 元的请求
	// 请求金额: 3000.00, 结果: 总监批准了 3000.00 元的请求
	// 请求金额: 12000.00, 结果: CFO批准了 12000.00 元的请求
	// 请求金额: 30000.00, 结果: 请求金额 30000.00 超出了CFO的审批权限，无法批准
}
