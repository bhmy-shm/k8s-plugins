package lib

func DeployHeader() []string {
	//NAME     READY   UP-TO-DATE   AVAILABLE   AGE
	//mygott   2/2     2            2           7d14h
	headers := []string{"名称", "命名空间", "状态", "UP-TO-DATE", "在用数量", "使用时长"}
	return headers
}
