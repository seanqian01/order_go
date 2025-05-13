<template>
    <div class="app-wrapper">
      <sidebar class="sidebar-container" />
      <div class="main-container">
        <navbar />
        <app-main />
      </div>
    </div>
  </template>
  
  <script setup>
  import Sidebar from './components/Sidebar.vue'
  import Navbar from './components/Navbar.vue'
  import AppMain from './components/AppMain.vue'
  </script>
  
  <style scoped>
  .app-wrapper {
    position: relative;
    height: 100vh;
    width: 100%;
  }
  .sidebar-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 210px;
    height: 100%;
    background-color: #304156;
    overflow-y: auto;
    transition: width 0.28s;
    z-index: 1001;
  }
  .main-container {
    margin-left: 210px;
    position: relative;
    min-height: 100%;
    transition: margin-left 0.28s;
  }
  </style>