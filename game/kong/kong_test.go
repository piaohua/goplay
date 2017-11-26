/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-07-04 18:51:25
 * Filename      : kong_test.go
 * Description   : 玩牌算法
 * *******************************************************/
package kong

import "testing"

// 测试
func TestAlgo(t *testing.T) {
	cs := []uint32{0x4c, 0x42, 0x18}
	n, c, s := algo2(cs)
	t.Logf("n %x, c %x, s %x\n", n, c, s)
}

// 测试
func TestCompare(t *testing.T) {
	a := []uint32{0x42, 0x18}
	b := []uint32{0x38, 0x12}
	t.Log(Compare(a, b))
}

// 测试
func TestVerify(t *testing.T) {
	cs := []uint32{0x4c, 0x42, 0x18}
	val := TWOEIGHT
	t.Log(AlgoVerify(cs, val))
}
