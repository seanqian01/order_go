<template>
    <div class="navbar">
      <div class="hamburger-container">
        <!-- 这里可以放置汉堡菜单按钮 -->
      </div>
      <div class="right-menu">
        <el-dropdown class="avatar-container" trigger="click" @command="handleCommand">
          <div class="avatar-wrapper">
            <span>管理员</span>
            <el-icon><arrow-down /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </template>
  
  <script setup>
  import { useRouter } from 'vue-router'
  import { ElMessageBox } from 'element-plus'
  
  const router = useRouter()
  
  const handleCommand = (command) => {
    if (command === 'logout') {
      ElMessageBox.confirm('确认退出系统吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 清除token
        localStorage.removeItem('token')
        // 跳转到登录页
        router.push('/login')
      }).catch(() => {})
    } else if (command === 'profile') {
      // 跳转到个人信息页面
      // router.push('/profile')
    }
  }
  </script>
  
  <style scoped>
  .navbar {
    height: 50px;
    overflow: hidden;
    position: relative;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 20px;
  }
  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;
  }
  .right-menu {
    display: flex;
    align-items: center;
    height: 100%;
    margin-left: auto; /* 确保右侧菜单靠右 */
  }
  .avatar-container {
    cursor: pointer;
    margin-right: 0; /* 移除右侧边距 */
  }
  .avatar-wrapper {
    display: flex;
    align-items: center;
    padding: 0 8px;
    height: 100%;
  }
  .avatar-wrapper span {
    margin-right: 5px;
    font-size: 14px;
  }
  </style>