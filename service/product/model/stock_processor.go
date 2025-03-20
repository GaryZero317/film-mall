package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

// StockOperation 定义库存操作类型
type StockOperationType int

const (
	DecrStock    StockOperationType = iota // 减少库存
	RestoreStock                           // 恢复库存
)

// StockOperation 库存操作请求
type StockOperation struct {
	Type     StockOperationType // 操作类型
	ID       int64              // 商品ID
	Quantity int64              // 操作数量
	Result   chan error         // 结果通道，用于异步获取操作结果
}

// StockProcessor 库存处理器
type StockProcessor struct {
	productModel ProductModel        // 商品模型
	queue        chan StockOperation // 操作队列
	workerCount  int                 // 工作协程数量
	mutex        sync.Mutex          // 防止重复初始化
	isRunning    bool                // 是否已启动
}

var (
	stockProcessorInstance *StockProcessor
	once                   sync.Once
)

// NewStockProcessor 创建库存处理器单例
func NewStockProcessor(productModel ProductModel, queueSize, workerCount int) *StockProcessor {
	once.Do(func() {
		stockProcessorInstance = &StockProcessor{
			productModel: productModel,
			queue:        make(chan StockOperation, queueSize),
			workerCount:  workerCount,
		}
	})
	return stockProcessorInstance
}

// Start 启动库存处理器
func (p *StockProcessor) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isRunning {
		return
	}

	p.isRunning = true
	logx.Info("[库存处理器] 启动")

	// 启动多个工作协程
	for i := 0; i < p.workerCount; i++ {
		threading.GoSafe(func() {
			p.worker(i)
		})
	}
}

// worker 工作协程函数
func (p *StockProcessor) worker(id int) {
	logx.Infof("[库存处理器] 工作协程 #%d 启动", id)

	for op := range p.queue {
		startTime := time.Now()
		var err error
		ctx := context.Background()

		switch op.Type {
		case DecrStock:
			logx.Infof("[库存处理器] 工作协程 #%d 处理减库存请求: 商品ID=%d, 数量=%d",
				id, op.ID, op.Quantity)
			err = p.productModel.DecrStock(ctx, op.ID, op.Quantity)
		case RestoreStock:
			logx.Infof("[库存处理器] 工作协程 #%d 处理恢复库存请求: 商品ID=%d, 数量=%d",
				id, op.ID, op.Quantity)
			err = p.productModel.RestoreStock(ctx, op.ID, op.Quantity)
		}

		// 记录处理结果
		elapsed := time.Since(startTime)
		if err != nil {
			logx.Errorf("[库存处理器] 工作协程 #%d 处理失败: 商品ID=%d, 耗时=%v, 错误=%v",
				id, op.ID, elapsed, err)
		} else {
			logx.Infof("[库存处理器] 工作协程 #%d 处理成功: 商品ID=%d, 耗时=%v",
				id, op.ID, elapsed)
		}

		// 发送结果
		if op.Result != nil {
			op.Result <- err
			close(op.Result)
		}
	}
}

// AsyncDecrStock 异步减少库存
func (p *StockProcessor) AsyncDecrStock(id, quantity int64) chan error {
	resultChan := make(chan error, 1)

	// 创建库存操作请求
	op := StockOperation{
		Type:     DecrStock,
		ID:       id,
		Quantity: quantity,
		Result:   resultChan,
	}

	// 发送到处理队列
	select {
	case p.queue <- op:
		logx.Infof("[库存处理器] 已将减库存请求加入队列: 商品ID=%d, 数量=%d", id, quantity)
	default:
		// 队列已满，直接报错
		go func() {
			resultChan <- fmt.Errorf("库存处理队列已满，无法处理请求")
			close(resultChan)
		}()
		logx.Errorf("[库存处理器] 队列已满，请求被拒绝: 商品ID=%d, 数量=%d", id, quantity)
	}

	return resultChan
}

// AsyncRestoreStock 异步恢复库存
func (p *StockProcessor) AsyncRestoreStock(id, quantity int64) chan error {
	resultChan := make(chan error, 1)

	// 创建库存操作请求
	op := StockOperation{
		Type:     RestoreStock,
		ID:       id,
		Quantity: quantity,
		Result:   resultChan,
	}

	// 发送到处理队列
	select {
	case p.queue <- op:
		logx.Infof("[库存处理器] 已将恢复库存请求加入队列: 商品ID=%d, 数量=%d", id, quantity)
	default:
		// 队列已满，直接报错
		go func() {
			resultChan <- fmt.Errorf("库存处理队列已满，无法处理请求")
			close(resultChan)
		}()
		logx.Errorf("[库存处理器] 队列已满，请求被拒绝: 商品ID=%d, 数量=%d", id, quantity)
	}

	return resultChan
}

// Stop 停止库存处理器
func (p *StockProcessor) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.isRunning {
		return
	}

	close(p.queue)
	p.isRunning = false
	logx.Info("[库存处理器] 已停止")
}
