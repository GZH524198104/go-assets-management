# 上传
# file.csv格式
# 第一列 座位编号, 第二列 x坐标, 第三列 y坐标, 第四列优先级 POST
curl -k -F "file=@file.csv" https://localhost:8080/seats

# 按百分比获取抽查座位号 参数是百分比 0% ~ 100% GET
curl https://localhost:8080/seats/50

# 上面的接口返回值中的 imageUrl 字段 通过GET访问 下载图片

