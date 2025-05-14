<template>
  <router-view />
</template>

<script setup>
// 无需额外导入
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif;
}

#app {
  height: 100%;
}

/* 响应式表格样式 */
@media screen and (max-width: 768px) {
  .el-table {
    width: 100%;
    overflow-x: auto;
  }
  
  .el-pagination {
    white-space: normal;
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }
  
  /* 表单在移动端的样式优化 */
  .el-form-item {
    margin-bottom: 18px;
  }
  
  .el-form-item__label {
    text-align: left;
    float: none;
    display: block;
    margin-bottom: 8px;
  }
  
  .el-form-item__content {
    margin-left: 0 !important;
  }
  
  /* 对话框在移动端的样式优化 */
  .el-dialog {
    width: 90% !important;
    max-width: 500px;
  }
  
  /* 按钮在移动端的样式优化 */
  .el-button {
    margin-left: 0;
    margin-right: 8px;
    margin-bottom: 8px;
  }
  
  /* 卡片在移动端的样式优化 */
  .el-card {
    margin-bottom: 15px;
  }
}
</style>