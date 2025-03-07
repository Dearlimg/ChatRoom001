package operate

//
//import (
//	"context"
//	"errors"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/redis/go-redis/v9"
//)

// 生成 Mock Redis 客户端（需先安装 gomock 和生成 mock）
// go get github.com/golang/mock/gomock
// mockgen -destination=mocks/mock_redis.go -package=mocks github.com/go-redis/redis/v9 Cmdable
//func TestRDB_AddEmails(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockRedis := mocks.NewMockCmdable(ctrl)
//	rdb := &RDB{rdb: mockRedis}
//
//	ctx := context.Background()
//
//	t.Run("成功添加多个邮箱", func(t *testing.T) {
//		emails := []string{"test1@example.com", "test2@example.com"}
//
//		// 预期调用 SAdd 并返回成功
//		mockRedis.EXPECT().SAdd(
//			ctx,
//			EmailKey,
//			gomock.Any(), // 可以具体验证参数
//		).Return(redis.NewIntResult(2, nil))
//
//		err := rdb.AddEmails(ctx, emails...)
//		if err != nil {
//			t.Errorf("预期无错误，但得到: %v", err)
//		}
//	})
//
//	t.Run("空邮箱列表不执行操作", func(t *testing.T) {
//		// 确保 SAdd 不会被调用
//		mockRedis.EXPECT().SAdd(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
//
//		err := rdb.AddEmails(ctx)
//		if err != nil {
//			t.Errorf("预期无错误，但得到: %v", err)
//		}
//	})
//
//	t.Run("处理Redis错误", func(t *testing.T) {
//		emails := []string{"test@example.com"}
//		expectedErr := errors.New("redis连接失败")
//
//		// 模拟返回错误
//		mockRedis.EXPECT().SAdd(
//			ctx,
//			EmailKey,
//			gomock.Eq([]interface{}{"test@example.com"}),
//		).Return(redis.NewIntResult(0, expectedErr))
//
//		err := rdb.AddEmails(ctx, emails...)
//		if err == nil || err.Error() != expectedErr.Error() {
//			t.Errorf("预期错误: %v, 实际得到: %v", expectedErr, err)
//		}
//	})
//}
