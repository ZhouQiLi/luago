package vm

/*
float point byte
ABC模式下b和c操作数只有9位, 即便是无符号整数最大可表示的值只有512。
因为lua经常被当作配置表(类似json)所以若初始容量只有512则可能导致频繁扩容从而影响数据加载效率。
所以fpb的编码方式就是用来解释ABC模式下B和C操作符的。具体规则为：
如果把某个字节用二进制写成 eeeeexxx，那么当eeeee==0时该字节表示的整数就是xxx，否则该字节表示的整数是（1xxx）*2^（eeeee-1）。
*/
func Int2fb(x int) int {
	e := 0
	if x < 8 {
		return x
	}
	for x >= (8 << 4) {
		x = (x + 0xf) >> 4 // x = ceil(x /16)
		e += 4
	}
	for x >= (8 << 1) {
		x = (x + 1) >> 1 //  x = ceil(x / 2)
		e++
	}

	return ((e + 1) << 3) | (x - 8)
}

func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8) << uint((x>>3)-1)
	}
}
